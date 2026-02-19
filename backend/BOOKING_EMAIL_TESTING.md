# Testing Booking Email Notifications

This guide explains how to test the email notification feature for hotel, car, flight, and transfer bookings.

## Prerequisites

### 1. SMTP2GO Configuration

Add these environment variables to your system or `.env` file:

```bash
# SMTP Configuration for SMTP2GO
SMTP_HOST=mail.smtp2go.com
SMTP_PORT=2525
SENDER_EMAIL=info@anytimetravel.app
SENDER_PASSWORD=your_smtp2go_password_here
```

**How to get SMTP2GO credentials:**
1. Go to https://www.smtp2go.com
2. Create an account or log in
3. Navigate to Settings → Users → Add SMTP User
4. Copy the SMTP username and password
5. Use the password in `SENDER_PASSWORD` environment variable

### 2. Start the Backend Server

```bash
cd /home/kongo/Anytime_Trip/backend
go run .
```

The server should start on port 8080 (or your configured port).

## Testing Methods

### Method 1: Using the Seed Data (Quickest)

The seed files already create bookings. To test email notifications:

1. **Enable Admin Notification Preferences** (Priority 1)
   - Log in to admin panel: http://localhost:8080/admin/login
   - Go to Settings → Notification Preferences
   - Enable "New Booking" notifications under Bookings section
   - Save changes

2. **Set User Email Preferences** (Priority 2)
   - User preferences default to `Email: false` for new users
   - You need to update a test user's preferences via API or database

3. **Run Seed Script to Create Test Booking**
   ```bash
   cd /home/kongo/Anytime_Trip/backend
   go run seeds/hotel_booking_seed.go
   ```

**Note:** The current implementation doesn't trigger emails from seed scripts. You need to use the API endpoints.

### Method 2: Using cURL (API Testing)

Since the booking creation endpoints aren't implemented yet in the handler, let's test the notification helper directly.

#### Create a Test Script

Create `/home/kongo/Anytime_Trip/backend/test_email.go`:

```go
//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/core/utils"
	"Anytime_Travel/backend/internal/database"
	"Anytime_Travel/backend/internal/models/app"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	apprepo "Anytime_Travel/backend/internal/repository/app"

	"github.com/google/uuid"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := apprepo.NewUserRepository(db.DB)
	adminNotifPrefsRepo := adminrepo.NewNotificationPreferencesRepository(db.DB)
	hotelBookingRepo := apprepo.NewHotelBookingRepository(db.DB)

	// Initialize email service and notification helper
	emailService := utils.NewEmailService()
	notificationHelper := utils.NewNotificationHelper(
		emailService,
		adminNotifPrefsRepo,
		userRepo,
	)

	ctx := context.Background()

	// Step 1: Find or create a test user
	testEmail := "test@example.com"
	user, err := userRepo.FindByEmail(ctx, testEmail)
	if err != nil {
		// Create test user
		log.Println("Creating test user...")
		user = &app.User{
			ID:          uuid.New().String(),
			Name:        "Test User",
			Email:       testEmail,
			PhoneNumber: "+1234567890",
			IsActive:    true,
			NotificationPreferences: app.NotificationPreferences{
				Email:   true, // Enable email notifications
				SMS:     false,
				Chatbot: false,
			},
			CreatedAt: time.Now(),
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatal("Failed to create test user:", err)
		}
	} else {
		// Update user to enable email notifications
		log.Println("Updating test user email preferences...")
		user.NotificationPreferences.Email = true
		if err := userRepo.Update(ctx, user.ID, user); err != nil {
			log.Fatal("Failed to update user:", err)
		}
	}

	log.Printf("Test user: %s (%s)\n", user.Name, user.Email)

	// Step 2: Create a test hotel booking
	log.Println("Creating test hotel booking...")
	booking := &app.HotelBooking{
		ID:          uuid.New().String(),
		BookingID:   "TEST-HB-" + time.Now().Format("20060102150405"),
		UserID:      user.ID,
		HotelID:     "hotel_001",
		Status:      app.HotelBookingStatusConfirmed,
		Customer: app.HotelCustomer{
			Name:  user.Name,
			Email: user.Email,
		},
		Details:       "Test Deluxe Room - 2 Nights",
		BookingDate:   time.Now(),
		Amount:        299.99,
		Currency:      "USD",
		PaymentStatus: app.HotelPaymentStatusPaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := hotelBookingRepo.Create(ctx, booking); err != nil {
		log.Fatal("Failed to create booking:", err)
	}

	log.Printf("Created booking: %s\n", booking.BookingID)

	// Step 3: Send email notification
	log.Println("Sending email notification...")
	if err := notificationHelper.SendHotelBookingEmail(ctx, booking); err != nil {
		log.Printf("Error sending email: %v\n", err)
	} else {
		log.Println("✓ Email notification sent successfully!")
	}

	log.Println("\nTest completed. Check the email inbox for:", user.Email)
}
```

