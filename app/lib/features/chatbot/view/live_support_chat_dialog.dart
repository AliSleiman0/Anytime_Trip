import 'dart:async';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../core/services/support_service.dart';
import '../../../core/services/websocket_service.dart';

class LiveSupportChatDialog extends StatefulWidget {
  final String ticketId;
  final String subject;

  const LiveSupportChatDialog({
    super.key,
    required this.ticketId,
    required this.subject,
  });

  @override
  State<LiveSupportChatDialog> createState() => _LiveSupportChatDialogState();
}

class _LiveSupportChatDialogState extends State<LiveSupportChatDialog> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  final List<ChatMessage> _messages = [];
  final WebSocketService _wsService = WebSocketService();
  final SupportService _supportService = SupportService();
  
  bool _isConnected = false;
  bool _isLoading = true;
  StreamSubscription? _messageSubscription;

  @override
  void initState() {
    super.initState();
    _initializeChat();
  }

  Future<void> _initializeChat() async {
    setState(() => _isLoading = true);

    try {
      print('[Chat] Loading history for ticket: ${widget.ticketId}');
      // Load chat history
      final history = await _supportService.getChatHistory(widget.ticketId);
      print('[Chat] History received: $history');
      
      // Parse ticket data safely
      final ticketData = history['ticket'];
      if (ticketData == null) {
        throw Exception('No ticket data in response');
      }
      
      print('[Chat] Parsing ticket data...');
      final ticket = SupportTicket.fromJson(ticketData as Map<String, dynamic>);
      print('[Chat] Ticket parsed, ${ticket.replies.length} replies');
      
      // Convert replies to chat messages
      if (ticket.replies.isNotEmpty) {
        for (var reply in ticket.replies) {
          _messages.add(ChatMessage(
            text: reply.message,
            isBot: reply.isAdmin,
            senderName: reply.isAdmin ? (reply.adminName ?? 'Admin') : 'You',
            timestamp: _formatTime(reply.createdAt),
          ));
        }
      }

      print('[Chat] Connecting to WebSocket...');
      // Connect to WebSocket
      await _wsService.connect(widget.ticketId);
      print('[Chat] WebSocket connected: ${_wsService.isConnected}');
      
      setState(() {
        _isConnected = _wsService.isConnected;
        _isLoading = false;
      });

      // Listen to incoming messages
      _messageSubscription = _wsService.messageStream.listen((message) {
        _handleIncomingMessage(message);
      });

      // Scroll to bottom
      _scrollToBottom();
    } catch (e, stackTrace) {
      print('[Chat] Error initializing chat: $e');
      print('[Chat] Stack trace: $stackTrace');
      setState(() => _isLoading = false);
      Get.snackbar(
        'Error',
        'Failed to connect to support chat: $e',
        snackPosition: SnackPosition.BOTTOM,
        backgroundColor: Colors.red,
        colorText: Colors.white,
      );
    }
  }

  void _handleIncomingMessage(Map<String, dynamic> message) {
    print('[Chat] Incoming message: $message');
    final isAdmin = message['is_admin'] as bool? ?? false;
    final messageText = message['message'] as String? ?? '';
    final senderName = isAdmin 
        ? (message['admin_name'] as String? ?? 'Admin')
        : 'You';
    
    // Only add admin messages to UI (user's own messages already added optimistically)
    if (messageText.isNotEmpty && isAdmin) {
      print('[Chat] Adding admin message to UI: $messageText');
      setState(() {
        _messages.add(ChatMessage(
          text: messageText,
          isBot: isAdmin,
          senderName: senderName,
          timestamp: _formatTime(DateTime.now()),
        ));
      });
      _scrollToBottom();
    } else {
      print('[Chat] Skipping user message (already in UI)');
    }
  }

  String _formatTime(DateTime time) {
    final hour = time.hour > 12 ? time.hour - 12 : (time.hour == 0 ? 12 : time.hour);
    final minute = time.minute.toString().padLeft(2, '0');
    final period = time.hour >= 12 ? 'pm' : 'am';
    return '$hour:$minute $period';
  }

  void _sendMessage() {
    if (_messageController.text.isEmpty) return;

    final messageText = _messageController.text;
    
    try {
      print('[Chat] Sending message: $messageText');
      
      // Clear input immediately
      _messageController.clear();

      // Add message to UI immediately (optimistic update)
      setState(() {
        _messages.add(ChatMessage(
          text: messageText,
          isBot: false,
          senderName: 'You',
          timestamp: _formatTime(DateTime.now()),
        ));
      });

      // Send through WebSocket
      if (_isConnected) {
        print('[Chat] Sending via WebSocket');
        _wsService.sendMessage(messageText, widget.ticketId);
      } else {
        print('[Chat] Sending via HTTP (WebSocket not connected)');
        // Fallback to HTTP if WebSocket not connected
        _supportService.replyToTicket(widget.ticketId, messageText);
      }

      _scrollToBottom();
    } catch (e) {
      print('Error sending message: $e');
      Get.snackbar(
        'Error',
        'Failed to send message',
        snackPosition: SnackPosition.BOTTOM,
        backgroundColor: Colors.red,
        colorText: Colors.white,
      );
    }
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  void dispose() {
    _messageSubscription?.cancel();
    _wsService.disconnect();
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      insetPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 24),
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
        ),
        constraints: BoxConstraints(
          maxHeight: MediaQuery.of(context).size.height * 0.8,
          maxWidth: 400,
        ),
        child: Column(
          children: [
            // Header
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
              decoration: const BoxDecoration(
                color: Color(0xFF336891),
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(16),
                  topRight: Radius.circular(16),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Row(
                      children: [
                        const Icon(
                          Icons.support_agent,
                          color: Colors.white,
                          size: 32,
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Text(
                                'Live Support',
                                style: TextStyle(
                                  color: Colors.white,
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                              Row(
                                children: [
                                  Container(
                                    width: 8,
                                    height: 8,
                                    decoration: BoxDecoration(
                                      shape: BoxShape.circle,
                                      color: _isConnected ? Colors.green : Colors.grey,
                                    ),
                                  ),
                                  const SizedBox(width: 6),
                                  Text(
                                    _isConnected ? 'Connected' : 'Offline',
                                    style: const TextStyle(
                                      color: Colors.white70,
                                      fontSize: 12,
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                  GestureDetector(
                    onTap: () => Get.back(),
                    child: Container(
                      width: 32,
                      height: 32,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: Colors.white.withOpacity(0.2),
                      ),
                      child: const Icon(
                        Icons.close,
                        color: Colors.white,
                        size: 18,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            
            // Ticket info
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                border: Border(
                  bottom: BorderSide(color: Colors.grey.shade200),
                ),
              ),
              child: Row(
                children: [
                  const Icon(Icons.confirmation_number, size: 16, color: Colors.grey),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      '${widget.ticketId} - ${widget.subject}',
                      style: const TextStyle(
                        fontSize: 12,
                        color: Colors.grey,
                        fontWeight: FontWeight.w500,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
            ),

            // Messages
            if (_isLoading)
              const Expanded(
                child: Center(
                  child: CircularProgressIndicator(
                    color: Color(0xFF336891),
                  ),
                ),
              )
            else
              Expanded(
                child: ListView.builder(
                  controller: _scrollController,
                  padding: const EdgeInsets.all(16),
                  itemCount: _messages.length,
                  itemBuilder: (context, index) {
                    final message = _messages[index];
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 16),
                      child: Column(
                        crossAxisAlignment: message.isBot
                            ? CrossAxisAlignment.start
                            : CrossAxisAlignment.end,
                        children: [
                          Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: message.isBot
                                ? MainAxisAlignment.start
                                : MainAxisAlignment.end,
                            children: [
                              if (message.isBot)
                                Container(
                                  width: 32,
                                  height: 32,
                                  margin: const EdgeInsets.only(right: 12),
                                  decoration: const BoxDecoration(
                                    shape: BoxShape.circle,
                                    color: Color(0xFF336891),
                                  ),
                                  child: const Icon(
                                    Icons.support_agent,
                                    color: Colors.white,
                                    size: 20,
                                  ),
                                ),
                              Flexible(
                                child: Column(
                                  crossAxisAlignment: message.isBot
                                      ? CrossAxisAlignment.start
                                      : CrossAxisAlignment.end,
                                  children: [
                                    Text(
                                      message.senderName,
                                      style: TextStyle(
                                        fontSize: 11,
                                        fontWeight: FontWeight.w600,
                                        color: Colors.grey.shade600,
                                      ),
                                    ),
                                    const SizedBox(height: 4),
                                    Container(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 12,
                                        vertical: 10,
                                      ),
                                      decoration: BoxDecoration(
                                        color: message.isBot
                                            ? const Color(0xFF336891)
                                            : const Color(0xFFE8F0F7),
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: Text(
                                        message.text,
                                        style: TextStyle(
                                          color: message.isBot
                                              ? Colors.white
                                              : Colors.black87,
                                          fontSize: 13,
                                          height: 1.4,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 6),
                          Text(
                            message.timestamp,
                            style: const TextStyle(
                              color: Colors.grey,
                              fontSize: 11,
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),

            // Input Area
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                border: Border(
                  top: BorderSide(
                    color: Colors.grey.shade200,
                    width: 1,
                  ),
                ),
              ),
              child: Row(
                children: [
                  Expanded(
                    child: Container(
                      decoration: BoxDecoration(
                        color: const Color(0xFFE8F0F7),
                        borderRadius: BorderRadius.circular(24),
                      ),
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      child: TextField(
                        controller: _messageController,
                        decoration: const InputDecoration(
                          border: InputBorder.none,
                          hintText: 'Type your message...',
                          hintStyle: TextStyle(
                            color: Color(0xFF999999),
                            fontSize: 13,
                          ),
                        ),
                        enabled: !_isLoading,
                        onSubmitted: (_) => _sendMessage(),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  GestureDetector(
                    onTap: _isLoading ? null : _sendMessage,
                    child: Container(
                      width: 44,
                      height: 44,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: _isLoading 
                            ? Colors.grey 
                            : const Color(0xFF336891),
                      ),
                      child: const Icon(
                        Icons.send,
                        color: Colors.white,
                        size: 20,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ChatMessage {
  final String text;
  final bool isBot;
  final String senderName;
  final String timestamp;

  ChatMessage({
    required this.text,
    required this.isBot,
    required this.senderName,
    required this.timestamp,
  });
}
