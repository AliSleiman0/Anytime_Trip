import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';

class SetNewPasswordPage extends StatefulWidget {
  const SetNewPasswordPage({super.key});

  @override
  State<SetNewPasswordPage> createState() => _SetNewPasswordPageState();
}

class _SetNewPasswordPageState extends State<SetNewPasswordPage> {
  String selectedLanguage = 'en';
  final newPasswordController = TextEditingController();
  final confirmPasswordController = TextEditingController();
  
  List<String> otpDigits = ['', '', '', '', '', ''];
  
  bool newPasswordVisible = false;
  bool confirmPasswordVisible = false;

  @override
  void dispose() {
    newPasswordController.dispose();
    confirmPasswordController.dispose();
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

  void _resendOtp() {
    // TODO: Call API to resend OTP
    Get.snackbar('Success', 'OTP sent again');
  }

  void _changePassword() {
    if (newPasswordController.text.isEmpty) {
      Get.snackbar('Error', 'Please enter new password');
      return;
    }
    
    if (confirmPasswordController.text.isEmpty) {
      Get.snackbar('Error', 'Please confirm password');
      return;
    }
    
    if (newPasswordController.text != confirmPasswordController.text) {
      Get.snackbar('Error', 'Passwords do not match');
      return;
    }

    String otp = otpDigits.join();
    if (otp.length < 6) {
      Get.snackbar('Error', 'Please enter OTP');
      return;
    }

    // TODO: Call API to change password
    Get.snackbar('Success', 'Password changed successfully');
    // Navigate to login
    Get.offNamed(AppRoutes.LOGIN);
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
                SizedBox(height: size.height * 0.08),
                // Logo
                Image.asset(
                  'assets/images/main-logo.png',
                  width: 70,
                  height: 70,
                ),
                SizedBox(height: size.height * 0.08),
                // Create New Password Input
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
                      controller: newPasswordController,
                      obscureText: !newPasswordVisible,
                      decoration: InputDecoration(
                        hintText: 'create_new_password'.tr,
                        border: InputBorder.none,
                        prefixIcon: const Icon(
                          Icons.lock,
                          color: Color(0xFF1e5a8e),
                        ),
                        contentPadding: const EdgeInsets.symmetric(vertical: 14, horizontal: 12),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Confirm New Password Input with Send OTP Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    children: [
                      Expanded(
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
                            controller: confirmPasswordController,
                            obscureText: !confirmPasswordVisible,
                            decoration: InputDecoration(
                              hintText: 'confirm_new_password'.tr,
                              border: InputBorder.none,
                              prefixIcon: const Icon(
                                Icons.lock,
                                color: Color(0xFF1e5a8e),
                              ),
                              contentPadding: const EdgeInsets.symmetric(vertical: 14, horizontal: 12),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Container(
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
                          onPressed: _resendOtp,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFF1e5a8e),
                            padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(25),
                            ),
                          ),
                          child: Text(
                            'send_otp'.tr,
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.bold,
                              color: Colors.white,
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),
                // OTP Input boxes
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      // Input OTP boxes
                      ...List.generate(6, (index) {
                        return Container(
                          width: 50,
                          height: 50,
                          decoration: BoxDecoration(
                            color: const Color(0xFFF5F5F5),
                            borderRadius: BorderRadius.circular(10),
                            border: Border.all(
                              color: Colors.grey.shade300,
                            ),
                          ),
                          alignment: Alignment.center,
                          child: TextField(
                            textAlign: TextAlign.center,
                            keyboardType: TextInputType.number,
                            maxLength: 1,
                            decoration: const InputDecoration(
                              border: InputBorder.none,
                              counterText: '',
                              contentPadding: EdgeInsets.zero,
                            ),
                            style: const TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: Colors.black,
                            ),
                            onChanged: (value) {
                              setState(() {
                                otpDigits[index] = value;
                              });
                              if (value.isNotEmpty && index < 5) {
                                FocusScope.of(context).nextFocus();
                              }
                            },
                          ),
                        );
                      }),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
                // Resend OTP Link
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Align(
                    alignment: Alignment.centerRight,
                    child: GestureDetector(
                      onTap: _resendOtp,
                      child: RichText(
                        text: TextSpan(
                          children: [
                            TextSpan(
                              text: "Didn't Receive Password? ",
                              style: const TextStyle(
                                color: Colors.grey,
                                fontSize: 12,
                              ),
                            ),
                            TextSpan(
                              text: 'send again',
                              style: const TextStyle(
                                color: Color(0xFFD32F2F),
                                fontSize: 12,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                SizedBox(height: size.height * 0.08),
                // Change Password Button
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
                      onPressed: _changePassword,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Text(
                        'change_password'.tr,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // Cancel Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    width: double.infinity,
                    child: OutlinedButton(
                      onPressed: () => Get.back(),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        side: const BorderSide(
                          color: Color(0xFF1e5a8e),
                          width: 2,
                        ),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Text(
                        'cancel'.tr,
                        style: const TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFF1e5a8e),
                        ),
                      ),
                    ),
                  ),
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
