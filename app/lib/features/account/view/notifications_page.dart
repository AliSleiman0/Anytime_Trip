import 'package:flutter/material.dart';

class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  // Notification preferences states
  late Map<String, bool> notificationPreferences;

  @override
  void initState() {
    super.initState();
    notificationPreferences = {
      'All Notifications': false,
      'Push Notifications': true,
      'Email Alerts': false,
      'SMS Updates': true,
      'In-App Messages': false,
      'Social Media Alerts': true,
      'Webhook Notifications': false,
      'Browser Notifications': true,
      'RSS Feed Updates': false,
      'Chatbot Messages': true,
    };
  }

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
            const Text(
              'Hello, User123',
              style: TextStyle(
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
                  'assets/images/defaultp.png',
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
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 16),
              ..._buildNotificationItems(),
            ],
          ),
        ),
      ),
    );
  }

  List<Widget> _buildNotificationItems() {
    return notificationPreferences.entries.map((entry) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 12),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Expanded(
              child: Text(
                entry.key,
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: Colors.black87,
                ),
              ),
            ),
            Transform.scale(
              scale: 0.8,
              child: Switch(
                value: entry.value,
                onChanged: (bool value) {
                  setState(() {
                    notificationPreferences[entry.key] = value;
                  });
                },
                activeColor: const Color(0xFF336891),
                inactiveThumbColor: const Color(0xFFB0BEC5),
                inactiveTrackColor: const Color(0xFFE0E0E0),
              ),
            ),
          ],
        ),
      );
    }).toList();
  }
}
