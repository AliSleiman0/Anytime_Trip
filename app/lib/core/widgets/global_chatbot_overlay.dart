import 'package:flutter/material.dart';
import '../../features/chatbot/view/chatbot_button.dart';

/// Global notifier to track if a modal/dropdown is open
final ValueNotifier<bool> isModalOpenNotifier = ValueNotifier<bool>(false);

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
            // Listen to modal state and conditionally show chatbot
            ValueListenableBuilder<bool>(
              valueListenable: isModalOpenNotifier,
              builder: (context, isModalOpen, _) {
                return isModalOpen
                    ? SizedBox.shrink() // Hide chatbot when modal is open
                    : const ChatbotButton();
              },
            ),
          ],
        );
      },
    );
  }
}
