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
      return UserModel.fromJson(response.data);
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
}
