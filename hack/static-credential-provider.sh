#!/bin/bash
set -e

# Check if jq is installed
if ! command -v jq &>/dev/null; then
    echo "Error: jq is not installed. Please install jq to proceed." >&2
    exit 1
fi

if [ -z ${KSCP_REGISTRY_USERNAME+x} ] || [ -z ${KSCP_REGISTRY_PASSWORD+x} ]; then
    echo "Username and password are required."
    exit 1
fi

# Read the input string from stdin
read -r input

# Attempt to parse the input with jq
if ! echo "$input" | jq . &>/dev/null; then
    echo "Error: Input is not valid JSON." >&2
    exit 1
fi

# Create the response
IMAGE=$(echo "$input" | jq -r .image | sed -E 's#(/[^:]*):.*#\1#g')
CREDENTIAL_RESPONSE=$(
    cat <<EOF
{
  "kind": "CredentialProviderResponse",
  "apiVersion": "credentialprovider.kubelet.k8s.io/v1",
  "cacheKeyType": "Image",
  "cacheDuration": "${KSCP_CACHE_DURATION-8h0m0s}",
  "auth": {
    "${IMAGE}": {
      "username": "$KSCP_REGISTRY_USERNAME",
      "password": "$KSCP_REGISTRY_USERNAME"
    }
  }
}
EOF
)
echo $CREDENTIAL_RESPONSE | jq -c
