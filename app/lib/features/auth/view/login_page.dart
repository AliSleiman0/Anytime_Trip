import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import '../../../app/routes/app_routes.dart';
import '../controller/auth_controller.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  String selectedLanguage = 'en';
  final emailController = TextEditingController();
  final phoneController = TextEditingController();
  final passwordController = TextEditingController();

  @override
  void dispose() {
    emailController.dispose();
    phoneController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  void _changeLanguage(String languageCode) {
    setState(() {
      selectedLanguage = languageCode;
    });
    if (languageCode == 'en') {
      Get.updateLocale(const Locale('en', 'US'));
    } else if (languageCode == 'ar') {
      Get.updateLocale(const Locale('ar', 'SA'));
    }
  }

  void _showLanguageMenu() {
    showMenu(
      context: context,
      position: const RelativeRect.fromLTRB(double.maxFinite, 60, 20, 0),
      items: [
        PopupMenuItem(
          value: 'en',
          child: Row(
            children: [
              Radio(
                value: 'en',
                groupValue: selectedLanguage,
                onChanged: (value) {
                  _changeLanguage(value as String);
                  Navigator.pop(context);
                },
              ),
              const Text('English'),
            ],
          ),
        ),
        PopupMenuItem(
          value: 'ar',
          child: Row(
            children: [
              Radio(
                value: 'ar',
                groupValue: selectedLanguage,
                onChanged: (value) {
                  _changeLanguage(value as String);
                  Navigator.pop(context);
                },
              ),
              const Text('العربية'),
            ],
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: GestureDetector(
              onTap: _showLanguageMenu,
              child: Container(
                width: 45,
                height: 45,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: Colors.white.withOpacity(0.7),
                ),
                child: const Icon(
                  Icons.language,
                  color: Color(0xFF1e5a8e),
                  size: 24,
                ),
              ),
            ),
          ),
        ],
      ),
      body: Stack(
        children: [
          // Background image
          Container(
            width: double.infinity,
            height: double.infinity,
            decoration: const BoxDecoration(
              image: DecorationImage(
                image: AssetImage('assets/images/background1.jpg'),
                fit: BoxFit.cover,
              ),
            ),
          ),
          SingleChildScrollView(
            child: Column(
              children: [
                SizedBox(height: size.height * 0.12),
                // Logo
                Image.asset(
                  'assets/images/main-logo.png',
                  width: 90,
                  height: 90,
                ),
                SizedBox(height: size.height * 0.08),
                // Email Input
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    decoration: BoxDecoration(
                      color: Color(0xFFF5F5F5),
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 8,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: TextField(
                      controller: emailController,
                      keyboardType: TextInputType.emailAddress,
                      decoration: InputDecoration(
                        hintText: 'enter_email_address'.tr,
                        border: InputBorder.none,
                        prefixIcon: const Icon(
                          Icons.email,
                          color: Color(0xFF1e5a8e),
                        ),
                        contentPadding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Phone Number Input
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: IntlPhoneField(
                    controller: phoneController,
                    decoration: InputDecoration(
                      hintText: '00 123 456',
                      border: InputBorder.none,
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(25),
                        borderSide: BorderSide.none,
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(25),
                        borderSide: BorderSide.none,
                      ),
                      filled: true,
                      fillColor: Color(0xFFF5F5F5),
                      contentPadding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
                    ),
                    initialCountryCode: 'LB',
                    onChanged: (phone) {
                      print(phone.completeNumber);
                    },
                  ),
                ),
                const SizedBox(height: 16),
                // Password Input
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    decoration: BoxDecoration(
                      color: Color(0xFFF5F5F5),
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 8,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: TextField(
                      controller: passwordController,
                      obscureText: true,
                      decoration: InputDecoration(
                        hintText: 'enter_password'.tr,
                        border: InputBorder.none,
                        prefixIcon: const Icon(
                          Icons.lock,
                          color: Color(0xFF1e5a8e),
                        ),
                        contentPadding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // Forgot Password Link
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Align(
                    alignment: Alignment.centerRight,
                    child: GestureDetector(
                      onTap: () {
                        Get.toNamed(AppRoutes.FORGOT_PASSWORD);
                      },
                      child: Text(
                        'forgot_password'.tr,
                        style: const TextStyle(
                          color: Color(0xFFD32F2F),
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ),
                ),
                SizedBox(height: size.height * 0.08),
                // Login Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    width: double.infinity,
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.15),
                          blurRadius: 8,
                          offset: const Offset(0, 4),
                        ),
                      ],
                    ),
                    child: ElevatedButton(
                      onPressed: () {
                        // Handle login and navigate to home
                        Get.offNamed(AppRoutes.HOME);
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Text(
                        'login'.tr,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Continue with Apple Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.1),
                          blurRadius: 6,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: OutlinedButton(
                      onPressed: () {},
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        side: const BorderSide(
                          color: Colors.black,
                          width: 1.5,
                        ),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset(
                            'assets/images/apple.png',
                            width: 24,
                            height: 24,
                          ),
                          const SizedBox(width: 12),
                          Text(
                            'continue_with_apple'.tr,
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // Continue with Google Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.1),
                          blurRadius: 6,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: OutlinedButton(
                      onPressed: () {},
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        side: const BorderSide(
                          color: Colors.black,
                          width: 1.5,
                        ),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset(
                            'assets/images/gmail.png',
                            width: 24,
                            height: 24,
                          ),
                          const SizedBox(width: 12),
                          Text(
                            'continue_with_google'.tr,
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 20),
                // Don't have an account
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      '${'dont_have_an_account'.tr} ',
                      style: const TextStyle(
                        color: Colors.black87,
                        fontSize: 14,
                      ),
                    ),
                    GestureDetector(
                      onTap: () => Get.toNamed(AppRoutes.REGISTER),
                      child: Text(
                        'signup'.tr,
                        style: const TextStyle(
                          color: Color(0xFFD32F2F),
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 32),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