#### Run the Test Script

```bash
cd /home/kongo/Anytime_Trip/backend
go run test_email.go
```

### Method 3: Test via MongoDB Direct Insert + Manual Trigger

1. **Insert a test booking directly in MongoDB:**

```javascript
// Connect to MongoDB
use anytime_travel

// Insert test booking
db.hotel_bookings.insertOne({
  "_id": "test_booking_001",
  "booking_id": "HB-TEST-001",
  "user_id": "user_001",  // Make sure this user exists
  "hotel_id": "hotel_001",
  "status": "confirmed",
  "customer": {
    "name": "John Doe",
    "email": "john.doe@example.com"
  },
  "details": "Test Suite - 3 Nights",
  "booking_date": new Date(),
  "amount": 450.00,
  "currency": "USD",
  "payment_status": "paid",
  "created_at": new Date(),
  "updated_at": new Date()
})
```

2. **Enable user email notifications:**

```javascript
db.users.updateOne(
  { "_id": "user_001" },
  { 
    $set: { 
      "notification_preferences.Email": true 
    } 
  }
)
```

3. Use the test script above to trigger the email.

## Testing Checklist

- [ ] SMTP2GO credentials configured in environment variables
- [ ] Backend server running
- [ ] Admin has enabled "New Booking" notifications
- [ ] Test user has email notifications enabled (`notification_preferences.Email: true`)
- [ ] Test booking created successfully
- [ ] Email sent without errors
- [ ] Email received in inbox (check spam folder)

## Expected Email Format

The email should contain:
- **Subject:** `Booking Confirmation - [BookingID]`
- **From:** `info@anytimetravel.app`
- **To:** User's email address
- **Content:**
  - Booking confirmation message
  - Booking ID
  - Booking Type (Hotel/Flight/Car/Transfer)
  - Booking Date
  - Total Amount
  - "View My Bookings" button

## Troubleshooting

### Email Not Sending

1. **Check SMTP credentials:**
   ```bash
   echo $SMTP_HOST
   echo $SMTP_PORT
   echo $SENDER_EMAIL
   # Don't echo password for security
   ```

2. **Check console output:**
   - Look for `[EMAIL]` or `[NOTIFICATION]` log messages
   - If you see "Password reset link for..." it means SMTP password is empty (development mode)

3. **Verify admin preferences:**
   ```javascript
   db.notification_preferences.find().pretty()
   ```
   - Check if any admin has `bookings.new_booking: true`

4. **Verify user preferences:**
   ```javascript
   db.users.findOne({ "_id": "user_001" }, { notification_preferences: 1 })
   ```
   - Check if `notification_preferences.Email: true`

### Email Blocked by Admin

If you see:
```
[NOTIFICATION] Email blocked by admin preferences for booking XYZ
```

Solution: Enable admin notification preferences in the admin panel.

### Email Blocked by User

If you see:
```
[NOTIFICATION] Email blocked by user preferences for user XYZ
```

Solution: Update user's notification preferences to enable email.

### SMTP Authentication Error

If you see SMTP authentication errors:
- Verify your SMTP2GO username and password
- Check if your SMTP2GO account is active
- Try generating a new SMTP user in SMTP2GO dashboard

## Development Mode (No Email Sent)

If `SENDER_PASSWORD` is empty, the system runs in development mode and only logs emails to console:

```
[EMAIL] New booking notification for user@example.com: HB-001 (Hotel Booking)
```

This is useful for testing without actual email delivery.

## Production Checklist

Before deploying to production:

- [ ] Valid SMTP2GO account with sufficient credits
- [ ] SMTP credentials stored securely (not in code)
- [ ] Test email delivery to real addresses
- [ ] Configure proper sender domain (info@anytimetravel.app)
- [ ] Set up SPF and DKIM records for your domain
- [ ] Test all booking types (hotel, flight, car, transfer)
- [ ] Verify admin notification preferences work correctly
- [ ] Verify user notification preferences work correctly
- [ ] Monitor SMTP2GO usage and delivery rates
