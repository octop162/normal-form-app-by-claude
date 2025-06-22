#!/bin/bash

# 会員登録webフォーム - ヘルスチェックスクリプト
# Usage: ./scripts/health-check.sh [option]
# Options: all, db, backend, frontend, detailed

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Functions
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_header() {
    echo -e "${CYAN}=== $1 ===${NC}"
}

check_command() {
    local cmd=$1
    local desc=$2
    
    if command -v $cmd &> /dev/null; then
        local version=$($cmd --version 2>/dev/null | head -1 || echo "Unknown version")
        print_success "$desc: $version"
        return 0
    else
        print_error "$desc: Not installed"
        return 1
    fi
}

check_port() {
    local port=$1
    local service=$2
    local timeout=${3:-3}
    
    if timeout $timeout bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null; then
        print_success "$service (port $port): Listening"
        return 0
    else
        print_error "$service (port $port): Not listening"
        return 1
    fi
}

check_http_endpoint() {
    local url=$1
    local service=$2
    local timeout=${3:-5}
    
    local response=$(curl -s -w "HTTP_CODE:%{http_code},TIME:%{time_total}" --max-time $timeout "$url" 2>/dev/null || echo "FAILED")
    
    if [[ $response == *"HTTP_CODE:200"* ]]; then
        local time=$(echo "$response" | grep -o "TIME:[0-9.]*" | cut -d: -f2)
        print_success "$service: Healthy (${time}s)"
        return 0
    elif [[ $response == *"HTTP_CODE:"* ]]; then
        local code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
        print_warning "$service: HTTP $code"
        return 1
    else
        print_error "$service: Connection failed"
        return 1
    fi
}

check_prerequisites() {
    print_header "Prerequisites Check"
    
    local all_good=0
    
    check_command "docker" "Docker" || all_good=1
    check_command "docker-compose" "Docker Compose" || all_good=1
    check_command "go" "Go" || all_good=1
    check_command "node" "Node.js" || all_good=1
    check_command "npm" "npm" || all_good=1
    check_command "curl" "curl" || all_good=1
    
    echo ""
    return $all_good
}

check_database() {
    print_header "Database Health Check"
    
    local db_status=0
    
    # Check if PostgreSQL container is running
    if docker-compose ps postgres 2>/dev/null | grep -q "Up"; then
        print_success "PostgreSQL container: Running"
        
        # Check if port is accessible
        if check_port 5432 "PostgreSQL"; then
            # Try to connect to database
            if docker exec normal-form-db psql -U postgres -d normal_form_db -c "SELECT 1;" &>/dev/null; then
                print_success "PostgreSQL connection: OK"
                
                # Check health_check table
                local health_status=$(docker exec normal-form-db psql -U postgres -d normal_form_db -t -c "SELECT status FROM health_check LIMIT 1;" 2>/dev/null | tr -d ' \n' || echo "")
                if [ "$health_status" = "ok" ]; then
                    print_success "Database health check: OK"
                else
                    print_warning "Database health check: No health record found"
                    db_status=1
                fi
            else
                print_error "PostgreSQL connection: Failed"
                db_status=1
            fi
        else
            db_status=1
        fi
    else
        print_error "PostgreSQL container: Not running"
        db_status=1
    fi
    
    echo ""
    return $db_status
}

