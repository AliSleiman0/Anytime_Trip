# My Bookings API Endpoints

## Overview
API endpoints to fetch and manage user bookings from the mobile app.

## Base URL
`/api/app`

## Authentication
All endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

### 1. Get All My Bookings
Get all bookings (car, flight, hotel, transfer) for the authenticated user.

**Endpoint:** `GET /my-bookings`

**Response:**
```json
{
  "success": true,
  "data": {
    "car_bookings": [...],
    "flight_bookings": [...],
    "hotel_bookings": [...],
    "transfer_bookings": [...]
  }
}
```

---

### 2. Get My Car Bookings
**Endpoint:** `GET /my-bookings/car`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "...",
      "booking_id": "CBK-12345678",
      "status": "confirmed",
      "car_type": "SUV",
      "pickup": {
        "location": "BEY Airport",
        "date": "2026-02-13T10:00:00Z",
        "time": "10:00 AM"
      },
      "pricing": {
        "total": 115.00
      },
      ...
    }
  ]
}
```

---

### 3. Get My Flight Bookings
**Endpoint:** `GET /my-bookings/flight`

**Response:** Similar structure with flight-specific fields

---

### 4. Get My Hotel Bookings
**Endpoint:** `GET /my-bookings/hotel`

**Response:** Similar structure with hotel-specific fields

---

### 5. Get My Transfer Bookings
**Endpoint:** `GET /my-bookings/transfer`

**Response:** Similar structure with transfer-specific fields

---

### 6. Get Specific Booking by ID
Get details of a specific booking.

**Endpoint:** `GET /bookings/:type/:id`

**Parameters:**
- `type`: Booking type (`car`, `flight`, `hotel`, `transfer`)
- `id`: Booking ID

**Example:** `GET /bookings/car/abc123`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "abc123",
    "booking_id": "CBK-12345678",
    ...
  }
}
```

**Error Response (404):**
```json
{
  "error": "Booking not found"
}
```

---

### 7. Cancel Booking
Cancel a booking.

**Endpoint:** `POST /bookings/:type/:id/cancel`

**Parameters:**
- `type`: Booking type (`car`, `flight`, `hotel`, `transfer`)
- `id`: Booking ID

**Example:** `POST /bookings/car/abc123/cancel`

**Response:**
```json
{
  "success": true,
  "message": "Booking cancelled successfully"
}
```

**Error Responses:**
- `401`: Unauthorized (missing/invalid token)
- `404`: Booking not found or doesn't belong to user
- `500`: Failed to cancel booking

---

## Usage in Mobile App

### Flutter/Dart Example:
```dart
// Get all bookings
final response = await http.get(
  Uri.parse('$baseUrl/api/app/my-bookings'),
  headers: {
    'Authorization': 'Bearer $token',
  },
);

if (response.statusCode == 200) {
  final data = jsonDecode(response.body);
  final carBookings = data['data']['car_bookings'];
  final flightBookings = data['data']['flight_bookings'];
  // ... use the data
}

// Cancel a booking
final cancelResponse = await http.post(
  Uri.parse('$baseUrl/api/app/bookings/car/$bookingId/cancel'),
  headers: {
    'Authorization': 'Bearer $token',
  },
);
```

---

## Implementation Notes

1. All endpoints validate that the booking belongs to the authenticated user
2. Cancelled bookings remain in the database with `status: "cancelled"`
3. Bookings are sorted by creation date (most recent first)
4. The `/my-bookings` endpoint returns empty arrays `[]` if no bookings exist for a type

---

## Next Steps

To integrate into your Flutter app:
1. Create a `BookingService` class to handle API calls
2. Create models for each booking type
3. Update your "My Bookings" screen to fetch data from these endpoints
4. Add pull-to-refresh functionality
5. Implement booking details view
6. Add cancel booking confirmation dialog
