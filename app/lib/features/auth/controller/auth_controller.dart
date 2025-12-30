import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/utils/helpers.dart';
import '../service/auth_api.dart';
import '../model/auth_ui_model.dart';

class AuthController extends GetxController {
  final AuthApi _authApi = Get.find();

  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final nameController = TextEditingController();

  final isLoading = false.obs;
  final isPasswordVisible = false.obs;

  @override
  void onClose() {
    emailController.dispose();
    passwordController.dispose();
    nameController.dispose();
    super.onClose();
  }

  void togglePasswordVisibility() {
    isPasswordVisible.value = !isPasswordVisible.value;
  }

  Future<void> login() async {
    try {
      isLoading.value = true;
      final credentials = LoginModel(
        email: emailController.text,
        password: passwordController.text,
      );

      final response = await _authApi.login(credentials);
      
      Helpers.showSnackbar('Success', 'Login successful');
      Get.offAllNamed(AppRoutes.PRODUCT);
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> register() async {
    try {
      isLoading.value = true;
      final userData = RegisterModel(
        name: nameController.text,
        email: emailController.text,
        password: passwordController.text,
      );

      final response = await _authApi.register(userData);
      
      Helpers.showSnackbar('Success', 'Registration successful');
      Get.offAllNamed(AppRoutes.LOGIN);
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> logout() async {
    try {
      await _authApi.logout();
      Get.offAllNamed(AppRoutes.LOGIN);
    } catch (e) {
      Helpers.showSnackbar('Error', e.toString(), isError: true);
    }
  }
}
