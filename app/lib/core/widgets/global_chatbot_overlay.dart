import 'package:flutter/material.dart';
import '../../features/chatbot/view/chatbot_button.dart';

/// A global overlay widget that adds the chatbot button to any screen
/// This widget wraps the child and adds a positioned chatbot button
/// making it accessible from anywhere in the app
class GlobalChatbotOverlay extends StatelessWidget {
  final Widget child;

  const GlobalChatbotOverlay({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Builder(
      builder: (BuildContext context) {
        return Stack(
          children: [
            child,
            const ChatbotButton(),
          ],
        );
      },
    );
  }
}
