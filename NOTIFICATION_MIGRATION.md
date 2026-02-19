# Notification Preferences Migration

## Changes Made

The notification preferences have been simplified from 10 options to 3 essential options:

### Old Structure (Removed):
- All Notifications
- Push Notifications
- Email Alerts
- SMS Updates
- In-App Messages
- Social Media Alerts
- Webhook Notifications
- Browser Notifications
- RSS Feed Updates
- Chatbot Messages

### New Structure (Active):
1. **Email** - Email notifications
2. **SMS** - SMS text message notifications
3. **Chatbot** - In-app chatbot notifications

## Files Updated

### Frontend (Flutter App)
- `app/lib/features/account/view/notifications_page.dart`
  - Updated to show only 3 notification options
  - Fixed initialization to prevent red error screen
  - Filters any saved preferences to only keep the 3 new options

### Backend (Go)
- `backend/internal/models/app/user.go`
  - Updated `NotificationPreferences` struct to only include Email, SMS, and Chatbot fields
  - Changed field names to match frontend expectations

### Migration Script
- `backend/scripts/migrate_notification_preferences.go`
  - Script to migrate existing user data from old structure to new structure
  - Maps old field names to new ones

## Running the Migration

To migrate existing user data in the database:

```bash
cd backend/scripts
go run migrate_notification_preferences.go
```

This will update all existing users' notification preferences to use the new structure while preserving their previous Email, SMS, and Chatbot settings.

## API Endpoints (Unchanged)

The following endpoints continue to work as before:

- `GET /app/notification-preferences` - Get user's notification preferences
- `POST /app/notification-preferences` - Save/update user's notification preferences

## Default Values

New users (or users without saved preferences) will have all three options enabled by default:
- Email: true
- SMS: true
- Chatbot: true
