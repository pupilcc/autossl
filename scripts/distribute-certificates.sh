#!/bin/bash
set -u

# ============================================================
# Configuration
# ============================================================

SSL_GROUP="ssl-cert"

declare -A urls_and_paths=(
    ["https://ssl.example.com/dl/REPLACE_WITH_CERTIFICATE_CODE.crt"]="/etc/ssl/example.com/fullchain.pem"
    ["https://ssl.example.com/dl/REPLACE_WITH_PRIVATE_KEY_CODE.key"]="/etc/ssl/example.com/privkey.pem"
)

# Track whether any certificate file changed.
execute_operation=false


# ============================================================
# Check dependencies
# ============================================================

for cmd in curl sha256sum getent groupadd systemctl mktemp; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "Error: required command not found: $cmd"
        exit 1
    fi
done


# ============================================================
# Create the ssl-cert group
# ============================================================

if ! getent group "$SSL_GROUP" >/dev/null 2>&1; then
    echo "Creating system group: $SSL_GROUP"

    if ! groupadd --system "$SSL_GROUP"; then
        echo "Error: failed to create group: $SSL_GROUP"
        exit 1
    fi
fi


# ============================================================
# Process all certificate files
# ============================================================

for url in "${!urls_and_paths[@]}"; do

    local_path="${urls_and_paths[$url]}"
    target_dir="$(dirname "$local_path")"
    filename="$(basename "$local_path")"

    echo
    echo "============================================================"
    echo "File: $local_path"
    echo "============================================================"

    # --------------------------------------------------------
    # Create the target directory.
    #
    # Permissions: root:ssl-cert 750
    # --------------------------------------------------------

    if ! mkdir -p "$target_dir"; then
        echo "Error: failed to create directory: $target_dir"
        continue
    fi

    chown root:"$SSL_GROUP" "$target_dir"
    chmod 750 "$target_dir"


    # --------------------------------------------------------
    # Create the temporary file in the target directory so the
    # final mv stays on the same filesystem and is atomic.
    # --------------------------------------------------------

    temp_file=$(mktemp "${target_dir}/.${filename}.tmp.XXXXXX")

    if [ -z "$temp_file" ]; then
        echo "Error: failed to create temporary file"
        continue
    fi

    chmod 600 "$temp_file"


    # --------------------------------------------------------
    # Download
    # --------------------------------------------------------

    if ! curl \
        -fsS \
        --retry 3 \
        --retry-delay 2 \
        --connect-timeout 10 \
        -o "$temp_file" \
        "$url"; then

        echo "Warning: failed to download: $local_path"

        rm -f "$temp_file"
        continue
    fi


    # --------------------------------------------------------
    # Reject empty downloads
    # --------------------------------------------------------

    if [ ! -s "$temp_file" ]; then
        echo "Warning: downloaded file is empty: $local_path"

        rm -f "$temp_file"
        continue
    fi


    # --------------------------------------------------------
    # Perform basic validation with OpenSSL when available
    # --------------------------------------------------------

    case "$local_path" in

        *privkey.pem)
            if command -v openssl >/dev/null 2>&1; then
                if ! openssl pkey \
                    -in "$temp_file" \
                    -noout \
                    >/dev/null 2>&1; then

                    echo "Warning: invalid private key: $local_path"

                    rm -f "$temp_file"
                    continue
                fi
            fi
            ;;

        *fullchain.pem)
            if command -v openssl >/dev/null 2>&1; then
                if ! openssl x509 \
                    -in "$temp_file" \
                    -noout \
                    >/dev/null 2>&1; then

                    echo "Warning: invalid certificate: $local_path"

                    rm -f "$temp_file"
                    continue
                fi
            fi
            ;;

        *)
            echo "Warning: unknown certificate file type: $local_path"

            rm -f "$temp_file"
            continue
            ;;

    esac


    # --------------------------------------------------------
    # Compare SHA-256 hashes
    # --------------------------------------------------------

    remote_hash=$(sha256sum "$temp_file" | awk '{print $1}')

    if [ -f "$local_path" ]; then
        local_hash=$(sha256sum "$local_path" | awk '{print $1}')
    else
        local_hash=""
    fi

    echo "Remote SHA256: $remote_hash"
    echo "Local  SHA256: $local_hash"


    # --------------------------------------------------------
    # The file is unchanged. Remove the temporary file and
    # enforce the expected ownership and permissions.
    # --------------------------------------------------------

    if [ "$remote_hash" = "$local_hash" ]; then

        echo "No change: $local_path"

        rm -f "$temp_file"

        case "$local_path" in

            *privkey.pem)
                chown root:"$SSL_GROUP" "$local_path"
                chmod 640 "$local_path"
                ;;

            *fullchain.pem)
                chown root:root "$local_path"
                chmod 644 "$local_path"
                ;;

        esac

        continue
    fi


    # --------------------------------------------------------
    # The file changed
    # --------------------------------------------------------

    echo "Updating: $local_path"


    # --------------------------------------------------------
    # Set ownership and permissions before replacement
    # --------------------------------------------------------

    case "$local_path" in

        *privkey.pem)
            chown root:"$SSL_GROUP" "$temp_file"
            chmod 640 "$temp_file"
            ;;

        *fullchain.pem)
            chown root:root "$temp_file"
            chmod 644 "$temp_file"
            ;;

    esac


    # --------------------------------------------------------
    # Atomically replace the target file
    # --------------------------------------------------------

    if ! mv -f "$temp_file" "$local_path"; then

        echo "Error: failed to replace: $local_path"

        rm -f "$temp_file"
        continue
    fi


    # --------------------------------------------------------
    # Enforce final ownership and permissions
    # --------------------------------------------------------

    case "$local_path" in

        *privkey.pem)
            chown root:"$SSL_GROUP" "$local_path"
            chmod 640 "$local_path"
            ;;

        *fullchain.pem)
            chown root:root "$local_path"
            chmod 644 "$local_path"
            ;;

    esac

    execute_operation=true

done


# ============================================================
# Reload Nginx when a certificate changed
# ============================================================

if [ "$execute_operation" = true ]; then

    echo
    echo "============================================================"
    echo "Certificate changed."
    echo "============================================================"

    if ! command -v nginx >/dev/null 2>&1; then
        echo "Warning: nginx command not found"
        exit 0
    fi

    echo "Testing nginx configuration..."

    if nginx -t; then

        echo "Reloading nginx..."

        if systemctl reload nginx; then
            echo "Nginx reloaded successfully."
        else
            echo "Error: failed to reload nginx"
            exit 1
        fi

    else

        echo "Error: nginx configuration test failed."
        echo "Nginx was NOT reloaded."

        exit 1
    fi

else

    echo
    echo "No certificate changes detected."

fi


echo
echo "Certificate update process completed."
