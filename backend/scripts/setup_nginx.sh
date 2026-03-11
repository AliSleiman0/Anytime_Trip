#!/bin/bash
# Setup Nginx reverse proxy for Anytime Trip backend

set -e

echo "🔧 Setting up Nginx reverse proxy..."

# Install Nginx if not already installed
if ! command -v nginx &> /dev/null; then
    echo "📦 Installing Nginx..."
    sudo apt-get update
    sudo apt-get install -y nginx
fi

# Get the domain from argument or use default
DOMAIN=${1:-"anytimetrip.online"}
BACKEND_PORT=${2:-"8080"}

echo "🌐 Configuring Nginx for domain: $DOMAIN"
echo "🔌 Backend port: $BACKEND_PORT"

# Create Nginx configuration
sudo tee /etc/nginx/sites-available/anytime-trip > /dev/null <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name $DOMAIN www.$DOMAIN;

    client_max_body_size 10M;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy to backend
    location / {
        proxy_pass http://localhost:$BACKEND_PORT;
        proxy_http_version 1.1;
        
        # WebSocket support
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        
        # Forward real IP
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Static files (if served separately)
    location /static/ {
        alias /root/anytime-trip/backend/static/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    location /uploads/ {
        alias /root/anytime-trip/backend/static/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Health check endpoint (optional)
    location /nginx-health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF

# Enable the site
sudo ln -sf /etc/nginx/sites-available/anytime-trip /etc/nginx/sites-enabled/

# Remove default site if exists
sudo rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
echo "🧪 Testing Nginx configuration..."
sudo nginx -t

# Reload Nginx
echo "🔄 Reloading Nginx..."
sudo systemctl enable nginx
sudo systemctl restart nginx

echo "✅ Nginx reverse proxy setup complete!"
echo "🌐 Your app should now be accessible at: http://$DOMAIN"
echo ""
echo "📝 Next steps:"
echo "   1. Ensure your backend is running on port $BACKEND_PORT"
echo "   2. For HTTPS, run: sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN"
