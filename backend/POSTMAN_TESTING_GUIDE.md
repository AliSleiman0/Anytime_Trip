# Testing Booking Email Notifications with Postman

## Quick Start

### 1. Import the Postman Collection

1. Open Postman
2. Click **Import** button
3. Select the file: `backend/postman_booking_email_test.json`
4. The collection "Anytime Travel - Booking Email Notifications" will appear

### 2. Set Up Environment (Optional)

Create a new environment or use the requests directly. The collection will auto-save:
- `auth_token` - JWT token after login
- `user_id` - User ID after login

### 3. Start the Backend Server

```bash
cd /home/kongo/Anytime_Trip/backend
go run .
```

The server should start on `http://localhost:8080`

### 4. Configure SMTP (Optional)

For actual email delivery, set these environment variables before starting the server:

```bash
export SMTP_HOST=mail.smtp2go.com
export SMTP_PORT=2525
export SENDER_EMAIL=info@anytimetravel.app
export SENDER_PASSWORD=your_smtp2go_password_here
```

**Without SMTP credentials:** Emails will be logged to console only (development mode).

---

## Step-by-Step Testing Guide

### Step 1: Create and Activate User Account

**Request:** `1. User Signup`

Creates a test user account.

**Expected Response:**
```json
{
  "message": "User created successfully. Please verify your phone number.",
  "user": {
    "id": "...",
    "email": "testuser@example.com",
    ...
  }
}
```

**Note:** The user account needs OTP verification. You can either:
- Skip activation and directly update the user in MongoDB to set `is_active: true`
- Or use an existing active user

**Quick MongoDB Fix:**
```javascript
db.users.updateOne(
  { email: "testuser@example.com" },
  { $set: { is_active: true } }
)
```

### Step 2: Login

**Request:** `2. User Login`

Logs in and automatically saves the JWT token.

**Expected Response:**
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": "...",
    "email": "testuser@example.com",
    ...
  }
}
```

The token is automatically saved to `{{auth_token}}` variable.

### Step 3: Enable Email Notifications

**Request:** `3. Enable Email Notifications`

Enables email notifications for the user (Priority 2 check).

**Expected Response:**
```json
{
  "message": "Notification preferences updated successfully"
}
```

### Step 4: Enable Admin Notifications

**IMPORTANT:** Before creating bookings, enable admin notification preferences:

1. Open browser: `http://localhost:8080/admin/login`
2. Login with:
   - Email: `admin@anytime.com`
   - Password: `admin123`
3. Go to **Settings → Notification Preferences**
4. Under **Bookings** section, enable **"New Booking"**
5. Click **Save Changes**

This is the **Priority 1** check that must pass for emails to be sent.

### Step 5: Create Hotel Booking

**Request:** `4. Create Hotel Booking`

Creates a hotel booking and triggers email notification.

**Expected Response:**
```json
{
  "message": "Hotel booking created successfully",
  "booking": {
    "id": "...",
    "booking_id": "HBK-...",
    "amount": 299.99,
    ...
  }
}
```

**Check the console output** for:
```
[NOTIFICATION] Sending new booking email to testuser@example.com for booking HBK-...
[NOTIFICATION] Successfully sent new booking email for HBK-...
```

Or if emails are blocked:
```
[NOTIFICATION] Email blocked by admin preferences for booking HBK-...
[NOTIFICATION] Email blocked by user preferences for user ...
```

### Step 6: Create Other Bookings (Optional)

Repeat for:
- **Flight Booking:** Request `5. Create Flight Booking`
- **Car Booking:** Request `6. Create Car Booking`
- **Transfer Booking:** Request `7. Create Transfer Booking`

---

## Verification Checklist

After creating a booking:

### ✅ Check Backend Console

Look for these log messages:

**Success:**
```
[NOTIFICATION] Sending new booking email to testuser@example.com for booking HBK-abc123
[NOTIFICATION] Successfully sent new booking email for HBK-abc123
```

