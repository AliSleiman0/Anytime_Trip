import 'package:flutter/material.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/endpoints.dart';
import '../service/account_api.dart';

class SecurityPage extends StatefulWidget {
  const SecurityPage({super.key});

  @override
  State<SecurityPage> createState() => _SecurityPageState();
}

class _SecurityPageState extends State<SecurityPage> {
  final _api = AccountApi();
  bool _isLoading = true;
  bool _twoFactorAuth = false;

  @override
  void initState() {
    super.initState();
    _loadSecurityPreferences();
  }

  Future<void> _loadSecurityPreferences() async {
    try {
      setState(() => _isLoading = true);
      final prefs = await _api.getSecurityPreferences();
      if (prefs != null) {
        setState(() {
          _twoFactorAuth = prefs['two_factor_auth'] ?? false;
          _isLoading = false;
        });
      } else {
        setState(() => _isLoading = false);
      }
    } catch (e) {
      print('[SECURITY_PREFS_ERROR] Failed to load security preferences: $e');
      setState(() => _isLoading = false);
    }
  }

  Future<void> _saveSecurityPreferences() async {
    try {
      final preferences = {
        'two_factor_auth': _twoFactorAuth,
      };

      await _api.saveSecurityPreferences(preferences);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Two-Factor Authentication updated')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to save preferences: $e')),
        );
      }
    }
  }

  void _on2FAToggle(bool value) {
    setState(() => _twoFactorAuth = value);
    _saveSecurityPreferences();
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
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const SizedBox(height: 16),
                    const Text(
                      'Security Settings',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: Colors.black87,
                      ),
                    ),
                    const SizedBox(height: 16),
                    _buildSecurityToggle(
                      'Enable Two-Factor Authentication',
                      'Add an extra layer of security to your account by requiring a second verification method',
                      _twoFactorAuth,
                      _on2FAToggle,
                    ),
                    const SizedBox(height: 24),
                  ],
                ),
              ),
            ),
    );
  }

  Widget _buildSecurityToggle(
    String title,
    String description,
    bool value,
    Function(bool) onChanged,
  ) {
    return Container(
      decoration: BoxDecoration(
        color: const Color(0xFFD4DFE8),
        borderRadius: BorderRadius.circular(12),
      ),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  description,
                  style: const TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w400,
                    color: Colors.black54,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 16),
          Transform.scale(
            scale: 0.8,
            child: Switch(
              value: value,
              onChanged: onChanged,
              activeColor: const Color(0xFFD24124),
              inactiveThumbColor: const Color(0xFFB0BEC5),
              inactiveTrackColor: const Color(0xFFE0E0E0),
            ),
          ),
        ],
      ),
    );
  }
}
