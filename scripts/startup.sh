#!/usr/bin/env bash
set -euo pipefail

APP_PORT="${APP_PORT:-3000}"
SITE_NAME="${SITE_NAME:-server}"

SERVICE_NAME="${SERVICE_NAME:-server}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/server}"
APP_BINARY="${APP_BINARY:-server}"
ENV_PORT="${ENV_PORT:-3000}"
OVERWRITE_UNIT="${OVERWRITE_UNIT:-false}"

UNIT_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
NGINX_AVAIL="/etc/nginx/sites-available/${SITE_NAME}"
NGINX_ENABLED="/etc/nginx/sites-enabled/${SITE_NAME}"

echo "==> Checking OS (Debian/Ubuntu required)"
if ! command -v apt-get >/dev/null 2>&1; then
  echo "ERROR: This script supports Debian/Ubuntu only."
  exit 1
fi

echo "==> Installing nginx"
sudo apt-get update -y
sudo apt-get install -y nginx

echo "==> Enabling nginx"
sudo systemctl enable nginx
sudo systemctl restart nginx

echo "==> Ensuring deploy directory: ${DEPLOY_PATH}"
sudo mkdir -p "${DEPLOY_PATH}"
sudo chown -R "$(whoami)":"$(whoami)" "${DEPLOY_PATH}" || true

echo "==> Writing nginx config"
sudo bash -c "cat > '${NGINX_AVAIL}'" <<EOF
server {
  listen 80;
  server_name _;

  client_max_body_size 50m;

  location / {
    proxy_pass http://127.0.0.1:${APP_PORT};
    proxy_http_version 1.1;

    proxy_set_header Host \$host;
    proxy_set_header X-Real-IP \$remote_addr;
    proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;

    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;
  }
}
EOF

sudo ln -sf "${NGINX_AVAIL}" "${NGINX_ENABLED}"

if [ -f /etc/nginx/sites-enabled/default ]; then
  sudo rm -f /etc/nginx/sites-enabled/default
fi

echo "==> Testing nginx"
sudo nginx -t
sudo systemctl reload nginx

echo "==> Configuring systemd service"

if [ -f "${UNIT_FILE}" ] && [ "${OVERWRITE_UNIT}" != "true" ]; then
  echo "Systemd unit already exists. Skipping (overwrite=false)."
else
  sudo bash -c "cat > '${UNIT_FILE}'" <<EOF
[Unit]
Description=${SERVICE_NAME} Service
After=network.target

[Service]
Type=simple
WorkingDirectory=${DEPLOY_PATH}
ExecStart=${DEPLOY_PATH}/${APP_BINARY}
Restart=always
RestartSec=3
Environment=PORT=${ENV_PORT}

[Install]
WantedBy=multi-user.target
EOF

  sudo systemctl daemon-reload
  sudo systemctl enable "${SERVICE_NAME}"
fi

echo "==> Restarting service if binary exists"
if [ -x "${DEPLOY_PATH}/${APP_BINARY}" ]; then
  sudo systemctl restart "${SERVICE_NAME}"
  sudo systemctl --no-pager status "${SERVICE_NAME}" -l || true
else
  echo "Binary not found at ${DEPLOY_PATH}/${APP_BINARY}"
  echo "Deploy your binary first, then restart:"
  echo "sudo systemctl restart ${SERVICE_NAME}"
fi

echo "Startup complete."