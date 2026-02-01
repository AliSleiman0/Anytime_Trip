import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/account_ui_model.dart';

class AccountApi {
  final ApiClient _apiClient = ApiClient();

  Future<ProfileModel> getProfile() async {
    try {
      final response = await _apiClient.get(Endpoints.profile);
      return ProfileModel.fromJson(response.data);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> saveNotificationPreferences(Map<String, bool> preferences) async {
    try {
      await _apiClient.post(
        Endpoints.saveNotificationPreferences,
        data: preferences,
      );
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, bool>?> getNotificationPreferences() async {
    try {
      final response = await _apiClient.get(Endpoints.getNotificationPreferences);
      if (response.data['preferences'] != null) {
        final prefs = response.data['preferences'] as Map<String, dynamic>;
        return prefs.cast<String, bool>();
      }
      return null;
    } catch (e) {
      rethrow;
    }
  }

  // Payment Methods APIs
  Future<List<PaymentMethod>> getPaymentMethods() async {
    try {
      final response = await _apiClient.get(Endpoints.paymentMethods);
      final List<dynamic> data = response.data['data'] ?? [];
      return data.map((pm) => PaymentMethod.fromJson(pm as Map<String, dynamic>)).toList();
    } catch (e) {
      print('[PAYMENT_METHODS_ERROR] Failed to get payment methods: $e');
      rethrow;
    }
  }

  Future<PaymentMethod> addPaymentMethod(PaymentMethodRequest request) async {
    try {
      final response = await _apiClient.post(
        Endpoints.paymentMethods,
        data: request.toJson(),
      );
      return PaymentMethod.fromJson(response.data['data'] as Map<String, dynamic>);
    } catch (e) {
      rethrow;
    }
  }

  Future<PaymentMethod> updatePaymentMethod(String id, UpdatePaymentMethodRequest request) async {
    try {
      final response = await _apiClient.put(
        Endpoints.paymentMethodDetails(id),
        data: request.toJson(),
      );
      return PaymentMethod.fromJson(response.data['data'] as Map<String, dynamic>);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> deletePaymentMethod(String id) async {
    try {
      await _apiClient.delete(Endpoints.paymentMethodDetails(id));
    } catch (e) {
      rethrow;
    }
  }

  Future<void> setDefaultPaymentMethod(String id) async {
    try {
      await _apiClient.post(Endpoints.setDefaultPaymentMethod(id));
    } catch (e) {
      rethrow;
    }
  }

  // Security Preferences APIs
  Future<Map<String, dynamic>?> getSecurityPreferences() async {
    try {
      final response = await _apiClient.get(Endpoints.getSecurityPreferences);
      if (response.data['preferences'] != null) {
        return response.data['preferences'] as Map<String, dynamic>;
      }
      return null;
    } catch (e) {
      rethrow;
    }
  }

  Future<void> saveSecurityPreferences(Map<String, dynamic> preferences) async {
    try {
      await _apiClient.post(
        Endpoints.saveSecurityPreferences,
        data: preferences,
      );
    } catch (e) {
      rethrow;
    }
  }
}
