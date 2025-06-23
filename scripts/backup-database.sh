#!/bin/bash

# Database backup script for normal-form-app
# This script creates encrypted backups of the PostgreSQL database

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/tmp/backups"
S3_BUCKET="${S3_BACKUP_BUCKET:-normal-form-app-backup-${AWS_ACCOUNT_ID}}"
KMS_KEY_ID="${BACKUP_KMS_KEY_ID:-alias/normal-form-app-backup}"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"

# Database configuration
DB_HOST="${DB_HOST:-}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-}"
DB_USER="${DB_USER:-}"
DB_PASSWORD="${DB_PASSWORD:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function to check prerequisites
check_prerequisites() {
    local tools=("pg_dump" "gzip" "aws" "gpg")
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "$tool is required but not installed"
            exit 1
        fi
    done
    
    if [[ -z "$DB_HOST" || -z "$DB_NAME" || -z "$DB_USER" || -z "$DB_PASSWORD" ]]; then
        log_error "Database connection parameters are required"
        exit 1
    fi
    
    log_info "Prerequisites check passed"
}

# Function to create backup directory
create_backup_dir() {
    if [[ ! -d "$BACKUP_DIR" ]]; then
        mkdir -p "$BACKUP_DIR"
        log_info "Created backup directory: $BACKUP_DIR"
    fi
}

# Function to test database connection
test_db_connection() {
    log_info "Testing database connection..."
    
    export PGPASSWORD="$DB_PASSWORD"
    
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &> /dev/null; then
        log_success "Database connection successful"
    else
        log_error "Failed to connect to database"
        exit 1
    fi
}

# Function to create database backup
create_backup() {
    local backup_file="$BACKUP_DIR/db_backup_${DB_NAME}_${BACKUP_DATE}.sql"
    local compressed_file="${backup_file}.gz"
    local encrypted_file="${compressed_file}.gpg"
    
    log_info "Starting database backup..."
    
    export PGPASSWORD="$DB_PASSWORD"
    
    # Create SQL dump
    if pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
        --verbose --clean --no-owner --no-privileges \
        --format=plain > "$backup_file"; then
        log_success "Database dump created: $backup_file"
    else
        log_error "Failed to create database dump"
        exit 1
    fi
    
    # Compress the backup
    if gzip "$backup_file"; then
        log_success "Backup compressed: $compressed_file"
    else
        log_error "Failed to compress backup"
        exit 1
    fi
    
    # Encrypt the backup
    if gpg --trust-model always --symmetric --cipher-algo AES256 \
        --output "$encrypted_file" "$compressed_file"; then
        log_success "Backup encrypted: $encrypted_file"
        rm "$compressed_file"  # Remove unencrypted file
    else
        log_error "Failed to encrypt backup"
        exit 1
    fi
    
    echo "$encrypted_file"
}

# Function to upload backup to S3
upload_to_s3() {
    local backup_file="$1"
    local s3_key="database/$(basename "$backup_file")"
    
    log_info "Uploading backup to S3: s3://$S3_BUCKET/$s3_key"
    
    if aws s3 cp "$backup_file" "s3://$S3_BUCKET/$s3_key" \
        --server-side-encryption aws:kms \
        --ssekms-key-id "$KMS_KEY_ID" \
        --storage-class STANDARD_IA; then
        log_success "Backup uploaded to S3"
    else
        log_error "Failed to upload backup to S3"
        exit 1
    fi
    
    # Verify upload
    if aws s3 ls "s3://$S3_BUCKET/$s3_key" &> /dev/null; then
        log_success "Backup verified in S3"
    else
        log_error "Backup verification failed"
        exit 1
    fi
}

# Function to cleanup old backups
cleanup_old_backups() {
    log_info "Cleaning up backups older than $RETENTION_DAYS days..."
    
    # Cleanup local backups
    find "$BACKUP_DIR" -name "db_backup_*.sql.gz.gpg" -type f -mtime +$RETENTION_DAYS -delete || true
    
    # Cleanup S3 backups (handled by S3 lifecycle policy)
    local cutoff_date=$(date -d "$RETENTION_DAYS days ago" +%Y-%m-%d)
    log_info "S3 backups older than $cutoff_date will be handled by lifecycle policy"
    
    log_success "Cleanup completed"
}

# Function to create backup metadata
create_metadata() {
    local backup_file="$1"
    local metadata_file="${backup_file}.metadata.json"
    
    cat > "$metadata_file" << EOF
{
    "backup_date": "$BACKUP_DATE",
    "database_name": "$DB_NAME",
    "database_host": "$DB_HOST",
    "backup_file": "$(basename "$backup_file")",
    "file_size": $(stat -c%s "$backup_file"),
    "checksum": "$(sha256sum "$backup_file" | cut -d' ' -f1)",
    "encryption": "GPG-AES256",
    "compression": "gzip",
    "created_by": "$(whoami)",
    "script_version": "1.0",
    "retention_days": $RETENTION_DAYS
}
EOF
    
    # Upload metadata to S3
    local metadata_s3_key="database/metadata/$(basename "$metadata_file")"
    aws s3 cp "$metadata_file" "s3://$S3_BUCKET/$metadata_s3_key" \
        --server-side-encryption aws:kms \
        --ssekms-key-id "$KMS_KEY_ID"
    
    log_success "Metadata created and uploaded"
}

# Function to send notification
send_notification() {
    local status="$1"
    local message="$2"
    local sns_topic="${BACKUP_NOTIFICATION_SNS_ARN:-}"
    
    if [[ -n "$sns_topic" ]]; then
        local subject="Database Backup $status - $(date '+%Y-%m-%d %H:%M:%S')"
        
        aws sns publish \
            --topic-arn "$sns_topic" \
            --subject "$subject" \
            --message "$message" || true
    fi
}

# Function to validate backup
validate_backup() {
    local backup_file="$1"
    
    log_info "Validating backup integrity..."
    
    # Check if file exists and is not empty
    if [[ ! -f "$backup_file" || ! -s "$backup_file" ]]; then
        log_error "Backup file is missing or empty"
        return 1
    fi
    
    # Test GPG decryption (without actually decrypting)
    if gpg --list-packets "$backup_file" &> /dev/null; then
        log_success "Backup encryption validation passed"
    else
        log_error "Backup encryption validation failed"
        return 1
    fi
    
    log_success "Backup validation completed"
}

# Main backup function
main() {
    local start_time=$(date +%s)
    
    log_info "Starting database backup process..."
    
    # Setup
    check_prerequisites
    create_backup_dir
    test_db_connection
    
    # Create backup
    local backup_file
    backup_file=$(create_backup)
    
    # Validate backup
    if validate_backup "$backup_file"; then
        # Upload to S3
        upload_to_s3 "$backup_file"
        
        # Create metadata
        create_metadata "$backup_file"
        
        # Cleanup
        cleanup_old_backups
        
        # Calculate duration
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        local success_message="Database backup completed successfully in ${duration} seconds"
        log_success "$success_message"
        send_notification "SUCCESS" "$success_message"
        
        # Cleanup local backup file
        rm "$backup_file"
        
    else
        local error_message="Database backup validation failed"
        log_error "$error_message"
        send_notification "FAILED" "$error_message"
        exit 1
    fi
}

# Trap for cleanup on exit
cleanup_on_exit() {
    if [[ -d "$BACKUP_DIR" ]]; then
        find "$BACKUP_DIR" -name "db_backup_*" -type f -delete 2>/dev/null || true
    fi
}

trap cleanup_on_exit EXIT

# Run main function
main "$@"