#!/bin/bash

# Setup Frontend Static Analysis Tools for normal-form-app
# Usage: ./scripts/setup-frontend-tools.sh

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

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js first."
        exit 1
    fi
    
    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm first."
        exit 1
    fi
    
    print_status "Node.js version: $(node --version)"
    print_status "npm version: $(npm --version)"
    echo ""
}

setup_git_hooks() {
    print_status "Setting up Git hooks with Husky..."
    
    cd "$FRONTEND_DIR"
    
    # Initialize husky if not already done
    if [ ! -d ".husky/_" ]; then
        npx husky init
    fi
    
    # Create pre-commit hook if it doesn't exist
    if [ ! -f ".husky/pre-commit" ]; then
        echo "#!/bin/sh" > .husky/pre-commit
        echo ". \"\$(dirname \"\$0\")/_/husky.sh\"" >> .husky/pre-commit
        echo "" >> .husky/pre-commit
        echo "cd frontend && npx lint-staged" >> .husky/pre-commit
        chmod +x .husky/pre-commit
    fi
    
    print_success "Git hooks configured"
}

install_additional_tools() {
    print_status "Installing additional development tools..."
    
    cd "$FRONTEND_DIR"
    
    # Install Playwright for E2E testing (optional)
    print_status "Installing Playwright for E2E testing..."
    npm install --save-dev @playwright/test
    
    # Install testing utilities
    print_status "Installing testing utilities..."
    npm install --save-dev @testing-library/react @testing-library/jest-dom @testing-library/user-event vitest jsdom
    
    print_success "Additional tools installed"
}

create_playwright_config() {
    print_status "Creating Playwright configuration..."
    
    cat > "$FRONTEND_DIR/playwright.config.ts" << 'EOF'
import { defineConfig, devices } from '@playwright/test';

/**
 * @see https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  testDir: './e2e',
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: 'html',
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: 'http://localhost:5173',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },

    /* Test against mobile viewports. */
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },

    /* Test against branded browsers. */
    // {
    //   name: 'Microsoft Edge',
    //   use: { ...devices['Desktop Edge'], channel: 'msedge' },
    // },
    // {
    //   name: 'Google Chrome',
    //   use: { ...devices['Desktop Chrome'], channel: 'chrome' },
    // },
  ],

  /* Run your local dev server before starting the tests */
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
EOF
    
    print_success "Playwright configuration created"
}

create_vitest_config() {
    print_status "Creating Vitest configuration..."
    
    cat > "$FRONTEND_DIR/vitest.config.ts" << 'EOF'
/// <reference types="vitest" />
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/components': path.resolve(__dirname, './src/components'),
      '@/pages': path.resolve(__dirname, './src/pages'),
      '@/hooks': path.resolve(__dirname, './src/hooks'),
      '@/services': path.resolve(__dirname, './src/services'),
      '@/types': path.resolve(__dirname, './src/types'),
      '@/utils': path.resolve(__dirname, './src/utils'),
      '@/constants': path.resolve(__dirname, './src/constants'),
    },
  },
});
EOF
    
    # Create test setup file
    mkdir -p "$FRONTEND_DIR/src/test"
    cat > "$FRONTEND_DIR/src/test/setup.ts" << 'EOF'
import '@testing-library/jest-dom';
EOF
    
    print_success "Vitest configuration created"
}

update_package_scripts() {
    print_status "Updating package.json scripts..."
    
    # This would be done manually or with a JSON parser
    # For now, we'll just show what should be added
    print_status "Please ensure the following scripts are in package.json:"
    echo "  \"test\": \"vitest\","
    echo "  \"test:ui\": \"vitest --ui\","
    echo "  \"test:coverage\": \"vitest --coverage\","
    echo "  \"e2e\": \"playwright test\","
    echo "  \"e2e:ui\": \"playwright test --ui\","
    echo "  \"e2e:debug\": \"playwright test --debug\""
}

verify_setup() {
    print_status "Verifying setup..."
    
    cd "$FRONTEND_DIR"
    
    # Test TypeScript
    print_status "Testing TypeScript compilation..."
    if npm run type-check; then
        print_success "TypeScript check passed"
    else
        print_error "TypeScript check failed"
    fi
    
    # Test ESLint
    print_status "Testing ESLint..."
    if npm run lint; then
        print_success "ESLint check passed"
    else
        print_error "ESLint check failed"
    fi
    
    # Test Prettier
    print_status "Testing Prettier..."
    if npm run format:check; then
        print_success "Prettier check passed"
    else
        print_warning "Some files need formatting"
    fi
}

main() {
    print_status "Setting up Frontend Static Analysis Tools..."
    echo ""
    
    check_prerequisites
    setup_git_hooks
    install_additional_tools
    create_playwright_config
    create_vitest_config
    update_package_scripts
    verify_setup
    
    echo ""
    print_success "Frontend static analysis tools setup complete!"
    print_status "Available commands:"
    echo "  npm run lint           - Run ESLint"
    echo "  npm run lint:fix       - Fix ESLint issues"
    echo "  npm run format         - Format code with Prettier"
    echo "  npm run format:check   - Check code formatting"
    echo "  npm run type-check     - Run TypeScript check"
    echo "  npm run check-all      - Run all checks"
    echo "  npm run test           - Run unit tests"
    echo "  npm run e2e            - Run E2E tests"
    echo ""
}

main "$@"