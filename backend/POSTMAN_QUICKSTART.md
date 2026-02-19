# 🚀 Quick Start: Test Booking Emails with Postman

## 📥 Setup (2 minutes)

1. **Import Collection**
   - Open Postman
   - Import: `backend/postman_booking_email_test.json`

2. **Start Backend**
   ```bash
   cd /home/kongo/Anytime_Trip/backend
   go run .
   ```

3. **Optional: Set SMTP** (for real emails)
   ```bash
   export SENDER_PASSWORD=your_smtp2go_password
   ```

## 🎯 Test Flow (5 steps)

```
1. Signup → Create user
2. Login → Get JWT token (auto-saved)
3. Enable Email Notifications → Turn on emails
4. Enable Admin Notifications → Admin panel
5. Create Booking → Email sent! ✉️
```

## 📋 Quick Steps

### In Postman:
1. ▶️ Run: `1. User Signup`
2. ▶️ Run: `2. User Login` (saves token)
3. ▶️ Run: `3. Enable Email Notifications`

### In MongoDB (activate user):
```javascript
db.users.updateOne(
  { email: "testuser@example.com" },
  { $set: { is_active: true } }
)
```

### In Browser:
1. Go to: http://localhost:8080/admin/login
2. Login: `admin@anytime.com` / `admin123`
3. Settings → Notifications → Enable "New Booking" ✅
4. Save

### Back in Postman:
5. ▶️ Run: `4. Create Hotel Booking`

## ✅ Success Check

**Console shows:**
```
[NOTIFICATION] Sending new booking email to testuser@example.com
[NOTIFICATION] Successfully sent new booking email
```

**Postman Response:**
```json
{
  "message": "Hotel booking created successfully",
  "booking": { "booking_id": "HBK-..." }
}
```

**Email inbox:** Check `testuser@example.com` (or spam folder)

## 🔧 Troubleshooting

| Problem | Solution |
|---------|----------|
| "Unauthorized" | Run `2. User Login` again |
| Email blocked by admin | Enable in admin panel |
| Email blocked by user | Run request #3 |
| User not active | Update in MongoDB (see above) |
| No email received | Check console logs for errors |

## 📝 Change Test Email

In each request, replace `testuser@example.com` with your email to receive actual emails.

## 📚 Full Guide

See: [`POSTMAN_TESTING_GUIDE.md`](POSTMAN_TESTING_GUIDE.md)

---

**That's it! Import collection → Run requests → Get emails** 🎉
