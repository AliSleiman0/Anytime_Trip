import 'package:flutter/material.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/endpoints.dart';
import '../service/account_api.dart';

class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  // Notification preferences states
  Map<String, bool> notificationPreferences = {
    'Email': true,
    'SMS': true,
    'Support': true,
  };
  final _storage = StorageService();
  final _api = AccountApi();

  @override
  void initState() {
    super.initState();
    _loadNotificationPreferences();
  }

  Future<void> _loadNotificationPreferences() async {
    // Initialize with only 3 options
    final allowedKeys = {'Email', 'SMS', 'Support'};
    
    // Helper function to migrate old 'Chatbot' key to 'Support'
    Map<String, bool> migratePreferences(Map<String, bool> prefs) {
      final migrated = Map<String, bool>.from(prefs);
      if (migrated.containsKey('Chatbot') && !migrated.containsKey('Support')) {
        migrated['Support'] = migrated['Chatbot']!;
        migrated.remove('Chatbot');
      }
      return migrated;
    }
    
    try {
      // Try to load from API first
      final apiPrefs = await _api.getNotificationPreferences();
      
      setState(() {
        // Start with defaults
        notificationPreferences = {
          'Email': true,
          'SMS': true,
          'Support': true,
        };
        
        // Merge saved values for only the allowed keys
        if (apiPrefs != null && apiPrefs.isNotEmpty) {
          final migratedApiPrefs = migratePreferences(apiPrefs);
          for (var key in allowedKeys) {
            if (migratedApiPrefs.containsKey(key)) {
              notificationPreferences[key] = migratedApiPrefs[key]!;
            }
          }
        } else {
          // If no API data, check local storage
          final saved = _storage.getNotificationPreferences();
          if (saved != null && saved.isNotEmpty) {
            final migratedSaved = migratePreferences(saved);
            for (var key in allowedKeys) {
              if (migratedSaved.containsKey(key)) {
                notificationPreferences[key] = migratedSaved[key]!;
              }
            }
          }
        }
        
        // Save the cleaned preferences
        _saveNotificationPreferences();
      });
    } catch (e) {
      print('Error loading preferences from API: $e');
      // Fallback to local storage
      final saved = _storage.getNotificationPreferences();
      setState(() {
        // Start with defaults
        notificationPreferences = {
          'Email': true,
          'SMS': true,
          'Support': true,
        };
        
        // Merge saved values for only the allowed keys
        if (saved != null && saved.isNotEmpty) {
          final migratedSaved = migratePreferences(saved);
          for (var key in allowedKeys) {
            if (migratedSaved.containsKey(key)) {
              notificationPreferences[key] = migratedSaved[key]!;
            }
          }
        }
      });
    }
  }

  Future<void> _saveNotificationPreferences() async {
    try {
      // Save to local storage
      await _storage.saveNotificationPreferences(notificationPreferences);
      // Save to backend API
      await _api.saveNotificationPreferences(notificationPreferences);
    } catch (e) {
      print('Error saving preferences: $e');
      // Still save to local storage even if API fails
      await _storage.saveNotificationPreferences(notificationPreferences);
    }
  }

  String _getUserName() {
    final user = StorageService().getUser();
    return user?['name'] ?? 'User';
  }

  String? _getProfileImage() {
    final user = StorageService().getUser();
    return user?['profile_image'];
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
            Text(
              'Hello, ${_getUserName()}',
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
                child: _getProfileImage() != null && _getProfileImage()!.isNotEmpty
                    ? Image.network(
                        '${Endpoints.baseUrl.replaceAll('/api', '')}${_getProfileImage()}',
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return Image.asset(
                            'assets/images/defaultp.png',
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) {
                              return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                            },
                          );
                        },
                      )
                    : Image.asset(
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
                  // Save to storage when changed
                  _saveNotificationPreferences();
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
