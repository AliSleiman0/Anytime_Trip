#!/bin/bash

# Quick Test Script for Booking Email Notifications
# This script helps you quickly test the email notification feature

echo "=========================================="
echo "Booking Email Notification Quick Test"
echo "=========================================="
echo ""

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    echo "Error: Please run this script from the backend directory"
    exit 1
fi

# Step 1: Check environment variables
echo "[Step 1/4] Checking SMTP configuration..."
if [ -z "$SMTP_HOST" ]; then
    echo "⚠️  SMTP_HOST not set. Using default: mail.smtp2go.com"
    export SMTP_HOST="mail.smtp2go.com"
fi

if [ -z "$SMTP_PORT" ]; then
    echo "⚠️  SMTP_PORT not set. Using default: 2525"
    export SMTP_PORT="2525"
fi

if [ -z "$SENDER_EMAIL" ]; then
    echo "⚠️  SENDER_EMAIL not set. Using default: info@anytimetravel.app"
    export SENDER_EMAIL="info@anytimetravel.app"
fi

if [ -z "$SENDER_PASSWORD" ]; then
    echo "⚠️  SENDER_PASSWORD not set. Running in DEVELOPMENT MODE (no emails sent)"
    echo "   Set SENDER_PASSWORD to enable actual email delivery"
else
    echo "✓ SMTP credentials configured"
fi

echo ""
echo "SMTP Configuration:"
echo "  Host: $SMTP_HOST"
echo "  Port: $SMTP_PORT"
echo "  From: $SENDER_EMAIL"
echo ""

# Step 2: Build the test script
echo "[Step 2/4] Building test script..."
if ! go build -o bin/test_email test_email.go; then
    echo "✗ Failed to build test script"
    exit 1
fi
echo "✓ Test script built successfully"
echo ""

# Step 3: Run the test
echo "[Step 3/4] Running email notification test..."
echo "----------------------------------------"
./bin/test_email
TEST_RESULT=$?
echo "----------------------------------------"
echo ""

# Step 4: Summary
echo "[Step 4/4] Test Summary"
if [ $TEST_RESULT -eq 0 ]; then
    echo "✓ Test completed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Check the console output above for email details"
    echo "2. If SENDER_PASSWORD is set, check your email inbox (test@example.com)"
    echo "3. Check spam folder if email not in inbox"
    echo ""
    echo "To test with different settings:"
    echo "  - Enable admin notifications: Admin Panel → Settings → Notifications"
    echo "  - Change test email in: backend/test_email.go (line 41)"
else
    echo "✗ Test failed with exit code: $TEST_RESULT"
    echo ""
    echo "Troubleshooting:"
    echo "1. Check if MongoDB is running"
    echo "2. Verify MONGO_URI environment variable is set"
    echo "3. Check logs above for specific errors"
fi

echo ""
echo "For detailed testing guide, see: BOOKING_EMAIL_TESTING.md"
echo "=========================================="
