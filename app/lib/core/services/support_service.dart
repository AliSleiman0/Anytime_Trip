import 'package:dio/dio.dart';
import '../network/api_client.dart';
import '../network/endpoints.dart';
import '../storage/storage_service.dart';

class SupportTicket {
  final String? id;
  final String ticketId;
  final String customerName;
  final String customerEmail;
  final String category;
  final String subject;
  final String description;
  final String status;
  final String? priority;
  final List<TicketReply> replies;
  final DateTime createdAt;
  final DateTime updatedAt;

  SupportTicket({
    this.id,
    required this.ticketId,
    required this.customerName,
    required this.customerEmail,
    required this.category,
    required this.subject,
    required this.description,
    required this.status,
    this.priority,
    this.replies = const [],
    required this.createdAt,
    required this.updatedAt,
  });

  factory SupportTicket.fromJson(Map<String, dynamic> json) {
    try {
      print('Parsing SupportTicket from: $json');
      
      // Handle MongoDB ObjectID format
      String? idStr;
      if (json['id'] != null) {
        if (json['id'] is String) {
          idStr = json['id'] as String;
        } else if (json['id'] is Map && json['id']['\$oid'] != null) {
          idStr = json['id']['\$oid'] as String;
        }
      }

      // Parse replies safely
      List<TicketReply> replies = [];
      if (json['replies'] != null && json['replies'] is List) {
        for (var item in json['replies']) {
          if (item is Map<String, dynamic>) {
            try {
              replies.add(TicketReply.fromJson(item));
            } catch (e) {
              print('Error parsing reply: $e');
            }
          }
        }
      }

      return SupportTicket(
        id: idStr,
        ticketId: json['ticket_id']?.toString() ?? '',
        customerName: json['customer_name']?.toString() ?? '',
        customerEmail: json['customer_email']?.toString() ?? '',
        category: json['category']?.toString() ?? '',
        subject: json['subject']?.toString() ?? '',
        description: json['description']?.toString() ?? '',
        status: json['status']?.toString() ?? 'Open',
        priority: json['priority']?.toString(),
        replies: replies,
        createdAt: json['created_at'] != null 
            ? DateTime.parse(json['created_at'].toString())
            : DateTime.now(),
        updatedAt: json['updated_at'] != null
            ? DateTime.parse(json['updated_at'].toString())
            : DateTime.now(),
      );
    } catch (e) {
      print('Error in SupportTicket.fromJson: $e');
      print('JSON data: $json');
      rethrow;
    }
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'ticket_id': ticketId,
      'customer_name': customerName,
      'customer_email': customerEmail,
      'category': category,
      'subject': subject,
      'description': description,
      'status': status,
      'priority': priority,
      'replies': replies.map((r) => r.toJson()).toList(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class TicketReply {
  final String? id;
  final String message;
  final bool isAdmin;
  final String? adminName;
  final String? userName;
  final DateTime createdAt;

  TicketReply({
    this.id,
    required this.message,
    required this.isAdmin,
    this.adminName,
    this.userName,
    required this.createdAt,
  });

  factory TicketReply.fromJson(Map<String, dynamic> json) {
    try {
      // Handle MongoDB ObjectID format
      String? idStr;
      if (json['id'] != null) {
        if (json['id'] is String) {
          idStr = json['id'] as String;
        } else if (json['id'] is Map && json['id']['\$oid'] != null) {
          idStr = json['id']['\$oid'] as String;
        }
      }

      return TicketReply(
        id: idStr,
        message: json['message']?.toString() ?? '',
        isAdmin: json['is_admin'] == true,
        adminName: json['admin_name']?.toString(),
        userName: json['user_name']?.toString(),
        createdAt: json['created_at'] != null
            ? DateTime.parse(json['created_at'].toString())
            : DateTime.now(),
      );
    } catch (e) {
      print('Error in TicketReply.fromJson: $e');
      print('JSON data: $json');
      rethrow;
    }
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'message': message,
      'is_admin': isAdmin,
      'admin_name': adminName,
      'user_name': userName,
      'created_at': createdAt.toIso8601String(),
    };
  }
}

class SupportService {
  final ApiClient _apiClient = ApiClient();
  final StorageService _storage = StorageService();

  static final SupportService _instance = SupportService._internal();
  factory SupportService() => _instance;
  SupportService._internal();

  /// Create a new support ticket
  Future<SupportTicket> createTicket({
    required String category,
    required String subject,
    required String description,
  }) async {
    try {
      final user = await _storage.getUser();
      if (user == null) {
        throw Exception('User not logged in');
      }

      final response = await _apiClient.post(
        Endpoints.createSupportTicket,
        data: {
          'customer_name': user['full_name'] ?? user['name'] ?? 'User',
          'customer_email': user['email'] ?? '',
          'category': category,
          'subject': subject,
          'description': description,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        print('Create ticket response: ${response.data}');
        final ticketData = response.data['ticket'];
        if (ticketData == null) {
          throw Exception('No ticket data in response');
        }
        
        print('Ticket data: $ticketData');
        return SupportTicket.fromJson(ticketData as Map<String, dynamic>);
      } else {
        throw Exception('Failed to create ticket');
      }
    } on DioException catch (e) {
      print('DioException: ${e.response?.data}');
      if (e.response?.data != null && e.response!.data['error'] != null) {
        throw Exception(e.response!.data['error']);
      }
      throw Exception('Failed to create ticket: ${e.message}');
    } catch (e) {
      print('Error creating ticket: $e');
      rethrow;
    }
  }

  /// Get all user's support tickets
  Future<List<SupportTicket>> getUserTickets() async {
    try {
      print('[SupportService] Getting user tickets...');
      print('[SupportService] Endpoint: ${Endpoints.userSupportTickets}');
      
      // Check if user is logged in
      final token = await _storage.getToken();
      print('[SupportService] Has token: ${token != null}');
      if (token != null) {
        print('[SupportService] Token length: ${token.length}');
      }
      
      final response = await _apiClient.get(Endpoints.userSupportTickets);

      print('[SupportService] Tickets response status: ${response.statusCode}');
      print('[SupportService] Tickets response data: ${response.data}');

      if (response.statusCode == 200) {
        final data = response.data;
        if (data['tickets'] == null || data['tickets'] is! List) {
          print('[SupportService] No tickets found or invalid format');
          return [];
        }
        
        final tickets = <SupportTicket>[];
        for (var item in data['tickets']) {
          if (item is Map<String, dynamic>) {
            try {
              tickets.add(SupportTicket.fromJson(item));
            } catch (e) {
              print('[SupportService] Error parsing ticket: $e');
            }
          }
        }
        print('[SupportService] Parsed ${tickets.length} tickets');
        return tickets;
      } else {
        throw Exception('Failed to fetch tickets');
      }
    } on DioException catch (e) {
      print('[SupportService] DioException in getUserTickets:');
      print('[SupportService] - Type: ${e.type}');
      print('[SupportService] - Message: ${e.message}');
      print('[SupportService] - Error: ${e.error}');
      print('[SupportService] - Response status: ${e.response?.statusCode}');
      print('[SupportService] - Response data: ${e.response?.data}');
      
      if (e.response?.data != null && e.response!.data['error'] != null) {
        throw Exception(e.response!.data['error']);
      }
      throw Exception('Failed to fetch tickets: ${e.message ?? e.error?.toString() ?? "Unknown error"}');
    } catch (e) {
      print('[SupportService] Error in getUserTickets: $e');
      rethrow;
    }
  }

  /// Get ticket details by ticket ID
  Future<SupportTicket> getTicketDetails(String ticketId) async {
    try {
      final response = await _apiClient.get(
        Endpoints.supportTicketDetails(ticketId),
      );

      if (response.statusCode == 200) {
        return SupportTicket.fromJson(response.data);
      } else {
        throw Exception('Failed to fetch ticket details');
      }
    } on DioException catch (e) {
      if (e.response?.data != null && e.response!.data['error'] != null) {
        throw Exception(e.response!.data['error']);
      }
      throw Exception('Failed to fetch ticket details: ${e.message}');
    }
  }

  /// Get chat history for a ticket
  Future<Map<String, dynamic>> getChatHistory(String ticketId) async {
    try {
      print('[SupportService] Getting chat history for ticket: $ticketId');
      final response = await _apiClient.get(
        Endpoints.chatHistory(ticketId),
      );

      print('[SupportService] Chat history response status: ${response.statusCode}');
      print('[SupportService] Chat history response data type: ${response.data.runtimeType}');
      print('[SupportService] Chat history response data: ${response.data}');

      if (response.statusCode == 200) {
        // Ensure response.data is a Map
        if (response.data is Map<String, dynamic>) {
          return response.data as Map<String, dynamic>;
        } else {
          throw Exception('Invalid response format: ${response.data.runtimeType}');
        }
      } else {
        throw Exception('Failed to fetch chat history');
      }
    } on DioException catch (e) {
      print('[SupportService] DioException in getChatHistory: ${e.message}');
      print('[SupportService] Response: ${e.response?.data}');
      if (e.response?.data != null && e.response!.data['error'] != null) {
        throw Exception(e.response!.data['error']);
      }
      throw Exception('Failed to fetch chat history: ${e.message}');
    } catch (e) {
      print('Error in getChatHistory: $e');
      rethrow;
    }
  }

  /// Reply to a ticket (HTTP fallback if WebSocket not available)
  Future<void> replyToTicket(String ticketId, String message) async {
    try {
      final response = await _apiClient.post(
        Endpoints.replyToTicket(ticketId),
        data: {'message': message},
      );

      if (response.statusCode != 200 && response.statusCode != 201) {
        throw Exception('Failed to send reply');
      }
    } on DioException catch (e) {
      if (e.response?.data != null && e.response!.data['error'] != null) {
        throw Exception(e.response!.data['error']);
      }
      throw Exception('Failed to send reply: ${e.message}');
    }
  }
}
