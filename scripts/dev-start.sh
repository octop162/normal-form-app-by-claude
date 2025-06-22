#!/bin/bash

# 会員登録webフォーム - 開発環境起動スクリプト
# Usage: ./scripts/dev-start.sh [option]
# Options: all, db, backend, frontend, docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Functions
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

wait_for_service() {
    local service_name=$1
    local url=$2
    local max_attempts=30
    local attempt=1

    print_status "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to start within ${max_attempts} seconds"
    return 1
}

check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if required commands exist
    for cmd in docker docker-compose go node npm; do
        if ! command -v $cmd &> /dev/null; then
            print_error "$cmd is not installed"
            exit 1
        fi
    done
    
    # Check Go version
    go_version=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
    if [[ $(echo "$go_version" | cut -d'.' -f1) -lt 1 ]] || [[ $(echo "$go_version" | cut -d'.' -f2) -lt 21 ]]; then
        print_warning "Go version $go_version detected. Recommended: 1.21+"
    fi
    
    # Check Node version  
    node_version=$(node --version | cut -d'v' -f2)
    if [[ $(echo "$node_version" | cut -d'.' -f1) -lt 18 ]]; then
        print_warning "Node.js version $node_version detected. Recommended: 18+"
    fi
    
    print_success "Prerequisites check completed"
}

setup_environment() {
    print_status "Setting up environment..."
    
    # Copy .env if not exists
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            print_success "Created .env from .env.example"
        else
            print_warning ".env.example not found, creating basic .env"
            cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_NAME=normal_form_db
DB_USER=postgres
DB_PASSWORD=postgres
LOG_LEVEL=debug
SESSION_TIMEOUT=4h
PORT=8080
NODE_ENV=development
GO_ENV=development
EOF
        fi
    fi
    
    # Install Go dependencies
    if [ -f "go.mod" ]; then
        print_status "Installing Go dependencies..."
        go mod download
        go mod tidy
    fi
    
    # Install Node dependencies
    if [ -f "frontend/package.json" ]; then
        print_status "Installing Node.js dependencies..."
        cd frontend
        if [ ! -d "node_modules" ]; then
            npm install
        else
            print_status "Node modules already installed, skipping..."
        fi
        cd ..
    fi
    
    print_success "Environment setup completed"
}

start_database() {
    print_status "Starting PostgreSQL database..."
    
    # Stop existing containers
    docker-compose stop postgres 2>/dev/null || true
    
    # Start PostgreSQL
    docker-compose up -d postgres
    
    # Wait for database to be ready
    wait_for_service "PostgreSQL" "http://localhost:5432" || {
        print_error "Failed to start PostgreSQL"
        return 1
    }
    
    print_success "PostgreSQL started successfully"
}

start_backend() {
    print_status "Starting Go backend..."
    
    # Kill existing Go processes
    pkill -f "go run cmd/server/main.go" 2>/dev/null || true
    pkill -f "cmd/server/main.go" 2>/dev/null || true
    
    # Wait a moment for processes to terminate
    sleep 2
    
    # Start Go server in background
    nohup go run cmd/server/main.go > logs/backend.log 2>&1 &
    GO_PID=$!
    echo $GO_PID > logs/backend.pid
    
    # Wait for backend to be ready
    wait_for_service "Go Backend" "http://localhost:8080/health" || {
        print_error "Failed to start Go backend"
        return 1
    }
    
    print_success "Go backend started successfully (PID: $GO_PID)"
}

start_frontend() {
    print_status "Starting React frontend..."
    
    # Kill existing Node processes
    pkill -f "vite" 2>/dev/null || true
    pkill -f "npm run dev" 2>/dev/null || true
    
    # Wait a moment for processes to terminate
    sleep 2
    
    # Start React dev server in background
    cd frontend
    nohup npm run dev > ../logs/frontend.log 2>&1 &
    REACT_PID=$!
    echo $REACT_PID > ../logs/frontend.pid
    cd ..
    
    # Wait for frontend to be ready
    wait_for_service "React Frontend" "http://localhost:5173" || {
        print_error "Failed to start React frontend"
        return 1
    }
    
    print_success "React frontend started successfully (PID: $REACT_PID)"
}

start_docker_environment() {
    print_status "Starting complete Docker environment..."
    
    # Stop existing containers
    docker-compose down 2>/dev/null || true
    
    # Start all services with Docker
    docker-compose --profile backend --profile frontend up -d
    
    # Wait for all services
    wait_for_service "PostgreSQL" "http://localhost:5432"
    wait_for_service "Go Backend" "http://localhost:8080/health"
    wait_for_service "React Frontend" "http://localhost:3000"
    
    print_success "Docker environment started successfully"
}

create_logs_directory() {
    mkdir -p logs
}

show_status() {
    echo ""
    echo "=== Development Environment Status ==="
    echo ""
    
    # Check PostgreSQL
    if docker-compose ps postgres | grep -q "Up"; then
        echo -e "PostgreSQL: ${GREEN}Running${NC} (localhost:5432)"
    else
        echo -e "PostgreSQL: ${RED}Not running${NC}"
    fi
    
    # Check Go backend
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "Go Backend: ${GREEN}Running${NC} (localhost:8080)"
    else
        echo -e "Go Backend: ${RED}Not running${NC}"
    fi
    
    # Check React frontend
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        echo -e "React Frontend: ${GREEN}Running${NC} (localhost:5173)"
    elif curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo -e "React Frontend: ${GREEN}Running${NC} (localhost:3000 - Docker)"
    else
        echo -e "React Frontend: ${RED}Not running${NC}"
    fi
    
    echo ""
    echo "To access the application:"
    echo "  Frontend: http://localhost:5173 (or http://localhost:3000 for Docker)"
    echo "  Backend API: http://localhost:8080"
    echo "  Health Check: http://localhost:8080/health"
    echo ""
}

show_help() {
    echo "会員登録webフォーム - 開発環境起動スクリプト"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  all       Start all services (database + backend + frontend) [default]"
    echo "  db        Start only PostgreSQL database"
    echo "  backend   Start only Go backend (requires database)"
    echo "  frontend  Start only React frontend"
    echo "  docker    Start complete environment using Docker"
    echo "  status    Show current services status"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Start all services"
    echo "  $0 all               # Start all services"
    echo "  $0 db                # Start only database"
    echo "  $0 backend           # Start only backend"
    echo "  $0 frontend          # Start only frontend"
    echo "  $0 docker            # Use Docker for all services"
    echo ""
}

# Main execution
main() {
    local option="${1:-all}"
    
    case "$option" in
        "help"|"-h"|"--help")
            show_help
            exit 0
            ;;
        "status")
            show_status
            exit 0
            ;;
        "all")
            create_logs_directory
            check_prerequisites
            setup_environment
            start_database
            start_backend
            start_frontend
            show_status
            ;;
        "db"|"database")
            check_prerequisites
            start_database
            ;;
        "backend"|"api")
            create_logs_directory
            check_prerequisites
            setup_environment
            start_backend
            ;;
        "frontend"|"ui")
            create_logs_directory
            check_prerequisites
            setup_environment
            start_frontend
            ;;
        "docker")
            check_prerequisites
            setup_environment
            start_docker_environment
            show_status
            ;;
        *)
            print_error "Unknown option: $option"
            show_help
            exit 1
            ;;
    esac
}

# Trap to handle script interruption
trap 'print_warning "Script interrupted by user"; exit 1' INT

# Run main function
main "$@"