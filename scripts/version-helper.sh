#!/bin/bash
#
# Version Helper Script - Suggests next semantic version based on git history
# Usage: ./scripts/version-helper.sh [--suggest] [--list]
#

set -e

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

get_latest_version() {
    # Get latest git tag that follows semver pattern
    git tag --list --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' | head -n1
}

parse_version() {
    local version=$1
    # Remove 'v' prefix and split into components
    version=${version#v}
    IFS='.' read -r major minor patch <<< "$version"
    echo "$major $minor $patch"
}

suggest_next_version() {
    local latest=$(get_latest_version)
    
    if [[ -z $latest ]]; then
        log_warning "No previous version tags found"
        echo "Suggested first version: v1.0.0"
        return
    fi
    
    log_info "Latest version: $latest"
    
    # Parse current version
    read -r major minor patch <<< $(parse_version "$latest")
    
    echo
    log_info "Suggested next versions:"
    echo "  🔧 Patch (bug fixes):     v$major.$minor.$((patch + 1))"
    echo "  ✨ Minor (new features):  v$major.$((minor + 1)).0"
    echo "  💥 Major (breaking):      v$((major + 1)).0.0"
    echo
    
    # Check recent commits for hints
    local commits_since_tag=$(git rev-list ${latest}..HEAD --count 2>/dev/null || echo "0")
    if [[ $commits_since_tag -gt 0 ]]; then
        echo "📝 $commits_since_tag commits since $latest:"
        git log --oneline ${latest}..HEAD | head -5
        
        # Simple heuristics
        local commit_messages=$(git log --oneline ${latest}..HEAD --format="%s")
        
        if echo "$commit_messages" | grep -qi "break\|major\|BREAKING"; then
            log_warning "Detected potential BREAKING CHANGES - consider major version bump"
        elif echo "$commit_messages" | grep -qi "feat\|feature\|add"; then
            log_info "Detected new features - consider minor version bump"
        else
            log_info "Looks like bug fixes - consider patch version bump"
        fi
    fi
}

list_versions() {
    log_info "All version tags (latest first):"
    git tag --list --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' | head -10
    
    local total_tags=$(git tag --list | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' | wc -l)
    if [[ $total_tags -gt 10 ]]; then
        echo "... and $((total_tags - 10)) more"
    fi
}

show_help() {
    cat << EOF
Version Helper Script

USAGE:
    ./scripts/version-helper.sh [OPTIONS]

OPTIONS:
    --suggest    Suggest next version based on git history (default)
    --list       List all existing version tags
    --help       Show this help

EXAMPLES:
    ./scripts/version-helper.sh              # Suggest next version
    ./scripts/version-helper.sh --list       # List all versions
    
This script helps you choose the next semantic version by:
- Showing the latest version
- Analyzing recent commits
- Suggesting patch/minor/major increments
EOF
}

# Parse arguments
case "${1:---suggest}" in
    --help|-h)
        show_help
        ;;
    --list)
        list_versions
        ;;
    --suggest|"")
        suggest_next_version
        ;;
    *)
        echo "Unknown option: $1"
        show_help
        exit 1
        ;;
esac