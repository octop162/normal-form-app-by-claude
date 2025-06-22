#!/bin/bash

# 会員登録webフォーム - 開発環境停止スクリプト
# Usage: ./scripts/dev-stop.sh [option]
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

stop_database() {
    print_status "Stopping PostgreSQL database..."
    
    docker-compose stop postgres 2>/dev/null || true
    docker-compose rm -f postgres 2>/dev/null || true
    
    print_success "PostgreSQL stopped"
}

stop_backend() {
    print_status "Stopping Go backend..."
    
    # Try to kill using PID file first
    if [ -f "logs/backend.pid" ]; then
        PID=$(cat logs/backend.pid)
        if kill -0 $PID 2>/dev/null; then
            kill $PID
            print_success "Stopped Go backend (PID: $PID)"
        fi
        rm -f logs/backend.pid
    fi
    
    # Fallback: kill by process name
    pkill -f "go run cmd/server/main.go" 2>/dev/null && print_status "Killed Go processes by name" || true
    pkill -f "cmd/server/main.go" 2>/dev/null || true
    
    print_success "Go backend stopped"
}

stop_frontend() {
    print_status "Stopping React frontend..."
    
    # Try to kill using PID file first
    if [ -f "logs/frontend.pid" ]; then
        PID=$(cat logs/frontend.pid)
        if kill -0 $PID 2>/dev/null; then
            kill $PID
            print_success "Stopped React frontend (PID: $PID)"
        fi
        rm -f logs/frontend.pid
    fi
    
    # Fallback: kill by process name
    pkill -f "vite" 2>/dev/null && print_status "Killed Vite processes" || true
    pkill -f "npm run dev" 2>/dev/null && print_status "Killed npm dev processes" || true
    
    print_success "React frontend stopped"
}

stop_docker_environment() {
    print_status "Stopping Docker environment..."
    
    docker-compose down
    
    print_success "Docker environment stopped"
}

stop_all_processes() {
    print_status "Stopping all development services..."
    
    stop_frontend
    stop_backend
    stop_database
    
    # Clean up any remaining processes
    print_status "Cleaning up remaining processes..."
    pkill -f "node.*vite" 2>/dev/null || true
    pkill -f "go run" 2>/dev/null || true
    
    # Clean up log files
    if [ -d "logs" ]; then
        rm -f logs/*.pid
    fi
    
    print_success "All services stopped"
}

force_cleanup() {
    print_warning "Performing force cleanup..."
    
    # Force kill all related processes
    pkill -9 -f "go run" 2>/dev/null || true
    pkill -9 -f "vite" 2>/dev/null || true
    pkill -9 -f "npm run dev" 2>/dev/null || true
    pkill -9 -f "node.*vite" 2>/dev/null || true
    
    # Stop all Docker containers
    docker-compose down --remove-orphans 2>/dev/null || true
    
    # Clean up log files and PID files
    rm -f logs/*.pid 2>/dev/null || true
    
    print_success "Force cleanup completed"
}

remove_data() {
    print_warning "Removing all data (including database)..."
    
    # Stop everything first
    stop_all_processes
    
    # Remove Docker volumes (this will delete database data)
    docker-compose down -v --remove-orphans
    
    # Remove log files
    rm -rf logs/ 2>/dev/null || true
    
    print_warning "All data removed. Database will be recreated on next start."
}

show_status() {
    echo ""
    echo "=== Development Environment Status ==="
    echo ""
    
    # Check PostgreSQL
    if docker-compose ps postgres 2>/dev/null | grep -q "Up"; then
        echo -e "PostgreSQL: ${GREEN}Running${NC}"
    else
        echo -e "PostgreSQL: ${RED}Stopped${NC}"
    fi
    
    # Check Go backend
    if pgrep -f "go run cmd/server/main.go" > /dev/null; then
        echo -e "Go Backend: ${GREEN}Running${NC}"
    else
        echo -e "Go Backend: ${RED}Stopped${NC}"
    fi
    
    # Check React frontend
    if pgrep -f "vite" > /dev/null; then
        echo -e "React Frontend: ${GREEN}Running${NC}"
    else
        echo -e "React Frontend: ${RED}Stopped${NC}"
    fi
    
    echo ""
}

show_help() {
    echo "会員登録webフォーム - 開発環境停止スクリプト"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  all       Stop all services (database + backend + frontend) [default]"
    echo "  db        Stop only PostgreSQL database"
    echo "  backend   Stop only Go backend"
    echo "  frontend  Stop only React frontend"
    echo "  docker    Stop Docker environment"
    echo "  force     Force stop all processes (use if normal stop fails)"
    echo "  clean     Stop all services and remove data/logs"
    echo "  status    Show current services status"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Stop all services"
    echo "  $0 all               # Stop all services"
    echo "  $0 db                # Stop only database"
    echo "  $0 backend           # Stop only backend"
    echo "  $0 frontend          # Stop only frontend"
    echo "  $0 force             # Force stop everything"
    echo "  $0 clean             # Stop and remove all data"
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
            stop_all_processes
            show_status
            ;;
        "db"|"database")
            stop_database
            ;;
        "backend"|"api")
            stop_backend
            ;;
        "frontend"|"ui")
            stop_frontend
            ;;
        "docker")
            stop_docker_environment
            ;;
        "force")
            force_cleanup
            show_status
            ;;
        "clean"|"reset")
            remove_data
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
trap 'print_warning "Script interrupted by user"; exit 0' INT

# Run main function
main "$@"