import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'chatbot_dialog.dart';

class ChatbotButton extends StatelessWidget {
  final VoidCallback? onPressed;

  const ChatbotButton({
    super.key,
    this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return Positioned(
      bottom: 20,
      right: 20,
      child: GestureDetector(
        onTap: () {
          print("Chatbot button tapped!"); // Debug log
          onPressed?.call();
          Get.dialog(
            const ChatbotDialog(),
            barrierDismissible: true,
          );
        },
        child: Container(
          width: 60,
          height: 60,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: Colors.white,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.3),
                blurRadius: 10,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: ClipOval(
            child: Image.asset(
              'assets/images/chatbot.png',
              fit: BoxFit.cover,
              errorBuilder: (context, error, stackTrace) {
                return const Icon(
                  Icons.chat_bubble,
                  color: Color(0xFF336891),
                  size: 30,
                );
              },
            ),
          ),
        ),
      ),
    );
  }
}
