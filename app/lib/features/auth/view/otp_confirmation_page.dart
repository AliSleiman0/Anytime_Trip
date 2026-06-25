import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/utils/helpers.dart';
import '../controller/auth_controller.dart';

class OtpConfirmationPage extends StatefulWidget {
  const OtpConfirmationPage({super.key});

  @override
  State<OtpConfirmationPage> createState() => _OtpConfirmationPageState();
}

class _OtpConfirmationPageState extends State<OtpConfirmationPage> {
  String selectedLanguage = 'en';
  final authController = Get.find<AuthController>();
  
  // Local variables from navigation arguments
  String userEmail = '';
  String userPhone = '';
  
  List<String> otpDigits = ['', '', '', ''];

  @override
  void initState() {
    super.initState();
    // Get email and phone from navigation arguments
    final args = Get.arguments as Map<String, dynamic>?;
    if (args != null) {
      userEmail = args['email'] ?? '';
      userPhone = args['phone'] ?? '';
    }
    
    print('DEBUG [OTP Page initState]: email=$userEmail, phone=$userPhone');
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
    if (userEmail.isEmpty) {
      Helpers.showSnackbar('Error', 'Email is missing', isError: true);
      return;
    }

    // Clear existing OTP inputs
    setState(() {
      otpDigits = ['', '', '', ''];
    });

    // Send new OTP via email
    authController.sendOTP(userEmail, '', 'email');
  }

  Future<void> _confirmOtp() async {
    String otp = otpDigits.join();

    if (otp.length < 4) {
      Helpers.showSnackbar('Invalid OTP', 'Please enter the 4-digit code', isError: true);
      return;
    }

    if (userEmail.isEmpty) {
      Helpers.showSnackbar('Error', 'Email is missing', isError: true);
      return;
    }

    print('DEBUG [_confirmOtp]: Verifying OTP - email=$userEmail, code=$otp');

    // Verify OTP via email
    final verified = await authController.verifyOTP(userEmail, '', otp, 'email');
    
    if (verified) {
      // Activate account
      final activated = await authController.activateAccount(userEmail);
      
      if (activated) {
        // Redirect to login page
        Get.offAllNamed(AppRoutes.LOGIN);
        Helpers.showSnackbar('Account Activated', 'You can now login with your credentials', isError: false);
      }
    }
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
                const Text(
                  'Travel',
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                    letterSpacing: -0.5,
                  ),
                ),
                SizedBox(height: size.height * 0.08),
                // Title
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Column(
                    children: [
                      const Text(
                        'Verify OTP',
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Enter the 4-digit code sent to\n$userEmail',
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                          fontSize: 14,
                          color: Colors.white70,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 32),
                // Display OTP input boxes for verification code
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      // Input OTP boxes (4 digits)
                      ...List.generate(4, (index) {
                        return Container(
                          width: 60,
                          height: 60,
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
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                              color: Colors.black,
                            ),
                            onChanged: (value) {
                              setState(() {
                                otpDigits[index] = value;
                              });
                              if (value.isNotEmpty && index < 3) {
                                FocusScope.of(context).nextFocus();
                              }
                            },
                          ),
                        );
                      }),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
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
                              text: "Didn't Receive OTP? ",
                              style: const TextStyle(
                                color: Colors.white70,
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
                SizedBox(height: size.height * 0.12),
                // Enter Button
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
                      onPressed: _confirmOtp,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Text(
                        'Verify & Continue'.tr,
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
                // Back Button
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
                        'back'.tr,
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