**Development Mode (no SMTP):**
```
[EMAIL] New booking notification for testuser@example.com: HBK-abc123 (Hotel Booking)
```

**Admin Blocked:**
```
[NOTIFICATION] Email blocked by admin preferences for booking HBK-abc123
```

**User Blocked:**
```
[NOTIFICATION] Email blocked by user preferences for user user_123
```

### ✅ Check Email Inbox

If SMTP credentials are configured:
1. Check inbox for: `testuser@example.com`
2. Subject: `Booking Confirmation - HBK-...`
3. From: `info@anytimetravel.app`
4. Check spam folder if not in inbox

### ✅ Check Database

Verify booking was created:
```javascript
db.hotel_bookings.find().sort({ created_at: -1 }).limit(1).pretty()
```

---

## Common Issues & Solutions

### Issue: "Unauthorized" Error

**Solution:** Run request `2. User Login` again to refresh the token.

### Issue: Email Not Sending

**Check Priority 1 - Admin Preferences:**
```javascript
db.notification_preferences.find().pretty()
```
Ensure at least one admin has `bookings.new_booking: true`

**Check Priority 2 - User Preferences:**
```javascript
db.users.findOne({ email: "testuser@example.com" }, { notification_preferences: 1 })
```
Ensure `notification_preferences.Email: true`

### Issue: User Not Active

After signup, activate the user:
```javascript
db.users.updateOne(
  { email: "testuser@example.com" },
  { $set: { is_active: true } }
)
```

### Issue: SMTP Errors

Check your SMTP2GO credentials:
- Verify password is correct
- Check SMTP2GO account is active
- Try regenerating SMTP user credentials

---

## Testing Different Scenarios

### Test 1: Email Blocked by Admin
1. Disable admin notification: Settings → Notifications → Uncheck "New Booking"
2. Create a booking
3. Expected: `[NOTIFICATION] Email blocked by admin preferences`

### Test 2: Email Blocked by User
1. Run request `3. Enable Email Notifications` with `"Email": false`
2. Create a booking
3. Expected: `[NOTIFICATION] Email blocked by user preferences`

### Test 3: Development Mode
1. Don't set `SENDER_PASSWORD`
2. Create a booking
3. Expected: Email logged to console only

### Test 4: Production Mode
1. Set valid SMTP credentials
2. Create a booking
3. Expected: Actual email sent to inbox

---

## API Endpoints Summary

| Endpoint | Method | Auth Required | Purpose |
|----------|--------|---------------|---------|
| `/api/app/signup` | POST | No | Create user account |
| `/api/app/login` | POST | No | Get JWT token |
| `/api/app/notification-preferences` | POST | Yes | Enable/disable email |
| `/api/app/bookings/hotel` | POST | Yes | Create hotel booking |
| `/api/app/bookings/flight` | POST | Yes | Create flight booking |
| `/api/app/bookings/car` | POST | Yes | Create car booking |
| `/api/app/bookings/transfer` | POST | Yes | Create transfer booking |

---

## Quick Test Commands

```bash
# Start server with SMTP (production mode)
export SENDER_PASSWORD=your_password
cd /home/kongo/Anytime_Trip/backend
go run .

# Start server without SMTP (development mode)
cd /home/kongo/Anytime_Trip/backend
go run .

# Watch logs in real-time
cd /home/kongo/Anytime_Trip/backend
go run . | grep -E "NOTIFICATION|EMAIL|BOOKING"
```

---

## Expected Email Template

When an email is successfully sent, the user receives:

**Subject:** Booking Confirmation - [BookingID]

**Content:**
```
✓ Booking Confirmed

Dear [User Name],

Thank you for choosing Anytime Travel! We're pleased to confirm your booking.

Booking Details:
- Booking ID: HBK-abc123
- Booking Type: Hotel Booking
- Booking Date: Monday, March 15, 2026
- Total Amount: $299.99

[View My Bookings Button]

Best regards,
Anytime Travel Team
```

---

**Ready to test?** Import the Postman collection and start with request #1! 🚀
