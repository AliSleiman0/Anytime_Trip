import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;
import '../network/endpoints.dart';
import '../storage/storage_service.dart';

class WebSocketService {
  WebSocketChannel? _channel;
  final StreamController<Map<String, dynamic>> _messageController =
      StreamController<Map<String, dynamic>>.broadcast();
  final StorageService _storage = StorageService();
  
  bool _isConnected = false;
  String? _currentTicketId;
  Timer? _reconnectTimer;
  int _reconnectAttempts = 0;
  static const int maxReconnectAttempts = 5;

  static final WebSocketService _instance = WebSocketService._internal();
  factory WebSocketService() => _instance;
  WebSocketService._internal();

  bool get isConnected => _isConnected;
  Stream<Map<String, dynamic>> get messageStream => _messageController.stream;

  /// Connect to WebSocket for a specific ticket
  Future<void> connect(String ticketId) async {
    if (_isConnected && _currentTicketId == ticketId) {
      return; // Already connected to this ticket
    }

    await disconnect();
    
    try {
      final user = await _storage.getUser();
      if (user == null) {
        throw Exception('User not logged in');
      }

      print('[WebSocket] User data: $user');
      
      final userId = user['user_id'] ?? user['id'];
      final userName = user['full_name'] ?? user['name'] ?? 'User';

      print('[WebSocket] userId: $userId, userName: $userName');

      // Get WebSocket URL from base URL
      final baseUrl = Endpoints.baseUrl;
      final wsUrl = baseUrl
          .replaceAll('http://', 'ws://')
          .replaceAll('https://', 'wss://')
          .replaceAll('/api', '');

      final url = '$wsUrl/admin/chat/ws?ticket_id=$ticketId&user_id=$userId&user_name=${Uri.encodeComponent(userName)}';

      print('[WebSocket] Connecting to: $url');
      
      _channel = WebSocketChannel.connect(Uri.parse(url));
      _currentTicketId = ticketId;
      _isConnected = true;
      _reconnectAttempts = 0;

      // Listen to messages
      _channel!.stream.listen(
        _onMessage,
        onError: _onError,
        onDone: _onDone,
        cancelOnError: false,
      );

      print('[WebSocket] Connected successfully');
    } catch (e) {
      print('[WebSocket] Connection error: $e');
      _isConnected = false;
      _attemptReconnect();
      rethrow;
    }
  }

  /// Disconnect from WebSocket
  Future<void> disconnect() async {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    
    if (_channel != null) {
      await _channel!.sink.close(status.goingAway);
      _channel = null;
    }
    
    _isConnected = false;
    _currentTicketId = null;
    print('[WebSocket] Disconnected');
  }

  /// Send a message through WebSocket
  void sendMessage(String message, String ticketId) {
    if (!_isConnected || _channel == null) {
      print('[WebSocket] Cannot send message: not connected');
      throw Exception('WebSocket not connected');
    }

    final payload = {
      'message': message,
      'ticket_id': ticketId,
      'timestamp': DateTime.now().toIso8601String(),
    };

    final jsonMessage = json.encode(payload);
    print('[WebSocket] Sending message: $jsonMessage');
    _channel!.sink.add(jsonMessage);
    print('[WebSocket] Message sent successfully');
  }

  /// Handle incoming messages
  void _onMessage(dynamic data) {
    try {
      print('[WebSocket] Raw message received: $data');
      final message = json.decode(data.toString()) as Map<String, dynamic>;
      print('[WebSocket] Parsed message: $message');
      _messageController.add(message);
    } catch (e) {
      print('[WebSocket] Error parsing message: $e');
    }
  }

  /// Handle WebSocket errors
  void _onError(error) {
    print('[WebSocket] Error: $error');
    _isConnected = false;
    _attemptReconnect();
  }

  /// Handle WebSocket connection closure
  void _onDone() {
    print('[WebSocket] Connection closed');
    _isConnected = false;
    _attemptReconnect();
  }

  /// Attempt to reconnect
  void _attemptReconnect() {
    if (_reconnectAttempts >= maxReconnectAttempts) {
      print('[WebSocket] Max reconnect attempts reached');
      return;
    }

    if (_currentTicketId == null) {
      return; // Don't reconnect if no ticket ID
    }

    _reconnectAttempts++;
    final delay = Duration(seconds: _reconnectAttempts * 2);
    
    print('[WebSocket] Reconnecting in ${delay.inSeconds} seconds (attempt $_reconnectAttempts)');
    
    _reconnectTimer = Timer(delay, () {
      if (_currentTicketId != null) {
        connect(_currentTicketId!);
      }
    });
  }

  /// Dispose resources
  void dispose() {
    disconnect();
    _messageController.close();
  }
}
