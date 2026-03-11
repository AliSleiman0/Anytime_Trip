#!/bin/bash
# Setup SSL/HTTPS with Let's Encrypt for Anytime Trip

set -e

DOMAIN=${1:-"anytimetravel.app"}
EMAIL=${2:-"admin@anytimetravel.app"}

echo "🔒 Setting up SSL/HTTPS with Let's Encrypt..."
echo "🌐 Domain: $DOMAIN"
echo "📧 Email: $EMAIL"

# Install certbot if not already installed
if ! command -v certbot &> /dev/null; then
    echo "📦 Installing certbot..."
    sudo apt-get update
    sudo apt-get install -y certbot python3-certbot-nginx
fi

# Check if certificate already exists
if sudo certbot certificates 2>/dev/null | grep -q "$DOMAIN"; then
    echo "✅ SSL certificate already exists for $DOMAIN"
    echo "🔄 Renewing certificate if needed..."
    sudo certbot renew --nginx --quiet
else
    echo "🆕 Obtaining new SSL certificate..."
    # Obtain certificate (non-interactive)
    sudo certbot --nginx \
        -d "$DOMAIN" \
        -d "www.$DOMAIN" \
        --non-interactive \
        --agree-tos \
        --email "$EMAIL" \
        --redirect
fi

# Setup auto-renewal cron job if not exists
if ! sudo crontab -l 2>/dev/null | grep -q "certbot renew"; then
    echo "⏰ Setting up auto-renewal cron job..."
    (sudo crontab -l 2>/dev/null; echo "0 3 * * * certbot renew --quiet --post-hook 'systemctl reload nginx'") | sudo crontab -
fi

# Test Nginx configuration
echo "🧪 Testing Nginx configuration..."
sudo nginx -t

# Reload Nginx
echo "🔄 Reloading Nginx..."
sudo systemctl reload nginx

echo "✅ SSL/HTTPS setup complete!"
echo "🔒 Your site is now accessible at: https://$DOMAIN"
echo ""
echo "📝 Certificate details:"
sudo certbot certificates | grep -A 3 "$DOMAIN" || true
