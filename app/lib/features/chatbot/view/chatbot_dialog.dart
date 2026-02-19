import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'chatbot_chat_dialog.dart';
import 'create_support_ticket_dialog.dart';
import '../../support/view/support_tickets_screen.dart';

class ChatbotDialog extends StatelessWidget {
  const ChatbotDialog({super.key});

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(20),
      ),
      insetPadding: EdgeInsets.zero,
      child: Container(
        width: 340,
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // Header with Brain Icon
            ClipRRect(
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(20),
                topRight: Radius.circular(20),
              ),
              child: Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(vertical: 20),
                decoration: const BoxDecoration(
                  color: Color(0xFF336891),
                ),
              child: Column(
                children: [
                  Image.asset(
                    'assets/images/brain.png',
                    width: 80,
                    height: 80,
                  ),
                  const SizedBox(height: 12),
                  const Text(
                    'Help & Support',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    'Our Chatbot is here to assist you',
                    style: TextStyle(
                      color: Colors.white70,
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),            ),            // Content
            Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'What do you wanna know?',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: Colors.black87,
                    ),
                  ),
                  const SizedBox(height: 16),
                  // FAQ Items
                  _buildFaqItem(
                    context,
                    'How can I cancel or modify my booking?',
                    'You can cancel or modify your booking through the "My Bookings" section in the app. For flights and hotels, free cancellation is available up to 24 hours before departure. Car rentals and transfers can be cancelled up to 48 hours in advance without any charges.',
                    1,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'What payment methods do you accept?',
                    'We accept all major credit cards (Visa, Mastercard, American Express), debit cards, and digital wallets. All payments are securely processed and encrypted. You can save your payment method for faster checkout on future bookings.',
                    2,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'How long until I receive my confirmation email?',
                    'Confirmation emails are typically sent within 5-10 minutes after completing your booking. Please check your spam folder if you don\'t see it in your inbox. You can also view all your bookings in the app under "My Bookings" section.',
                    3,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'How do I track my booking status?',
                    'Go to the "My Bookings" section in the app to see all your reservations. Each booking shows its current status (Confirmed, Pending, Completed, or Cancelled). You\'ll also receive email and in-app notifications for any status updates.',
                    1,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'What documents do I need for travel?',
                    'For international flights, you need a valid passport (with at least 6 months validity) and any required visas for your destination. For domestic flights, a government-issued ID is sufficient. Check your destination\'s entry requirements as they may vary.',
                    2,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'Is travel insurance included in my booking?',
                    'Travel insurance is optional and can be added during checkout. We highly recommend it for international trips. Our insurance covers trip cancellations, medical emergencies, lost baggage, and flight delays. You can add it to existing bookings within 24 hours.',
                    3,
                  ),
                  const SizedBox(height: 12),
                  _buildFaqItem(
                    context,
                    'Can I get a refund if I cancel my booking?',
                    'Refund policies vary by booking type. Flights and hotels with "Free Cancellation" offer full refunds if cancelled within the allowed timeframe. Non-refundable bookings may incur cancellation fees. Check your booking details for specific terms and conditions.',
                    1,
                  ),
                  const SizedBox(height: 20),
                  // Talk with Chatbot
                  GestureDetector(
                    onTap: () {
                      Get.back();
                      Get.dialog(
                        const ChatbotChatDialog(),
                        barrierDismissible: true,
                      );
                    },
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        color: const Color(0xFFF5F5F5),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Talk with Chatbot',
                                  style: TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black87,
                                  ),
                                ),
                                SizedBox(height: 4),
                                Text(
                                  'Our Chatbot will respond immediately',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.grey,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Image.asset(
                            'assets/images/sent.png',
                            width: 40,
                            height: 40,
                            fit: BoxFit.contain,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  // Live Support
                  GestureDetector(
                    onTap: () {
                      Get.back();
                      Get.dialog(
                        const CreateSupportTicketDialog(),
                        barrierDismissible: true,
                      );
                    },
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        color: const Color(0xFF336891),
                        borderRadius: BorderRadius.circular(12),
                        boxShadow: [
                          BoxShadow(
                            color: const Color(0xFF336891).withOpacity(0.3),
                            blurRadius: 8,
                            offset: const Offset(0, 4),
                          ),
                        ],
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  children: [
                                    Icon(
                                      Icons.support_agent,
                                      color: Colors.white,
                                      size: 20,
                                    ),
                                    SizedBox(width: 8),
                                    Text(
                                      'Live Support Chat',
                                      style: TextStyle(
                                        fontSize: 14,
                                        fontWeight: FontWeight.w600,
                                        color: Colors.white,
                                      ),
                                    ),
                                  ],
                                ),
                                SizedBox(height: 4),
                                Text(
                                  'Connect with a real support agent',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.white70,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const Icon(
                            Icons.arrow_forward,
                            color: Colors.white,
                            size: 24,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  // My Tickets
                  GestureDetector(
                    onTap: () {
                      Get.back();
                      Get.to(() => const SupportTicketsScreen());
                    },
                    child: Container(
                      padding: const EdgeInsets.all(14),
                      decoration: BoxDecoration(
                        border: Border.all(
                          color: const Color(0xFF336891),
                          width: 2,
                        ),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  children: [
                                    Icon(
                                      Icons.history,
                                      color: Color(0xFF336891),
                                      size: 20,
                                    ),
                                    SizedBox(width: 8),
                                    Text(
                                      'My Support Tickets',
                                      style: TextStyle(
                                        fontSize: 14,
                                        fontWeight: FontWeight.w600,
                                        color: Color(0xFF336891),
                                      ),
                                    ),
                                  ],
                                ),
                                SizedBox(height: 4),
                                Text(
                                  'View your ticket history and replies',
                                  style: TextStyle(
                                    fontSize: 12,
                                    color: Colors.grey,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const Icon(
                            Icons.arrow_forward,
                            color: Color(0xFF336891),
                            size: 24,
                          ),
                        ],
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

  Widget _buildFaqItem(BuildContext context, String question, String answer, int itemNumber) {
    return GestureDetector(
      onTap: () {
        // Show FAQ answer in a simple dialog
        Get.dialog(
          Dialog(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(16),
            ),
            child: Container(
              padding: const EdgeInsets.all(20),
              constraints: const BoxConstraints(maxWidth: 340),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Image.asset(
                        'assets/images/brain.png',
                        width: 32,
                        height: 32,
                      ),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Text(
                          'FAQ Answer',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                            color: Color(0xFF336891),
                          ),
                        ),
                      ),
                      GestureDetector(
                        onTap: () => Get.back(),
                        child: Container(
                          width: 28,
                          height: 28,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: Colors.grey.shade200,
                          ),
                          child: const Icon(
                            Icons.close,
                            size: 16,
                            color: Colors.black54,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xFFE8F0F7),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      question,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: Color(0xFF336891),
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    answer,
                    style: const TextStyle(
                      fontSize: 13,
                      height: 1.5,
                      color: Colors.black87,
                    ),
                  ),
                  const SizedBox(height: 20),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () => Get.back(),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF336891),
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text(
                        'Got it',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
          barrierDismissible: true,
        );
      },
      child: Row(
        children: [
          Image.asset(
            'assets/images/$itemNumber.png',
            width: 24,
            height: 24,
            fit: BoxFit.contain,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              question,
              style: const TextStyle(
                fontSize: 13,
                color: Color(0xFF336891),
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
