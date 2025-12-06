#!/bin/sh
set -e

CONFIG_FILE="/usr/share/nginx/html/env-config.js"

# Generate runtime config
cat <<EOF > "$CONFIG_FILE"
window.__ENV__ = {
  VITE_API_URL: "${VITE_API_URL:-http://localhost:8080}",
  VITE_USE_MOCK_API: "${VITE_USE_MOCK_API:-false}"
};
EOF

echo "Runtime config generated:"
cat "$CONFIG_FILE"

# Execute the main command (nginx)
exec "$@"