import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/utils/helpers.dart';
import '../../../core/storage/storage_service.dart';
import '../service/auth_api.dart';
import '../model/auth_ui_model.dart';

class AuthController extends GetxController {
  final AuthApi _authApi = Get.find();
  final StorageService _storage = StorageService();

  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final nameController = TextEditingController();
  final phoneController = TextEditingController();
  final confirmPasswordController = TextEditingController();

  final isLoading = false.obs;
  final isPasswordVisible = false.obs;
  
  // Store user data after signup for OTP verification
  final userEmail = ''.obs;
  final userPhone = ''.obs;

  @override
  void onClose() {
    emailController.dispose();
    passwordController.dispose();
    nameController.dispose();
    phoneController.dispose();
    confirmPasswordController.dispose();
    super.onClose();
  }

  void togglePasswordVisibility() {
    isPasswordVisible.value = !isPasswordVisible.value;
  }

  Future<void> login() async {
    try {
      isLoading.value = true;
      
      // Validate inputs
      if (emailController.text.isEmpty || passwordController.text.isEmpty) {
        Helpers.showSnackbar('Error', 'Email and password are required', isError: true);
        return;
      }
      
      final credentials = LoginModel(
        email: emailController.text,
        password: passwordController.text,
      );

      final user = await _authApi.login(credentials);
      
      // Save token and user data to storage
      if (user.token != null) {
        await _storage.saveToken(user.token!);
        await _storage.saveUser(user.toJson());
      }
      
      Helpers.showSnackbar('Success', 'Welcome back!', isError: false);
      Get.offAllNamed(AppRoutes.HOME);
    } catch (e) {
      print('DEBUG Login Error: $e');
      String errorMessage = e.toString().replaceAll('Exception: ', '');
      Helpers.showSnackbar('Login Failed', errorMessage, isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> register({
    required String name,
    required String email,
    required String password,
    required String confirmPassword,
    required String phoneNumber,
    required String sex,
    required String country,
  }) async {
    try {
      isLoading.value = true;
      final userData = RegisterModel(
        name: name,
        email: email,
        password: password,
        confirmPassword: confirmPassword,
        phoneNumber: phoneNumber,
        sex: sex,
        country: country,
      );

      final response = await _authApi.register(userData);
      
      // Store user data for OTP verification
      userEmail.value = email;
      userPhone.value = phoneNumber;
      
      print('DEBUG Register: Stored email=$email, phone=$phoneNumber');
      print('DEBUG Register: userEmail.value=${userEmail.value}, userPhone.value=${userPhone.value}');
      
      Helpers.showSnackbar('Account Created', 'Please verify your phone number', isError: false);
      // Navigate to OTP confirmation page for signup verification
      Get.offAllNamed(AppRoutes.OTP_CONFIRMATION, arguments: {
        'email': email,
        'phone': phoneNumber,
      });
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> logout() async {
    try {
      await _authApi.logout();
      
      // Clear stored data
      await _storage.clearAll();
      
      Get.offAllNamed(AppRoutes.LOGIN);
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    }
  }

  Future<void> sendOTP(String email, String phone, String type) async {
    try {
      isLoading.value = true;
      final response = await _authApi.sendOTP(email, phone, type);
      
      // OTP will be sent to phone via SMS
      print('DEBUG: OTP sent to phone');
      Helpers.showSnackbar('OTP Sent', 'Check your phone for verification code', isError: false);
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> verifyOTP(String email, String phone, String code, String type) async {
    try {
      isLoading.value = true;
      final response = await _authApi.verifyOTP(email, phone, code, type);
      Helpers.showSnackbar('Verified', 'Account verified successfully', isError: false);
      return true;
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> activateAccount(String email) async {
    try {
      isLoading.value = true;
      final response = await _authApi.activateAccount(email);
      return true;
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> resetPassword(String phone, String otp, String newPassword) async {
    try {
      isLoading.value = true;
      final response = await _authApi.resetPassword(phone, otp, newPassword);
      Helpers.showSnackbar('Success', 'Password reset successfully', isError: false);
      return true;
    } catch (e) {
      String errorMessage = e.toString().replaceAll('Exception: ', '');
      Helpers.showSnackbar('Reset Failed', errorMessage, isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> forgotPassword(String email) async {
    try {
      isLoading.value = true;
      final response = await _authApi.forgotPassword(email);
      Helpers.showSnackbar('Success', 'Verification code sent to your email', isError: false);
      return true;
    } catch (e) {
      String errorMessage = e.toString().replaceAll('Exception: ', '');
      Helpers.showSnackbar('Error', errorMessage, isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> verifyResetCode(String email, String code) async {
    try {
      isLoading.value = true;
      final response = await _authApi.verifyResetCode(email, code);
      Helpers.showSnackbar('Success', 'Account verified successfully', isError: false);
      return true;
    } catch (e) {
      String errorMessage = e.toString().replaceAll('Exception: ', '');
      Helpers.showSnackbar('Error', errorMessage, isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  Future<bool> resetPasswordWithEmail(String email, String code, String newPassword) async {
    try {
      isLoading.value = true;
      final response = await _authApi.resetPasswordWithEmail(email, code, newPassword);
      Helpers.showSnackbar('Success', 'Password reset successfully', isError: false);
      return true;
    } catch (e) {
      String errorMessage = e.toString().replaceAll('Exception: ', '');
      Helpers.showSnackbar('Error', errorMessage, isError: true);
      return false;
    } finally {
      isLoading.value = false;
    }
  }
}
