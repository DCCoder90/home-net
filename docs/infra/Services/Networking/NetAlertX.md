# NetAlertX

## Router Config

/config/scripts/exportLeases.sh:

```bash
#!/bin/bash

# --- Configuration ---
SMB_SERVER="192.168.1.41"
SMB_SHARE_NAME="Misc"       # The name of the SMB share
SMB_REMOTE_DIR="dhcpExports" # The directory path within the SMB share
SMB_USER="username" # Replace with your SMB username
SMB_PASS="password" # Replace with your SMB password

# Common DHCP lease file locations on EdgeOS or other Linux systems
declare -a LEASE_FILE_PATHS=(
    "/var/run/dnsmasq.leases"          # Common for dnsmasq
    "/var/lib/dnsmasq/dnsmasq.leases"  # Another common dnsmasq path
    "/config/dhcpd.leases"             # EdgeOS often keeps ISC dhcpd leases in /config (persists across reboots)
    "/var/lib/dhcp/dhcpd.leases"       # Common for ISC dhcpd on Debian-based systems
    "/var/run/dhcpd.leases"            # Older path for ISC dhcpd
)

# Name of the file as it will be saved on the SMB share.
EXPORT_FILENAME="dhcp_leases_$(hostname).txt"

# --- Helper Functions ---
log_message() {
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $1"
    # echo "$(date +'%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
}

# --- Prerequisite Check ---
if ! command -v smbclient &> /dev/null; then
    log_message "Error: smbclient command not found. Please install it first."
    log_message "On EdgeOS, you might need to run: sudo apt-get update && sudo apt-get install smbclient"
    exit 1
fi

# --- Find DHCP Lease File ---
SOURCE_LEASE_FILE=""
for path_to_check in "${LEASE_FILE_PATHS[@]}"; do
    if [ -f "$path_to_check" ] && [ -s "$path_to_check" ]; then # Check if file exists AND is not empty
        SOURCE_LEASE_FILE="$path_to_check"
        log_message "Found DHCP lease file: $SOURCE_LEASE_FILE"
        break
    fi
done

if [ -z "$SOURCE_LEASE_FILE" ]; then
    log_message "Error: DHCP lease file not found at any specified common locations or it is empty."
    exit 1
fi

# --- Perform SMB Upload ---

FULL_REMOTE_TARGET="//${SMB_SERVER}/${SMB_SHARE_NAME}/${SMB_REMOTE_DIR}/${EXPORT_FILENAME}"

log_message "Attempting to upload '$SOURCE_LEASE_FILE' to '$FULL_REMOTE_TARGET'..."

smbclient "//${SMB_SERVER}/${SMB_SHARE_NAME}" -U "${SMB_USER}%${SMB_PASS}" \
    --socket-options='TCP_NODELAY SO_RCVBUF=8192 SO_SNDBUF=8192' \
    -m SMB3 \
    -c "cd \"${SMB_REMOTE_DIR}\"; put \"${SOURCE_LEASE_FILE}\" \"${EXPORT_FILENAME}\""

SMB_EXIT_CODE=$?

if [ $SMB_EXIT_CODE -eq 0 ]; then
    log_message "Successfully uploaded DHCP lease file to '$FULL_REMOTE_TARGET'"
else
    log_message "Error: Failed to upload DHCP lease file. smbclient exited with code $SMB_EXIT_CODE."
    log_message "Check SMB server, share name, remote path, credentials, and if the remote directory exists."
    exit 1
fi

log_message "Script finished successfully."
exit 0
```

Followed by running the following commands which will run the script hourly to export DHCP leases to Unraid:

```
configure;
set system task-scheduler task dhcp-export executable path /config/scripts/exportLeases.sh
set system task-scheduler task dhcp-export interval 1h
commit; save; exit;
```