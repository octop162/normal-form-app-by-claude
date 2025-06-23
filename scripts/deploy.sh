#!/bin/bash

# Deployment script for normal-form-app
# Usage: ./deploy.sh [environment] [action]
# Examples:
#   ./deploy.sh production deploy
#   ./deploy.sh staging rollback

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
AWS_REGION="${AWS_REGION:-ap-northeast-1}"
AWS_ACCOUNT_ID="${AWS_ACCOUNT_ID:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check required tools
check_dependencies() {
    local tools=("aws" "docker" "jq")
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "$tool is required but not installed"
            exit 1
        fi
    done
    
    log_info "All required tools are available"
}

# Function to validate AWS credentials
validate_aws_credentials() {
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials not configured or invalid"
        exit 1
    fi
    
    if [[ -z "$AWS_ACCOUNT_ID" ]]; then
        AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    fi
    
    log_info "AWS credentials validated for account: $AWS_ACCOUNT_ID"
}

# Function to build and push Docker images
build_and_push_images() {
    local environment=$1
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local git_hash=$(git rev-parse --short HEAD)
    local image_tag="${environment}-${timestamp}-${git_hash}"
    
    log_info "Building and pushing Docker images with tag: $image_tag"
    
    # Login to ECR
    aws ecr get-login-password --region "$AWS_REGION" | \
        docker login --username AWS --password-stdin \
        "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"
    
    # Build and push backend image
    log_info "Building backend image..."
    docker build -t "normal-form-app-backend:$image_tag" \
        -f "$PROJECT_ROOT/Dockerfile.backend" "$PROJECT_ROOT"
    
    docker tag "normal-form-app-backend:$image_tag" \
        "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:$image_tag"
    
    docker tag "normal-form-app-backend:$image_tag" \
        "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:latest"
    
    docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:$image_tag"
    docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:latest"
    
    # Build and push frontend image
    log_info "Building frontend image..."
    docker build -t "normal-form-app-frontend:$image_tag" \
        -f "$PROJECT_ROOT/frontend/Dockerfile" \
        --target production "$PROJECT_ROOT/frontend"
    
    docker tag "normal-form-app-frontend:$image_tag" \
        "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:$image_tag"
    
    docker tag "normal-form-app-frontend:$image_tag" \
        "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:latest"
    
    docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:$image_tag"
    docker push "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:latest"
    
    log_success "Docker images built and pushed successfully"
    echo "$image_tag"
}

# Function to update ECS service
update_ecs_service() {
    local environment=$1
    local cluster_name="normal-form-app-cluster-${environment}"
    local service_name="normal-form-app-service-${environment}"
    
    log_info "Updating ECS service: $service_name in cluster: $cluster_name"
    
    # Force new deployment
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --force-new-deployment \
        --region "$AWS_REGION"
    
    log_info "Waiting for service to become stable..."
    aws ecs wait services-stable \
        --cluster "$cluster_name" \
        --services "$service_name" \
        --region "$AWS_REGION"
    
    log_success "ECS service updated successfully"
}

# Function to perform health check
health_check() {
    local environment=$1
    local max_attempts=30
    local attempt=1
    
    # Get ALB DNS name
    local alb_name="normal-form-app-alb-${environment}"
    local alb_dns=$(aws elbv2 describe-load-balancers \
        --names "$alb_name" \
        --query 'LoadBalancers[0].DNSName' \
        --output text \
        --region "$AWS_REGION")
    
    log_info "Performing health check on https://$alb_dns"
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "https://$alb_dns/health" > /dev/null; then
            log_success "Health check passed!"
            return 0
        fi
        
        log_warning "Health check attempt $attempt/$max_attempts failed, retrying in 10 seconds..."
        sleep 10
        ((attempt++))
    done
    
    log_error "Health check failed after $max_attempts attempts"
    return 1
}

