#!/bin/bash

# Database migration script for normal-form-app
# Usage: ./scripts/migrate.sh [up|down|drop|version] [steps]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-normal_form_app}
DB_SSLMODE=${DB_SSLMODE:-disable}

# Construct database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Migration directory
MIGRATIONS_DIR="./migrations"

# Check if migrate command is available
if ! command -v migrate &> /dev/null; then
    print_error "migrate command not found. Please install golang-migrate:"
    echo "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    print_error "Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

show_usage() {
    echo "Usage: $0 [command] [steps]"
    echo ""
    echo "Commands:"
    echo "  up [N]      Apply all or N up migrations"
    echo "  down [N]    Apply all or N down migrations"
    echo "  drop        Drop everything inside database"
    echo "  force V     Set version V but don't run migration (ignores dirty state)"
    echo "  version     Print current migration version"
    echo "  create NAME Create new migration files"
    echo ""
    echo "Examples:"
    echo "  $0 up               # Apply all pending migrations"
    echo "  $0 up 1             # Apply next 1 migration"
    echo "  $0 down 1           # Rollback 1 migration"
    echo "  $0 version          # Show current version"
    echo "  $0 create add_user  # Create new migration files"
}

# Function to check database connection
check_db_connection() {
    print_status "Checking database connection..."
    
    if command -v psql &> /dev/null; then
        if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' &> /dev/null; then
            print_error "Cannot connect to database. Please check your database settings."
            return 1
        fi
        print_success "Database connection successful"
    else
        print_warning "psql not found, skipping connection check"
    fi
}

# Function to run migration
run_migrate() {
    local cmd="$1"
    local steps="$2"
    
    print_status "Running migration: $cmd $steps"
    print_status "Database URL: postgres://${DB_USER}:***@${DB_HOST}:${DB_PORT}/${DB_NAME}"
    
    case "$cmd" in
        "up")
            if [ -n "$steps" ]; then
                migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up "$steps"
            else
                migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
            fi
            ;;
        "down")
            if [ -n "$steps" ]; then
                migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down "$steps"
            else
                print_warning "Running down without steps will revert ALL migrations!"
                read -p "Are you sure? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down
                else
                    print_status "Migration cancelled"
                    exit 0
                fi
            fi
            ;;
        "drop")
            print_warning "This will drop everything in the database!"
            read -p "Are you sure? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" drop
            else
                print_status "Drop cancelled"
                exit 0
            fi
            ;;
        "force")
            if [ -z "$steps" ]; then
                print_error "Force command requires version number"
                exit 1
            fi
            migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$steps"
            ;;
        "version")
            migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
            ;;
        "create")
            if [ -z "$steps" ]; then
                print_error "Create command requires migration name"
                exit 1
            fi
            migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$steps"
            print_success "Created migration files for: $steps"
            ;;
        *)
            print_error "Unknown command: $cmd"
            show_usage
            exit 1
            ;;
    esac
}

# Main execution
main() {
    local command="${1:-up}"
    local steps="$2"
    
    if [ "$command" = "help" ] || [ "$command" = "-h" ] || [ "$command" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    print_status "Database Migration Tool"
    echo "Database: $DB_NAME@$DB_HOST:$DB_PORT"
    echo ""
    
    if [ "$command" != "create" ]; then
        check_db_connection
    fi
    
    run_migrate "$command" "$steps"
    
    if [ $? -eq 0 ]; then
        print_success "Migration completed successfully"
    else
        print_error "Migration failed"
        exit 1
    fi
}

main "$@"