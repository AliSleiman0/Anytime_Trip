import 'package:get_storage/get_storage.dart';

class StorageService {
  final GetStorage _box = GetStorage();

  // Keys
  static const String _tokenKey = 'auth_token';
  static const String _userKey = 'user_data';
  static const String _themeKey = 'theme_mode';
  static const String _notificationPreferencesKey = 'notification_preferences';

  // Token operations
  Future<void> saveToken(String token) async {
    await _box.write(_tokenKey, token);
  }

  String? getToken() {
    return _box.read(_tokenKey);
  }

  Future<void> removeToken() async {
    await _box.remove(_tokenKey);
  }

  // User operations
  Future<void> saveUser(Map<String, dynamic> user) async {
    await _box.write(_userKey, user);
  }

  Map<String, dynamic>? getUser() {
    return _box.read(_userKey);
  }

  Future<void> removeUser() async {
    await _box.remove(_userKey);
  }

  // Theme operations
  Future<void> saveThemeMode(String mode) async {
    await _box.write(_themeKey, mode);
  }

  String? getThemeMode() {
    return _box.read(_themeKey);
  }

  // Notification preferences operations
  Future<void> saveNotificationPreferences(Map<String, bool> preferences) async {
    await _box.write(_notificationPreferencesKey, preferences);
  }

  Map<String, bool>? getNotificationPreferences() {
    final data = _box.read(_notificationPreferencesKey);
    if (data is Map) {
      return Map<String, bool>.from(data);
    }
    return null;
  }

  Future<void> removeNotificationPreferences() async {
    await _box.remove(_notificationPreferencesKey);
  }

  // Clear all data
  Future<void> clearAll() async {
    await _box.erase();
  }
}