check_backend() {
    print_header "Backend Health Check"
    
    local backend_status=0
    
    # Check if Go process is running
    if pgrep -f "go run cmd/server/main.go" > /dev/null; then
        print_success "Go backend process: Running"
    elif pgrep -f "cmd/server/main.go" > /dev/null; then
        print_success "Go backend process: Running"
    else
        print_error "Go backend process: Not running"
        backend_status=1
    fi
    
    # Check health endpoint
    if check_http_endpoint "http://localhost:8080/health" "Health endpoint"; then
        # Test API endpoint
        if check_http_endpoint "http://localhost:8080/api/v1/ping" "API endpoint"; then
            # Get detailed health info
            local health_info=$(curl -s http://localhost:8080/health 2>/dev/null || echo "{}")
            local service=$(echo "$health_info" | grep -o '"service":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "unknown")
            local version=$(echo "$health_info" | grep -o '"version":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "unknown")
            print_success "Service info: $service v$version"
        else
            backend_status=1
        fi
    else
        backend_status=1
    fi
    
    echo ""
    return $backend_status
}

check_frontend() {
    print_header "Frontend Health Check"
    
    local frontend_status=0
    
    # Check if Vite/React process is running
    if pgrep -f "vite" > /dev/null; then
        print_success "React frontend process: Running"
    else
        print_error "React frontend process: Not running"
        frontend_status=1
    fi
    
    # Check frontend endpoints
    if check_http_endpoint "http://localhost:5173" "Frontend (Vite dev server)"; then
        true
    elif check_http_endpoint "http://localhost:3000" "Frontend (Docker)"; then
        true
    else
        print_error "Frontend: Not accessible on any port"
        frontend_status=1
    fi
    
    echo ""
    return $frontend_status
}

check_integration() {
    print_header "Integration Health Check"
    
    local integration_status=0
    
    # Check if all services are running
    print_status "Testing service integration..."
    
    # Test database → backend connectivity
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        print_success "Backend → Database: Connected"
    else
        print_error "Backend → Database: Connection issues"
        integration_status=1
    fi
    
    # Test API endpoints that might use database
    # (These will be implemented in later phases)
    print_status "API endpoints ready for phase 2 implementation"
    
    echo ""
    return $integration_status
}

detailed_system_info() {
    print_header "Detailed System Information"
    
    echo -e "${CYAN}System:${NC}"
    uname -a
    echo ""
    
    echo -e "${CYAN}Docker:${NC}"
    docker --version
    docker-compose --version
    echo ""
    
    echo -e "${CYAN}Go:${NC}"
    go version
    echo ""
    
    echo -e "${CYAN}Node.js:${NC}"
    node --version
    npm --version
    echo ""
    
    echo -e "${CYAN}Memory Usage:${NC}"
    free -h
    echo ""
    
    echo -e "${CYAN}Disk Usage:${NC}"
    df -h . | head -2
    echo ""
    
    echo -e "${CYAN}Running Processes:${NC}"
    ps aux | grep -E "(go run|vite|npm|postgres)" | grep -v grep || echo "No development processes found"
    echo ""
    
    echo -e "${CYAN}Docker Containers:${NC}"
    docker-compose ps 2>/dev/null || echo "No containers running"
    echo ""
    
    echo -e "${CYAN}Network Ports:${NC}"
    netstat -tlnp 2>/dev/null | grep -E ":(5432|8080|5173|3000)" || echo "No development ports listening"
    echo ""
}

run_quick_tests() {
    print_header "Quick Functional Tests"
    
    local test_status=0
    
    # Test backend health endpoint
    print_status "Testing backend health endpoint..."
    local health_response=$(curl -s http://localhost:8080/health 2>/dev/null || echo "FAILED")
    if echo "$health_response" | grep -q "ok"; then
        print_success "Health endpoint: Returns valid response"
    else
        print_error "Health endpoint: Invalid response"
        test_status=1
    fi
    
    # Test backend API endpoint
    print_status "Testing backend API endpoint..."
    local ping_response=$(curl -s http://localhost:8080/api/v1/ping 2>/dev/null || echo "FAILED")
    if echo "$ping_response" | grep -q "pong"; then
        print_success "API endpoint: Returns valid response"
    else
        print_error "API endpoint: Invalid response"
        test_status=1
    fi
    
    # Test database connectivity (if running)
    if docker-compose ps postgres 2>/dev/null | grep -q "Up"; then
        print_status "Testing database connectivity..."
        if docker exec normal-form-db psql -U postgres -d normal_form_db -c "SELECT COUNT(*) FROM health_check;" &>/dev/null; then
            local count=$(docker exec normal-form-db psql -U postgres -d normal_form_db -t -c "SELECT COUNT(*) FROM health_check;" 2>/dev/null | tr -d ' \n' || echo "0")
            print_success "Database query: Returns $count health records"
        else
            print_error "Database query: Failed"
            test_status=1
        fi
    fi
    
    echo ""
    return $test_status
}

show_summary() {
    print_header "Health Check Summary"
    
    local overall_status=0
    
    echo -e "${CYAN}Service Status:${NC}"
    
    # PostgreSQL
    if docker-compose ps postgres 2>/dev/null | grep -q "Up"; then
        echo -e "  PostgreSQL: ${GREEN}✓ Running${NC} (localhost:5432)"
    else
        echo -e "  PostgreSQL: ${RED}✗ Stopped${NC}"
        overall_status=1
    fi
    
    # Go Backend
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "  Go Backend: ${GREEN}✓ Running${NC} (localhost:8080)"
    else
        echo -e "  Go Backend: ${RED}✗ Not responding${NC}"
        overall_status=1
    fi
    
    # React Frontend
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        echo -e "  React Frontend: ${GREEN}✓ Running${NC} (localhost:5173)"
    elif curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo -e "  React Frontend: ${GREEN}✓ Running${NC} (localhost:3000)"
    else
        echo -e "  React Frontend: ${RED}✗ Not responding${NC}"
        overall_status=1
    fi
    
    echo ""
    
    if [ $overall_status -eq 0 ]; then
        print_success "All services are healthy! 🎉"
        echo ""
        echo -e "${CYAN}Ready for development:${NC}"
        echo "  • Frontend: http://localhost:5173"
        echo "  • Backend API: http://localhost:8080"
        echo "  • Health Check: http://localhost:8080/health"
    else
        print_error "Some services need attention ⚠️"
        echo ""
        echo "Run './scripts/dev-start.sh' to start missing services"
    fi
    
    echo ""
    return $overall_status
}

show_help() {
    echo "会員登録webフォーム - ヘルスチェックスクリプト"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  all       Check all services [default]"
    echo "  db        Check only database"
    echo "  backend   Check only backend API"
    echo "  frontend  Check only frontend"
    echo "  detailed  Show detailed system information"
    echo "  quick     Run quick functional tests"
    echo "  summary   Show summary only"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Check all services"
    echo "  $0 all               # Check all services"
    echo "  $0 db                # Check only database"
    echo "  $0 backend           # Check only backend"
    echo "  $0 detailed          # Show detailed system info"
    echo "  $0 quick             # Run quick tests"
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
        "all")
            check_prerequisites
            check_database
            check_backend
            check_frontend
            check_integration
            run_quick_tests
            show_summary
            ;;
        "db"|"database")
            check_database
            ;;
        "backend"|"api")
            check_backend
            ;;
        "frontend"|"ui")
            check_frontend
            ;;
        "detailed"|"system")
            detailed_system_info
            ;;
        "quick"|"test")
            run_quick_tests
            ;;
        "summary"|"status")
            show_summary
            ;;
        *)
            print_error "Unknown option: $option"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"