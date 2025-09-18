#!/bin/bash
#
# GitHub Actions Workflow Status Checker
# Usage: ./scripts/check-workflows.sh
#

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'  
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) is not installed"
        log_info "Install it from: https://cli.github.com/"
        exit 1
    fi
    
    # Check if authenticated
    if ! gh auth status &> /dev/null; then
        log_error "Not authenticated with GitHub CLI"
        log_info "Run: gh auth login"
        exit 1
    fi
}

show_workflow_status() {
    local repo=$1
    
    log_info "=== GitHub Actions Workflow Status ==="
    echo
    
    # Get latest workflow runs
    log_info "Recent workflow runs:"
    gh run list --repo "$repo" --limit 10 --json status,name,conclusion,createdAt,headBranch | \
        jq -r '.[] | "\(.status) | \(.name) | \(.conclusion // "running") | \(.headBranch) | \(.createdAt)"' | \
        while IFS='|' read -r status name conclusion branch created; do
            status=$(echo "$status" | xargs)
            name=$(echo "$name" | xargs)
            conclusion=$(echo "$conclusion" | xargs)
            branch=$(echo "$branch" | xargs)
            created=$(echo "$created" | xargs)
            
            case "$conclusion" in
                "success") 
                    echo -e "  ${GREEN}✅${NC} $name (${branch}) - $(date -d "$created" '+%Y-%m-%d %H:%M')"
                    ;;
                "failure")
                    echo -e "  ${RED}❌${NC} $name (${branch}) - $(date -d "$created" '+%Y-%m-%d %H:%M')"
                    ;;
                "running")
                    echo -e "  ${YELLOW}🔄${NC} $name (${branch}) - $(date -d "$created" '+%Y-%m-%d %H:%M')"
                    ;;
                *)
                    echo -e "  ${YELLOW}❓${NC} $name (${branch}) - $(date -d "$created" '+%Y-%m-%d %H:%M')"
                    ;;
            esac
        done
    
    echo
    log_info "Workflow files:"
    find .github/workflows -name "*.yml" -o -name "*.yaml" | while read -r file; do
        workflow_name=$(grep -m1 "^name:" "$file" | sed 's/name: *//' | tr -d '"'"'"'')
        echo "  📄 $file - $workflow_name"
    done
    
    echo
    log_info "GitHub Secrets required:"
    echo "  🔐 DOCKER_USERNAME - Docker Hub username"
    echo "  🔐 DOCKER_PASSWORD - Docker Hub access token"
    echo "  📝 GITHUB_TOKEN - Automatically provided by GitHub"
    
    echo
    log_info "Useful commands:"
    echo "  📊 View all runs:     gh run list --repo $repo"
    echo "  📋 View specific run: gh run view [RUN_ID] --repo $repo"
    echo "  🔄 Re-run failed:     gh run rerun [RUN_ID] --repo $repo"
    echo "  🏷️  List releases:     gh release list --repo $repo"
}

show_help() {
    cat << EOF
GitHub Actions Workflow Status Checker

USAGE:
    ./scripts/check-workflows.sh [REPO]

ARGUMENTS:
    REPO    Repository in format owner/repo (optional, auto-detected from git)

EXAMPLES:
    ./scripts/check-workflows.sh
    ./scripts/check-workflows.sh stalukder-plume/openusp

REQUIREMENTS:
    - GitHub CLI (gh) installed and authenticated
    - Access to the repository

This script shows:
    - Recent workflow runs and their status
    - Available workflow files
    - Required GitHub secrets
    - Useful commands for managing workflows
EOF
}

# Parse arguments
case "${1:-}" in
    --help|-h)
        show_help
        exit 0
        ;;
esac

# Check prerequisites
check_gh_cli

# Get repository
if [[ -n "${1:-}" ]]; then
    REPO="$1"
else
    # Auto-detect from git remote
    if git remote get-url origin &> /dev/null; then
        ORIGIN=$(git remote get-url origin)
        # Extract owner/repo from various URL formats
        if [[ $ORIGIN == git@github.com:* ]]; then
            REPO=$(echo "$ORIGIN" | sed 's/git@github.com://' | sed 's/\.git$//')
        elif [[ $ORIGIN == https://github.com/* ]]; then
            REPO=$(echo "$ORIGIN" | sed 's|https://github.com/||' | sed 's/\.git$//')
        else
            log_error "Could not parse repository from git remote: $ORIGIN"
            exit 1
        fi
    else
        log_error "Not in a git repository and no REPO specified"
        log_info "Usage: $0 [owner/repo]"
        exit 1
    fi
fi

log_info "Checking workflows for repository: $REPO"
echo

# Show status
show_workflow_status "$REPO"

log_success "Workflow status check completed!"