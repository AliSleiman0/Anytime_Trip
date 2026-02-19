# Booking Email Notification System

## Overview

This system sends email notifications to users when they create new bookings (hotel, flight, car, transfer). The system implements a two-tier preference checking system:

**Priority 1**: Admin Notification Center - Controls whether booking emails are sent globally  
**Priority 2**: User Notification Preferences - Controls whether individual users receive emails

## Email Configuration

### SMTP2GO Setup

The system is configured to use SMTP2GO with the following settings:

- **SMTP Host**: `mail.smtp2go.com`
- **SMTP Port**: `2525`
- **Sender Email**: `info@anytimetravel.app`
- **Sender Password**: Set via environment variable `SENDER_PASSWORD`

### Environment Variables

Add these to your `.env` file or environment:

```bash
SMTP_HOST=mail.smtp2go.com
SMTP_PORT=2525
SENDER_EMAIL=info@anytimetravel.app
SENDER_PASSWORD=your_smtp2go_password_here
```

## How It Works

### 1. Admin Notification Center (Priority 1)

Admins can control whether booking notification emails are sent through their notification preferences:

- **Location**: Admin Dashboard → Notification Settings → Bookings
- **Setting**: `new_booking` toggle
- **Behavior**: 
  - If ANY admin has `new_booking` enabled → emails are allowed
  - If ALL admins have `new_booking` disabled → emails are blocked
  - If no admin preferences exist → emails are allowed (default)

### 2. User Notification Preferences (Priority 2)

Individual users can control their own email notifications:

- **Location**: User Profile → Notification Preferences
- **Setting**: `Email` toggle
- **Behavior**:
  - If user has `Email` enabled → emails are sent (if admin allows)
  - If user has `Email` disabled → emails are blocked

### 3. Email Flow

```
User creates booking
    ↓
Check Admin Preferences (Priority 1)
    ↓ (if allowed)
Check User Preferences (Priority 2)
    ↓ (if allowed)
Send Email via SMTP2GO
```

## Implementation Details

### Files Created/Modified

1. **`/backend/core/utils/email.go`**
   - Updated SMTP configuration for SMTP2GO
   - Added `SendNewBookingNotification()` method
   - Professional HTML email template

2. **`/backend/core/utils/notification_helper.go`** (NEW)
   - Central notification logic with preference checking
   - Methods for each booking type:
     - `SendHotelBookingEmail()`
     - `SendFlightBookingEmail()`
     - `SendCarBookingEmail()`
     - `SendTransferBookingEmail()`

3. **`/backend/services/booking_service.go`** (NEW)
   - Service layer for booking operations
   - Handles booking creation + email notifications
   - Non-blocking email sending

4. **`/backend/internal/handlers/app/handler.go`**
   - Added booking repositories and notification helper
   - New handler methods:
     - `CreateHotelBooking()`
     - `CreateFlightBooking()`
     - `CreateCarBooking()`
     - `CreateTransferBooking()`

5. **`/backend/internal/routes/app/routes.go`**
   - Added booking API endpoints:
     - `POST /api/app/bookings/hotel`
     - `POST /api/app/bookings/flight`
     - `POST /api/app/bookings/car`
     - `POST /api/app/bookings/transfer`

6. **`/backend/internal/repository/admin/notification_preferences_repository.go`**
   - Added `FindAll()` method to fetch all admin preferences

7. **`/backend/main.go`**
   - Initialize email service and notification helper
   - Pass dependencies to AppHandler

## API Usage

### Create Hotel Booking

```bash
POST /api/app/bookings/hotel
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "hotel_id": "hotel_001",
  "status": "pending",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "details": "Deluxe room, 2 nights",
  "booking_date": "2026-03-15T00:00:00Z",
  "amount": 500.00,
  "currency": "USD",
  "payment_status": "pending"
}
```

### Create Flight Booking

```bash
POST /api/app/bookings/flight
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "flight_id": "flight_001",
  "status": "pending",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "details": "NYC-LAX AA120",
  "booking_date": "2026-03-20T00:00:00Z",
  "amount": 450.00,
  "currency": "USD",
  "payment_status": "pending"
}
```

### Create Car Booking

```bash
POST /api/app/bookings/car
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "car_id": "car_001",
  "status": "pending",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "details": "SUV rental - 3 days",
  "booking_date": "2026-03-10T00:00:00Z",
  "amount": 300.00,
  "currency": "USD",
  "payment_status": "pending"
}
```

### Create Transfer Booking

```bash
POST /api/app/bookings/transfer
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "transfer_id": "transfer_001",
  "status": "pending",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "details": "Airport pickup - Terminal 3",
  "booking_date": "2026-03-25T00:00:00Z",
  "amount": 75.00,
  "currency": "USD",
  "payment_status": "pending"
}
```

## Testing

### 1. Test with Development Mode (No SMTP Credentials)

If `SENDER_PASSWORD` is not set, emails will be logged to console instead of sent:

```
[EMAIL] New booking notification for john@example.com: HBK-abc123 (Hotel Booking)
```

### 2. Test with SMTP2GO

1. Set up SMTP2GO account and get credentials
2. Add credentials to environment variables
3. Create a booking via API
4. Check email inbox for confirmation

### 3. Test Admin Preferences

1. Go to Admin Dashboard → Notification Settings
2. Toggle "New Booking" under Bookings section
3. Create a booking
4. Verify email is sent/blocked based on setting

### 4. Test User Preferences

1. Go to User Profile → Notification Preferences
2. Toggle "Email" preference
3. Create a booking
4. Verify email is sent/blocked based on setting

## Email Template

The email includes:

- Professional header with Anytime Travel branding
- Booking confirmation message
- Booking details:
  - Booking ID
  - Booking Type (Hotel/Flight/Car/Transfer)
  - Booking Date
  - Total Amount
- Call-to-action button to view bookings
- Footer with contact information

## Error Handling

- Email sending is non-blocking - booking will succeed even if email fails
- Errors are logged but don't affect booking creation
- Admin preference errors default to allowing emails
- Missing user preferences default to database values

## Logs

The system logs important events:

```
[NOTIFICATION] Sending new booking email to john@example.com for booking HBK-abc123
[NOTIFICATION] Successfully sent new booking email for HBK-abc123
[NOTIFICATION] Email blocked by admin preferences for booking HBK-abc123
[NOTIFICATION] Email blocked by user preferences for user user_123
[BOOKING] Warning: Failed to send booking notification: connection timeout
```

## Future Enhancements

1. Add booking cancellation emails
2. Add booking confirmation emails (separate from initial booking)
3. Add booking reminder emails
4. Add SMS notifications
5. Add in-app notifications (chatbot)
6. Add email templates for different languages
7. Add email tracking/analytics

## Troubleshooting

### Emails not sending

1. Check `SENDER_PASSWORD` is set correctly
2. Verify SMTP2GO credentials are valid
3. Check admin notification preferences
4. Check user notification preferences
5. Review server logs for error messages

### Emails going to spam

1. Configure SPF/DKIM records for your domain
2. Use SMTP2GO's domain verification
3. Ensure sender email matches verified domain
4. Add proper unsubscribe links (future enhancement)

### Testing locally

1. Use a tool like MailHog or Mailtrap for local testing
2. Or leave `SENDER_PASSWORD` empty to see console logs
3. Or use your personal SMTP credentials for testing
