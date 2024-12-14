#!/bin/bash

# Configuration
RDAP_DIR="config"
IANA_BASE_URL="https://data.iana.org/rdap"
LOG_FILE="config/rdap-updates.log"
BACKUP_DIR="config/backup/$(date +%Y%m%d)"

# Create directories
mkdir -p "$RDAP_DIR" "$BACKUP_DIR"

# Function to log messages
log_message() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Function to download and verify a file
download_file() {
    local file=$1
    local temp_file="${RDAP_DIR}/${file}.tmp"
    local final_file="${RDAP_DIR}/${file}"
    
    log_message "Downloading ${file}..."
    
    # Download to temp file
    if curl -s -o "$temp_file" "${IANA_BASE_URL}/${file}"; then
        # Verify JSON syntax
        if jq empty "$temp_file" 2>/dev/null; then
            # Backup existing file if it exists
            if [ -f "$final_file" ]; then
                cp "$final_file" "${BACKUP_DIR}/${file}.bak"
            fi
            # Move temp file to final location
            mv "$temp_file" "$final_file"
            log_message "Successfully updated ${file}"
            return 0
        else
            log_message "ERROR: Invalid JSON in ${file}"
            rm -f "$temp_file"
            return 1
        fi
    else
        log_message "ERROR: Failed to download ${file}"
        rm -f "$temp_file"
        return 1
    fi
}

# List of files to download
FILES=(
    "asn.json"
    "dns.json"
    "ipv4.json"
    "ipv6.json"
    "object-tags.json"
)

# Download all files
for file in "${FILES[@]}"; do
    download_file "$file"
done
