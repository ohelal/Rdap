#!/bin/bash

set -euo pipefail

# Configuration
RDAP_DIR="config"
IANA_BASE_URL="https://data.iana.org/rdap"
LOG_FILE="config/rdap-updates.log"
BACKUP_DIR="config/backup/$(date +%Y%m%d)"
MAX_BACKUPS=7
VERIFY_SSL=true

# Command line options
FORCE=false
SKIP_VERIFY=false
SKIP_BACKUP=false
SKIP_SERVICE_RELOAD=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check dependencies
check_dependencies() {
    local missing_deps=()
    for cmd in curl jq dirname basename mktemp; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_deps+=("$cmd")
        fi
    done
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}Error: Missing required dependencies: ${missing_deps[*]}${NC}"
        exit 1
    fi
}

# Function to print usage
usage() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS]

Options:
    -f, --force              Force update even if files are current
    -s, --skip-verify       Skip SSL verification
    -b, --skip-backup       Skip creating backups
    -n, --no-service-reload Skip service reload
    -h, --help              Show this help message
EOF
}

# Function to parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--force)
                FORCE=true
                shift
                ;;
            -s|--skip-verify)
                SKIP_VERIFY=true
                VERIFY_SSL=false
                shift
                ;;
            -b|--skip-backup)
                SKIP_BACKUP=true
                shift
                ;;
            -n|--no-service-reload)
                SKIP_SERVICE_RELOAD=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                usage
                exit 1
                ;;
        esac
    done
}

# Function to log messages
log_message() {
    local level=$1
    local message=$2
    local color=$NC
    
    case $level in
        INFO)  color=$GREEN ;;
        WARN)  color=$YELLOW ;;
        ERROR) color=$RED ;;
    esac
    
    echo -e "${color}[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message${NC}" | tee -a "$LOG_FILE"
}

# Function to clean old backups
clean_old_backups() {
    if [ "$SKIP_BACKUP" = true ]; then
        return
    fi
    
    local backup_count=$(find "$(dirname "$BACKUP_DIR")" -maxdepth 1 -type d | wc -l)
    if [ "$backup_count" -gt "$MAX_BACKUPS" ]; then
        log_message "INFO" "Cleaning old backups..."
        find "$(dirname "$BACKUP_DIR")" -maxdepth 1 -type d -printf '%T@ %p\n' | \
            sort -n | head -n -"$MAX_BACKUPS" | cut -d' ' -f2- | \
            xargs rm -rf
    fi
}

# Function to verify file contents
verify_file() {
    local file=$1
    
    # Check file size
    if [ ! -s "$file" ]; then
        return 1
    fi
    
    # Verify JSON structure
    if ! jq empty "$file" >/dev/null 2>&1; then
        return 1
    fi
    
    # Additional checks could be added here
    return 0
}

# Function to download and verify a file
download_file() {
    local file=$1
    local temp_file
    temp_file=$(mktemp)
    local final_file="${RDAP_DIR}/${file}"
    local curl_opts=("-s" "-o" "$temp_file")
    
    if [ "$VERIFY_SSL" = false ]; then
        curl_opts+=("-k")
    fi
    
    log_message "INFO" "Downloading ${file}..."
    
    if ! curl "${curl_opts[@]}" "${IANA_BASE_URL}/${file}"; then
        log_message "ERROR" "Failed to download ${file}"
        rm -f "$temp_file"
        return 1
    fi
    
    if ! verify_file "$temp_file"; then
        log_message "ERROR" "Invalid or empty file: ${file}"
        rm -f "$temp_file"
        return 1
    fi
    
    if [ "$SKIP_BACKUP" = false ] && [ -f "$final_file" ]; then
        mkdir -p "$BACKUP_DIR"
        cp "$final_file" "${BACKUP_DIR}/${file}.bak"
    fi
    
    mv "$temp_file" "$final_file"
    log_message "INFO" "Successfully updated ${file}"
    return 0
}

# Function to reload service
reload_service() {
    if [ "$SKIP_SERVICE_RELOAD" = true ]; then
        return
    fi
    
    local pid_file="/var/run/rdap-service.pid"
    if [ -f "$pid_file" ]; then
        log_message "INFO" "Notifying RDAP service to reload configuration..."
        if kill -SIGHUP "$(cat "$pid_file")" 2>/dev/null; then
            log_message "INFO" "Service reload signal sent successfully"
        else
            log_message "WARN" "Failed to send reload signal to service"
        fi
    else
        log_message "WARN" "PID file not found, skipping service reload"
    fi
}

main() {
    check_dependencies
    parse_args "$@"
    
    # Create required directories
    mkdir -p "$RDAP_DIR"
    
    # Files to download
    local -a FILES=(
        "asn.json"
        "dns.json"
        "ipv4.json"
        "ipv6.json"
        "object-tags.json"
    )
    
    local update_count=0
    for file in "${FILES[@]}"; do
        if download_file "$file"; then
            ((update_count++))
        fi
    done
    
    clean_old_backups
    
    if [ $update_count -gt 0 ]; then
        reload_service
        log_message "INFO" "Updated $update_count files successfully"
    else
        log_message "INFO" "No files were updated"
    fi
}

# Run main function
main "$@"