import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../controller/account_controller.dart';
import '../../../core/widgets/unified_ui_components.dart';
import '../../../core/storage/storage_service.dart';
import '../../../app/routes/app_routes.dart';
import 'profile_page.dart';
import 'notifications_page.dart';
import 'payment_methods_page.dart';
import 'security_page.dart';
import 'bookings_page.dart';

class AccountPage extends GetView<AccountController> {
  const AccountPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e), size: 24),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Hello, User123',
              style: const TextStyle(
                color: Colors.black87,
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(width: 12),
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: const Color(0xFFE0E0E0),
              ),
              child: ClipOval(
                child: Image.asset(
                  'assets/images/profile_avatar.png',
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                  },
                ),
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              child: Column(
                children: [
                  const SizedBox(height: 8),
                  _buildMenuItemWithIcon(
                    context,
                    icon: Icons.person,
                    label: 'Profile',
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => const ProfilePage()),
                    ),
                  ),
                  _buildMenuItemWithIcon(
                    context,
                    icon: Icons.notifications,
                    label: 'Notifications',
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => const NotificationsPage()),
                    ),
                  ),
                  _buildMenuItemWithIcon(
                    context,
                    icon: Icons.credit_card,
                    label: 'Payment Methods',
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => const PaymentMethodsPage()),
                    ),
                  ),
                  _buildMenuItemWithIcon(
                    context,
                    icon: Icons.security,
                    label: 'Security',
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => const SecurityPage()),
                    ),
                  ),
                  _buildMenuItemWithIcon(
                    context,
                    icon: Icons.bookmark,
                    label: 'My Bookings',
                    onTap: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => const BookingsPage()),
                    ),
                  ),
                ],
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16),
            child: SizedBox(
              width: double.infinity,
              child: OutlinedButton(
                onPressed: () async {
                  // Clear stored user data
                  final storage = StorageService();
                  await storage.clearAll();
                  
                  // Navigate to login page
                  Get.offAllNamed(AppRoutes.LOGIN);
                },
                style: OutlinedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  side: const BorderSide(color: Color(0xFFD32F2F), width: 1.5),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.logout, color: Color(0xFFD32F2F), size: 20),
                    const SizedBox(width: 8),
                    const Text(
                      'Sign out',
                      style: TextStyle(
                        color: Color(0xFFD32F2F),
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMenuItemWithIcon(
    BuildContext context, {
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Icon(icon, color: const Color(0xFFD32F2F), size: 26),
                  const SizedBox(width: 20),
                  Text(
                    label,
                    style: const TextStyle(
                      color: Colors.black87,
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
              const Icon(Icons.arrow_forward_ios, color: Color(0xFF1e5a8e), size: 20),
            ],
          ),
        ),
      ),
    );
  }
}
