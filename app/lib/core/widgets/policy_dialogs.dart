import 'package:flutter/material.dart';
import 'package:get/get.dart';

class PolicyDialogs {
  /// Show Terms and Conditions Dialog
  static void showTermsAndConditions(BuildContext context) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Container(
            constraints: const BoxConstraints(maxWidth: 600),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Header
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: const BoxDecoration(
                    color: Color(0xFF1e5a8e),
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(16),
                      topRight: Radius.circular(16),
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'terms_conditions'.tr,
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close, color: Colors.white),
                        onPressed: () => Navigator.pop(context),
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                    ],
                  ),
                ),
                // Content
                Flexible(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildSectionTitle('1. Acceptance of Terms'),
                        _buildSectionContent(
                          'By accessing and using Travel services, you accept and agree to be bound by the terms and provisions of this agreement. If you do not agree to these terms, please do not use our services.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('2. Booking and Reservations'),
                        _buildSectionContent(
                          'All bookings are subject to availability and confirmation. Prices are subject to change without notice until payment is received. It is your responsibility to ensure all information provided is accurate and complete.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('3. Payment Terms'),
                        _buildSectionContent(
                          'Payment must be made in full at the time of booking unless otherwise agreed. We accept various payment methods including credit cards, debit cards, and digital wallets. All transactions are processed securely.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('4. Cancellations and Refunds'),
                        _buildSectionContent(
                          'Cancellation policies vary depending on the service provider and type of booking. Refunds, if applicable, will be processed according to the specific cancellation policy of each booking. Some bookings may be non-refundable.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('5. Travel Documents'),
                        _buildSectionContent(
                          'It is your responsibility to ensure you have valid travel documents including passports, visas, and any required health certificates. We are not responsible for denied boarding or entry due to invalid or missing documents.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('6. Liability'),
                        _buildSectionContent(
                          'Travel acts as an intermediary between you and service providers. We are not liable for any loss, damage, or injury resulting from services provided by third parties including airlines, hotels, and car rental companies.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('7. Force Majeure'),
                        _buildSectionContent(
                          'We are not liable for any failure to perform our obligations due to circumstances beyond our reasonable control, including natural disasters, war, terrorism, strikes, or government actions.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('8. Changes to Terms'),
                        _buildSectionContent(
                          'We reserve the right to modify these terms at any time. Changes will be effective immediately upon posting. Your continued use of our services constitutes acceptance of any changes.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('9. Contact Information'),
                        _buildSectionContent(
                          'For questions or concerns about these terms, please contact our customer service team through the app or visit our website.',
                        ),
                        const SizedBox(height: 20),
                        
                        Text(
                          'Last Updated: ${DateTime.now().year}',
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.grey,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                // Footer Button
                Padding(
                  padding: const EdgeInsets.all(20),
                  child: SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () => Navigator.pop(context),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text(
                        'I Understand',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Colors.white,
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  /// Show Privacy Policy Dialog
  static void showPrivacyPolicy(BuildContext context) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Container(
            constraints: const BoxConstraints(maxWidth: 600),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Header
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: const BoxDecoration(
                    color: Color(0xFF1e5a8e),
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(16),
                      topRight: Radius.circular(16),
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'privacy_policy'.tr,
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close, color: Colors.white),
                        onPressed: () => Navigator.pop(context),
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                    ],
                  ),
                ),
                // Content
                Flexible(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildSectionTitle('1. Information We Collect'),
                        _buildSectionContent(
                          'We collect information you provide directly to us, including name, email address, phone number, payment information, passport details, and travel preferences. We also collect information automatically such as device information, IP address, and usage data.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('2. How We Use Your Information'),
                        _buildSectionContent(
                          'We use your information to process bookings, communicate with you about reservations, provide customer support, send promotional offers (with your consent), improve our services, and comply with legal obligations.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('3. Information Sharing'),
                        _buildSectionContent(
                          'We share your information with service providers (airlines, hotels, car rental companies) necessary to complete your bookings. We may also share information with payment processors, business partners, and as required by law.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('4. Data Security'),
                        _buildSectionContent(
                          'We implement appropriate security measures to protect your personal information. However, no method of transmission over the internet is 100% secure. We use encryption for sensitive data and regularly update our security practices.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('5. Your Rights'),
                        _buildSectionContent(
                          'You have the right to access, correct, or delete your personal information. You can opt-out of marketing communications at any time. You may also request a copy of your data or object to certain processing activities.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('6. Cookies and Tracking'),
                        _buildSectionContent(
                          'We use cookies and similar technologies to improve your experience, analyze usage patterns, and personalize content. You can control cookie settings through your browser preferences.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('7. Data Retention'),
                        _buildSectionContent(
                          'We retain your information for as long as necessary to provide services, comply with legal obligations, resolve disputes, and enforce our agreements. Booking records are typically kept for 7 years for tax and legal purposes.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('8. Children\'s Privacy'),
                        _buildSectionContent(
                          'Our services are not directed to children under 13. We do not knowingly collect personal information from children. If you believe we have collected such information, please contact us immediately.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('9. International Transfers'),
                        _buildSectionContent(
                          'Your information may be transferred to and processed in countries other than your country of residence. We ensure appropriate safeguards are in place for such transfers.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('10. Changes to Privacy Policy'),
                        _buildSectionContent(
                          'We may update this privacy policy from time to time. We will notify you of significant changes via email or through the app. Your continued use after changes constitutes acceptance.',
                        ),
                        const SizedBox(height: 16),
                        
                        _buildSectionTitle('11. Contact Us'),
                        _buildSectionContent(
                          'For privacy-related questions or to exercise your rights, please contact our Data Protection Officer through the app or at privacy@travel.app.',
                        ),
                        const SizedBox(height: 20),
                        
                        Text(
                          'Last Updated: ${DateTime.now().year}',
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.grey,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                // Footer Button
                Padding(
                  padding: const EdgeInsets.all(20),
                  child: SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () => Navigator.pop(context),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: const Text(
                        'I Understand',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Colors.white,
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  static Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: const TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.bold,
        color: Color(0xFF1e5a8e),
      ),
    );
  }

  static Widget _buildSectionContent(String content) {
    return Padding(
      padding: const EdgeInsets.only(top: 8),
      child: Text(
        content,
        style: const TextStyle(
          fontSize: 14,
          color: Colors.black87,
          height: 1.5,
        ),
      ),
    );
  }
}
