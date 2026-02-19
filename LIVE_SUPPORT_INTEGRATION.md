# Live Support System Integration - Flutter App to Admin Web Dashboard

## Overview
This implementation connects the Flutter mobile app's chatbot feature with the admin web dashboard's live support system using WebSockets for real-time messaging.

## Architecture

### Backend (Go/Fiber)
- **WebSocket Hub**: Manages real-time connections between users and admins
- **Support Ticket System**: Stores and manages support conversations
- **Routes**: 
  - User routes: `/api/app/support/*`
  - Admin WebSocket: `/admin/chat/ws`
  - Chat history: `/admin/chat/history`

### Flutter App
- **WebSocket Service**: Handles real-time communication
- **Support Service**: Manages ticket creation and HTTP fallback
- **UI Components**: 
  - CreateSupportTicketDialog: Form to initiate support
  - LiveSupportChatDialog: Real-time chat interface
  - Updated ChatbotDialog: Entry point with live support option

## Implementation Details

### 1. Backend WebSocket System

**Location**: `/backend/internal/ws/hub.go`

The WebSocket hub manages connections using a ticket-based room system:
- Each support ticket has a unique room
- Both users and admins can connect to the same room
- Messages are broadcasted to all participants in the room

**Connection URL Format**:
```
ws://[host]/admin/chat/ws?ticket_id=[TICKET_ID]&user_id=[USER_ID]&user_name=[USER_NAME]
```

### 2. Flutter WebSocket Service

**Location**: `/app/lib/core/services/websocket_service.dart`

Features:
- Automatic reconnection with exponential backoff
- Connection state management
- Message streaming via Dart streams
- Graceful disconnection handling

**Usage**:
```dart
final wsService = WebSocketService();
await wsService.connect(ticketId);

wsService.messageStream.listen((message) {
  // Handle incoming message
});

wsService.sendMessage('Hello', ticketId);
```

### 3. Support Ticket Service

**Location**: `/app/lib/core/services/support_service.dart`

Provides HTTP API access for:
- Creating support tickets
- Fetching ticket history
- Getting chat history
- HTTP fallback for replies when WebSocket unavailable

### 4. User Flow

1. **User opens chatbot** → ChatbotDialog displays
2. **User clicks "Live Support Chat"** → CreateSupportTicketDialog opens
3. **User fills form** (category, subject, description)
4. **Ticket created** → Backend generates unique ticket ID
5. **WebSocket connection established** → LiveSupportChatDialog opens
6. **Real-time chat begins** → Messages sync instantly with admin dashboard

### 5. Admin Flow (Existing)

1. Admin opens support page
2. Sees all tickets in list
3. Clicks on ticket to view details
4. WebSocket connects to same ticket room
5. Both user and admin messages appear in real-time

## API Endpoints

### App Endpoints (Protected - Requires Auth)

```
POST   /api/app/support/tickets           - Create new support ticket
GET    /api/app/support/tickets           - Get user's tickets
GET    /api/app/support/tickets/:ticketId - Get ticket details
POST   /api/app/support/tickets/:ticketId/reply - Add reply (HTTP fallback)
```

### Admin Endpoints (Existing)

```
GET    /admin/chat/ws                     - WebSocket connection
GET    /admin/chat/history?ticket_id=X    - Get chat history
GET    /admin/support                     - Support dashboard page
```

## Database Schema

**Collection**: `support_tickets`

```javascript
{
  "_id": ObjectId,
  "ticket_id": "ABC123",           // 6-char short ID for display
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "category": "General Inquiry",
  "subject": "Need help with booking",
  "description": "Detailed issue description",
  "status": "Open",                // Open, In Progress, Resolved, Closed
  "priority": "Medium",             // Low, Medium, High
  "replies": [
    {
      "_id": ObjectId,
      "message": "Message text",
      "is_admin": false,
      "read_by_admin": false,
      "user_name": "John Doe",
      "created_at": ISODate
    }
  ],
  "created_at": ISODate,
  "updated_at": ISODate
}
```

## Message Flow

### User Sends Message

1. User types message in LiveSupportChatDialog
2. Message added to local UI (optimistic update)
3. WebSocket sends message to backend
4. Backend saves to database
5. Backend broadcasts to all room participants
6. Admin receives message in real-time

### Admin Sends Message

