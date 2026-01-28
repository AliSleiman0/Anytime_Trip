import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../model/auth_ui_model.dart';

class AuthApi {
  final ApiClient _apiClient = ApiClient();

  Future<UserModel> login(LoginModel credentials) async {
    try {
      final response = await _apiClient.post(
        Endpoints.login,
        data: credentials.toJson(),
      );
      
      // Backend returns: { message, token, user: {...} }
      final userData = response.data['user'] as Map<String, dynamic>;
      userData['token'] = response.data['token']; // Add token to user data
      
      return UserModel.fromJson(userData);
    } catch (e) {
      rethrow;
    }
  }

  Future<UserModel> register(RegisterModel userData) async {
    try {
      final response = await _apiClient.post(
        Endpoints.register,
        data: userData.toJson(),
      );
      return UserModel.fromJson(response.data);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> logout() async {
    try {
      await _apiClient.post(Endpoints.logout);
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> sendOTP(String email, String phone, String type) async {
    try {
      final data = <String, dynamic>{'type': type};
      if (type == 'email' && email.isNotEmpty) {
        data['email'] = email;
      } else if (type == 'phone' && phone.isNotEmpty) {
        data['phone'] = phone;
      }
      
      final response = await _apiClient.post(
        Endpoints.sendOTP,
        data: data,
      );
      return response.data;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> verifyOTP(String email, String phone, String code, String type) async {
    try {
      final data = <String, dynamic>{
        'code': code,
        'type': type,
      };
      
      if (type == 'email' && email.isNotEmpty) {
        data['email'] = email;
      } else if (type == 'phone' && phone.isNotEmpty) {
        data['phone'] = phone;
      }
      
      final response = await _apiClient.post(
        Endpoints.verifyOTP,
        data: data,
      );
      return response.data;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> activateAccount(String email) async {
    try {
      final response = await _apiClient.post(
        Endpoints.activateAccount,
        data: {'email': email},
      );
      return response.data;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> resetPassword(String phone, String otp, String newPassword) async {
    try {
      final response = await _apiClient.post(
        Endpoints.resetPassword,
        data: {
          'phone': phone,
          'otp': otp,
          'new_password': newPassword,
        },
      );
      return response.data;
    } catch (e) {
      rethrow;
    }
  }
}
