#!/bin/bash
#
# OpenUSP Docker Release Script with Semantic Versioning
# Usage: ./scripts/docker-release.sh [version] [--push] [--component=name]
#
# Examples:
#   ./scripts/docker-release.sh v1.2.3 --push              # Release all components
#   ./scripts/docker-release.sh v1.2.3 --component=controller # Build single component
#   ./scripts/docker-release.sh --help                     # Show help
#

set -e  # Exit on any error

# Configuration
REPO_PREFIX="stalukder-plume/openusp"
COMPONENTS=("controller" "apiserver" "cli" "cwmpacs")
COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git branch --show-current)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
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

show_help() {
    cat << EOF
OpenUSP Docker Release Script

USAGE:
    ./scripts/docker-release.sh [VERSION] [OPTIONS]

VERSION:
    Semantic version (e.g., v1.2.3, v2.0.0-beta.1)
    If not provided, will prompt for input

OPTIONS:
    --push              Push images to Docker Hub after building
    --component=NAME    Build only specified component (controller, apiserver, cli, cwmpacs)
    --dry-run          Show what would be done without executing
    --help             Show this help message

EXAMPLES:
    ./scripts/docker-release.sh v1.2.3 --push
    ./scripts/docker-release.sh v1.2.3 --component=controller
    ./scripts/docker-release.sh v2.0.0-rc.1 --dry-run

SEMANTIC VERSIONING:
    MAJOR.MINOR.PATCH (e.g., 1.2.3)
    - MAJOR: Breaking changes
    - MINOR: New features (backwards compatible)
    - PATCH: Bug fixes

    Pre-release: v1.2.3-alpha.1, v1.2.3-beta.2, v1.2.3-rc.1
EOF
}

validate_version() {
    local version=$1
    # Check if version follows semantic versioning pattern
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        log_error "Invalid version format: $version"
        log_info "Expected format: vMAJOR.MINOR.PATCH (e.g., v1.2.3)"
        log_info "Pre-release: v1.2.3-alpha.1, v1.2.3-beta.1, v1.2.3-rc.1"
        exit 1
    fi
}

check_git_status() {
    if [[ -n $(git status --porcelain) ]]; then
        log_warning "Working directory is not clean. Uncommitted changes detected."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Aborted by user"
            exit 1
        fi
    fi
}

build_component() {
    local component=$1
    local version=$2
    local dockerfile="build/$component/Dockerfile"
    
    if [[ ! -f $dockerfile ]]; then
        log_error "Dockerfile not found: $dockerfile"
        return 1
    fi
    
    log_info "Building $component..."
    
    local image_name="$REPO_PREFIX-$component"
    local tags=(
        "$image_name:$version"
        "$image_name:$COMMIT"
    )
    
    # Add 'latest' tag only for stable releases (no pre-release suffix)
    if [[ ! $version =~ -[a-zA-Z] ]]; then
        tags+=("$image_name:latest")
    fi
    
    # Build command with all tags
    local build_cmd="docker build -f $dockerfile"
    for tag in "${tags[@]}"; do
        build_cmd="$build_cmd -t $tag"
    done
    build_cmd="$build_cmd ."
    
    if [[ $DRY_RUN == true ]]; then
        log_info "DRY RUN: $build_cmd"
    else
        eval $build_cmd
        log_success "Built $component with tags: ${tags[*]}"
    fi
    
    # Store tags for pushing
    BUILT_TAGS+=("${tags[@]}")
}

push_images() {
    if [[ ${#BUILT_TAGS[@]} -eq 0 ]]; then
        log_warning "No images to push"
        return
    fi
    
    log_info "Pushing ${#BUILT_TAGS[@]} image tags to Docker Hub..."
    
    for tag in "${BUILT_TAGS[@]}"; do
        if [[ $DRY_RUN == true ]]; then
            log_info "DRY RUN: docker push $tag"
        else
            log_info "Pushing $tag..."
            docker push "$tag"
            log_success "Pushed $tag"
        fi
    done
}

show_summary() {
    local version=$1
    echo
    log_info "=== RELEASE SUMMARY ==="
    echo "Version: $version"
    echo "Commit: $COMMIT"
    echo "Branch: $BRANCH"
    echo "Components: ${TARGET_COMPONENTS[*]}"
    echo "Total tags built: ${#BUILT_TAGS[@]}"
    
    if [[ $PUSH_IMAGES == true ]]; then
        echo "Images pushed: ✅"
    else
        echo "Images pushed: ❌ (use --push to push)"
    fi
    
    echo
    log_info "To pull these images:"
    for component in "${TARGET_COMPONENTS[@]}"; do
        echo "  docker pull $REPO_PREFIX-$component:$version"
    done
}

# Parse arguments
VERSION=""
PUSH_IMAGES=false
SINGLE_COMPONENT=""
DRY_RUN=false
BUILT_TAGS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            show_help
            exit 0
            ;;
        --push)
            PUSH_IMAGES=true
            shift
            ;;
        --component=*)
            SINGLE_COMPONENT="${1#*=}"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        v*.*)
            VERSION=$1
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Get version if not provided
if [[ -z $VERSION ]]; then
    echo -e "${BLUE}Current git info:${NC}"
    echo "  Branch: $BRANCH"
    echo "  Commit: $COMMIT"
    echo
    read -p "Enter semantic version (e.g., v1.2.3): " VERSION
fi

# Validate inputs
validate_version "$VERSION"
check_git_status

# Determine target components
if [[ -n $SINGLE_COMPONENT ]]; then
    if [[ ! " ${COMPONENTS[@]} " =~ " ${SINGLE_COMPONENT} " ]]; then
        log_error "Unknown component: $SINGLE_COMPONENT"
        log_info "Available components: ${COMPONENTS[*]}"
        exit 1
    fi
    TARGET_COMPONENTS=("$SINGLE_COMPONENT")
else
    TARGET_COMPONENTS=("${COMPONENTS[@]}")
fi

# Show what we're about to do
echo
log_info "=== BUILD PLAN ==="
echo "Version: $VERSION"
echo "Components: ${TARGET_COMPONENTS[*]}"
echo "Push to registry: $(if [[ $PUSH_IMAGES == true ]]; then echo "YES"; else echo "NO"; fi)"
echo "Dry run: $(if [[ $DRY_RUN == true ]]; then echo "YES"; else echo "NO"; fi)"
echo

if [[ $DRY_RUN == false ]]; then
    read -p "Continue? (Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        log_info "Aborted by user"
        exit 0
    fi
fi

# Build images
log_info "Starting build process..."
for component in "${TARGET_COMPONENTS[@]}"; do
    build_component "$component" "$VERSION"
done

# Push if requested
if [[ $PUSH_IMAGES == true ]]; then
    push_images
fi

# Show summary
show_summary "$VERSION"

log_success "Release process completed! 🚀"