1. Admin types in support dashboard chat
2. WebSocket sends to backend
3. Backend saves with `is_admin: true`
4. Backend broadcasts to room
5. User receives in LiveSupportChatDialog
6. Displays with admin badge/styling

## Configuration

### Flutter App - Update Base URL

**File**: `/app/lib/core/network/endpoints.dart`

```dart
static const String baseUrl = 'http://10.0.2.2:8080/api';  // Android emulator
// static const String baseUrl = 'http://localhost:8080/api';  // iOS simulator
// static const String baseUrl = 'http://192.168.1.x:8080/api';  // Physical device
```

### Dependencies

**File**: `/app/pubspec.yaml`

```yaml
dependencies:
  web_socket_channel: ^2.4.0  # WebSocket support
  dio: ^5.4.0                 # HTTP client
  get: ^4.6.6                 # State management
```

Run: `flutter pub get`

## Security Considerations

1. **Authentication**: All endpoints require JWT authentication
2. **Authorization**: Users can only access their own tickets
3. **WebSocket Security**: User identity verified via query parameters
4. **Input Validation**: All user inputs sanitized on backend
5. **Rate Limiting**: Applied on backend routes (already configured)

## Error Handling

### WebSocket Connection Failures
- Automatic reconnection (max 5 attempts)
- Exponential backoff (2s, 4s, 6s, 8s, 10s)
- Fallback to HTTP for sending messages
- User-friendly error notifications

### Network Issues
- Loading states during API calls
- Error dialogs with actionable messages
- Graceful degradation to HTTP-only mode

## Testing

### Test User Flow

1. Login to Flutter app with test account
2. Open chatbot from main screen
3. Click "Live Support Chat"
4. Fill in ticket details
5. Verify ticket creation
6. Send test message
7. Check message appears

### Test Admin Flow

1. Login to admin dashboard
2. Navigate to Support page
3. Find newly created ticket
4. Click to open chat
5. Send reply
6. Verify user receives in real-time

## Troubleshooting

### WebSocket Not Connecting

**Issue**: Connection fails or immediately disconnects

**Solutions**:
1. Check backend server is running
2. Verify WebSocket URL format
3. Check firewall/network settings
4. Ensure user is authenticated
5. Check backend logs for errors

### Messages Not Appearing

**Issue**: Messages sent but not received

**Solutions**:
1. Verify both parties connected to same ticket_id
2. Check message format in WebSocket payload
3. Inspect backend broadcast logic
4. Check client message stream subscription
5. Verify database writes successful

### Ticket Creation Fails

**Issue**: Cannot create support ticket

**Solutions**:
1. Check user authentication token
2. Verify all required fields filled
3. Check backend validation rules
4. Inspect network request/response
5. Check database connection

## File Structure

```
app/lib/
├── core/
│   ├── network/
│   │   ├── api_client.dart
│   │   └── endpoints.dart (updated)
│   └── services/
│       ├── websocket_service.dart (new)
│       └── support_service.dart (new)
└── features/
    └── chatbot/
        └── view/
            ├── chatbot_dialog.dart (updated)
            ├── create_support_ticket_dialog.dart (new)
            └── live_support_chat_dialog.dart (new)

backend/
├── internal/
│   ├── handlers/
│   │   ├── admin/
│   │   │   ├── chat_ws.go
│   │   │   └── support_handler.go
│   │   └── app/
│   │       └── support_handler.go (new)
│   ├── repository/
│   │   └── app/
│   │       └── support_ticket_repository.go (updated)
│   ├── routes/
│   │   └── app/
│   │       └── routes.go (updated)
│   └── ws/
│       └── hub.go
└── main.go (updated)
```

## Next Steps / Enhancements

1. **Push Notifications**: Notify users of admin replies even when app closed
2. **File Attachments**: Allow users to send images/documents
3. **Typing Indicators**: Show when admin is typing
4. **Read Receipts**: Show when admin has read messages
5. **Ticket History**: Display past closed tickets in app
6. **AI Auto-Response**: Suggest answers before escalating to human agent
7. **Ticket Ratings**: Allow users to rate support experience
8. **Canned Responses**: Quick replies for common questions

## Support

For issues or questions, please refer to:
- Backend API Documentation: `/backend/API_DOCUMENTATION.md`
- WebSocket Testing: Check admin dashboard support page
- Flutter Debug: Enable verbose logging in WebSocket service
