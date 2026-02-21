import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'chatbot_rating_dialog.dart';

class ChatbotChatDialog extends StatefulWidget {
  const ChatbotChatDialog({super.key});

  @override
  State<ChatbotChatDialog> createState() => _ChatbotChatDialogState();
}

class _ChatbotChatDialogState extends State<ChatbotChatDialog> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  late final List<ChatMessage> _messages;
  
  // Predefined FAQ questions and answers
  final List<Map<String, String>> _faqItems = [
    {
      'question': 'How can I cancel or modify my booking?',
      'answer': 'You can cancel or modify your booking through the "My Bookings" section in the app. For flights and hotels, free cancellation is available up to 24 hours before departure. Car rentals and transfers can be cancelled up to 48 hours in advance without any charges.',
    },
    {
      'question': 'What payment methods do you accept?',
      'answer': 'We accept all major credit cards (Visa, Mastercard, American Express), debit cards, and digital wallets. All payments are securely processed and encrypted. You can save your payment method for faster checkout on future bookings.',
    },
    {
      'question': 'How long until I receive my confirmation email?',
      'answer': 'Confirmation emails are typically sent within 5-10 minutes after completing your booking. Please check your spam folder if you don\'t see it in your inbox. You can also view all your bookings in the app under "My Bookings" section.',
    },
    {
      'question': 'How do I track my booking status?',
      'answer': 'Go to the "My Bookings" section in the app to see all your reservations. Each booking shows its current status (Confirmed, Pending, Completed, or Cancelled). You\'ll also receive email and in-app notifications for any status updates.',
    },
    {
      'question': 'What documents do I need for travel?',
      'answer': 'For international flights, you need a valid passport (with at least 6 months validity) and any required visas for your destination. For domestic flights, a government-issued ID is sufficient. Check your destination\'s entry requirements as they may vary.',
    },
    {
      'question': 'Is travel insurance included in my booking?',
      'answer': 'Travel insurance is optional and can be added during checkout. We highly recommend it for international trips. Our insurance covers trip cancellations, medical emergencies, lost baggage, and flight delays. You can add it to existing bookings within 24 hours.',
    },
    {
      'question': 'Can I get a refund if I cancel my booking?',
      'answer': 'Refund policies vary by booking type. Flights and hotels with "Free Cancellation" offer full refunds if cancelled within the allowed timeframe. Non-refundable bookings may incur cancellation fees. Check your booking details for specific terms and conditions.',
    },
  ];

  @override
  void initState() {
    super.initState();
    _messages = [
      ChatMessage(
        text: _buildWelcomeMessage(),
        isBot: true,
        timestamp: _formatTime(DateTime.now()),
      ),
    ];
  }

  String _buildWelcomeMessage() {
    String message = "Hello! I'm here to help you. Please select a question by typing its number:\n\n";
    for (int i = 0; i < _faqItems.length; i++) {
      message += "${i + 1}. ${_faqItems[i]['question']}\n";
    }
    message += "\nOr type your own question and I'll try to help!";
    return message;
  }

  String _formatTime(DateTime time) {
    final hour = time.hour > 12 ? time.hour - 12 : (time.hour == 0 ? 12 : time.hour);
    final minute = time.minute.toString().padLeft(2, '0');
    final period = time.hour >= 12 ? 'pm' : 'am';
    return '$hour:$minute $period';
  }

  void _scrollToBottom() {
    Future.delayed(const Duration(milliseconds: 100), () {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _sendMessage() {
    if (_messageController.text.isEmpty) return;

    final userInput = _messageController.text.trim();
    
    setState(() {
      _messages.add(ChatMessage(
        text: userInput,
        isBot: false,
        timestamp: _formatTime(DateTime.now()),
      ));
      _messageController.clear();
    });
    
    _scrollToBottom();

    // Check if user typed a number corresponding to FAQ
    final number = int.tryParse(userInput);
    String botResponse;
    
    if (number != null && number >= 1 && number <= _faqItems.length) {
      // User selected a FAQ question
      final selectedFaq = _faqItems[number - 1];
      botResponse = "📌 ${selectedFaq['question']}\n\n${selectedFaq['answer']}\n\n---\nWould you like to ask another question? Type a number (1-${_faqItems.length}) or your own question.";
    } else {
      // Custom question - provide a generic response
      botResponse = "Thank you for your question! Our support team will assist you shortly. In the meantime, you can select from our FAQ by typing a number:\n\n";
      for (int i = 0; i < _faqItems.length; i++) {
        botResponse += "${i + 1}. ${_faqItems[i]['question']}\n";
      }
    }

    // Simulate bot response with delay
    Future.delayed(const Duration(milliseconds: 500), () {
      if (mounted) {
        setState(() {
          _messages.add(ChatMessage(
            text: botResponse,
            isBot: true,
            timestamp: _formatTime(DateTime.now()),
          ));
        });
        _scrollToBottom();
      }
    });
  }

  @override
  void dispose() {
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
                  Row(
                    children: [
                      Image.asset(
                        'assets/images/brain.png',
                        width: 32,
                        height: 32,
                      ),
                      const SizedBox(width: 12),
                      const Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Chatbot',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          Text(
                            '24/7 Support Bot',
                            style: TextStyle(
                              color: Colors.white70,
                              fontSize: 12,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  GestureDetector(
                    onTap: () {
                      Get.back();
                      Get.dialog(
                        const ChatbotRatingDialog(),
                        barrierDismissible: true,
                      );
                    },
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
            // Messages
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
                                  color: Color(0xFFE84C3D),
                                ),
                                child: Center(
                                  child: Image.asset(
                                    'assets/images/brain.png',
                                    width: 20,
                                    height: 20,
                                  ),
                                ),
                              ),
                            Flexible(
                              child: Container(
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
                        if (message.isBot && index == _messages.length - 1)
                          Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                GestureDetector(
                                  onTap: () {},
                                  child: Container(
                                    width: 32,
                                    height: 32,
                                    decoration: BoxDecoration(
                                      shape: BoxShape.circle,
                                      border: Border.all(
                                        color: const Color(0xFFE84C3D),
                                        width: 1.5,
                                      ),
                                    ),
                                    child: const Icon(
                                      Icons.thumb_up_outlined,
                                      size: 16,
                                      color: Color(0xFFE84C3D),
                                    ),
                                  ),
                                ),
                                const SizedBox(width: 8),
                                GestureDetector(
                                  onTap: () {},
                                  child: Container(
                                    width: 32,
                                    height: 32,
                                    decoration: BoxDecoration(
                                      shape: BoxShape.circle,
                                      border: Border.all(
                                        color: const Color(0xFFE84C3D),
                                        width: 1.5,
                                      ),
                                    ),
                                    child: const Icon(
                                      Icons.thumb_down_outlined,
                                      size: 16,
                                      color: Color(0xFFE84C3D),
                                    ),
                                  ),
                                ),
                              ],
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
                      child: Row(
                        children: [
                          Expanded(
                            child: TextField(
                              controller: _messageController,
                              decoration: const InputDecoration(
                                border: InputBorder.none,
                                hintText: 'Enter Text...',
                                hintStyle: TextStyle(
                                  color: Color(0xFF999999),
                                  fontSize: 13,
                                ),
                              ),
                              onSubmitted: (_) => _sendMessage(),
                            ),
                          ),
                          GestureDetector(
                            onTap: () {},
                            child: const Icon(
                              Icons.camera_alt_outlined,
                              color: Color(0xFF999999),
                              size: 20,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  GestureDetector(
                    onTap: _sendMessage,
                    child: Container(
                      width: 44,
                      height: 44,
                      decoration: const BoxDecoration(
                        shape: BoxShape.circle,
                        color: Color(0xFF336891),
                      ),
                      child: Image.asset(
                        'assets/images/sent.png',
                        fit: BoxFit.contain,
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
  final String timestamp;

  ChatMessage({
    required this.text,
    required this.isBot,
    required this.timestamp,
  });
}