# Function to rollback to previous version
rollback_service() {
    local environment=$1
    local cluster_name="normal-form-app-cluster-${environment}"
    local service_name="normal-form-app-service-${environment}"
    
    log_warning "Initiating rollback for service: $service_name"
    
    # Get previous task definition
    local current_task_def=$(aws ecs describe-services \
        --cluster "$cluster_name" \
        --services "$service_name" \
        --query 'services[0].taskDefinition' \
        --output text \
        --region "$AWS_REGION")
    
    log_info "Current task definition: $current_task_def"
    
    # Get task definition family
    local family=$(echo "$current_task_def" | cut -d':' -f1)
    local current_revision=$(echo "$current_task_def" | cut -d':' -f2)
    local previous_revision=$((current_revision - 1))
    
    if [[ $previous_revision -lt 1 ]]; then
        log_error "No previous version available for rollback"
        return 1
    fi
    
    local previous_task_def="${family}:${previous_revision}"
    
    log_info "Rolling back to task definition: $previous_task_def"
    
    # Update service with previous task definition
    aws ecs update-service \
        --cluster "$cluster_name" \
        --service "$service_name" \
        --task-definition "$previous_task_def" \
        --region "$AWS_REGION"
    
    log_info "Waiting for rollback to complete..."
    aws ecs wait services-stable \
        --cluster "$cluster_name" \
        --services "$service_name" \
        --region "$AWS_REGION"
    
    log_success "Rollback completed successfully"
}

# Function to run database migrations
run_migrations() {
    local environment=$1
    local cluster_name="normal-form-app-cluster-${environment}"
    local task_definition="normal-form-app-migration-${environment}"
    
    log_info "Running database migrations for environment: $environment"
    
    # Run migration task
    local task_arn=$(aws ecs run-task \
        --cluster "$cluster_name" \
        --task-definition "$task_definition" \
        --launch-type FARGATE \
        --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
        --query 'tasks[0].taskArn' \
        --output text \
        --region "$AWS_REGION")
    
    log_info "Migration task started: $task_arn"
    
    # Wait for task to complete
    aws ecs wait tasks-stopped \
        --cluster "$cluster_name" \
        --tasks "$task_arn" \
        --region "$AWS_REGION"
    
    # Check task exit code
    local exit_code=$(aws ecs describe-tasks \
        --cluster "$cluster_name" \
        --tasks "$task_arn" \
        --query 'tasks[0].containers[0].exitCode' \
        --output text \
        --region "$AWS_REGION")
    
    if [[ "$exit_code" == "0" ]]; then
        log_success "Database migrations completed successfully"
    else
        log_error "Database migrations failed with exit code: $exit_code"
        return 1
    fi
}

# Main deployment function
deploy() {
    local environment=$1
    
    log_info "Starting deployment to environment: $environment"
    
    # Pre-deployment checks
    check_dependencies
    validate_aws_credentials
    
    # Build and push images
    local image_tag=$(build_and_push_images "$environment")
    
    # Run database migrations
    if [[ "$environment" == "production" ]]; then
        run_migrations "$environment"
    fi
    
    # Update ECS service
    update_ecs_service "$environment"
    
    # Perform health check
    if health_check "$environment"; then
        log_success "Deployment to $environment completed successfully!"
        log_info "Image tag: $image_tag"
    else
        log_error "Deployment health check failed"
        log_warning "Consider rolling back with: $0 $environment rollback"
        return 1
    fi
}

# Main function
main() {
    local environment=${1:-""}
    local action=${2:-"deploy"}
    
    if [[ -z "$environment" ]]; then
        log_error "Environment is required"
        echo "Usage: $0 [environment] [action]"
        echo "Environments: staging, production"
        echo "Actions: deploy, rollback"
        exit 1
    fi
    
    case "$action" in
        "deploy")
            deploy "$environment"
            ;;
        "rollback")
            rollback_service "$environment"
            ;;
        *)
            log_error "Unknown action: $action"
            echo "Available actions: deploy, rollback"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"