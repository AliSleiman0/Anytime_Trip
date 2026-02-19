import 'dart:io';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:country_picker/country_picker.dart';
import '../controller/auth_controller.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/widgets/policy_dialogs.dart';
import '../../../core/utils/helpers.dart';
import '../../../services/auth_service.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  String selectedLanguage = 'en';
  final nameController = TextEditingController();
  final phoneController = TextEditingController();
  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final confirmPasswordController = TextEditingController();
  bool agreeToTerms = false;
  String? selectedSex;
  String? selectedCountry;
  final _authService = AuthService();
  bool _isGoogleLoading = false;
  bool _isAppleLoading = false;

  @override
  void dispose() {
    nameController.dispose();
    phoneController.dispose();
    emailController.dispose();
    passwordController.dispose();
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
                SizedBox(height: size.height * 0.06),
                // Full Name Input
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
                      controller: nameController,
                      decoration: InputDecoration(
                        hintText: 'enter_full_name'.tr,
                        border: InputBorder.none,
                        prefixIcon: const Icon(
                          Icons.person,
                          color: Color(0xFF1e5a8e),
                        ),
                        contentPadding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Sex/Gender Dropdown
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
                    child: DropdownButtonFormField<String>(
                      value: selectedSex,
                      decoration: InputDecoration(
                        hintText: 'Select Sex',
                        border: InputBorder.none,
                        prefixIcon: const Icon(
                          Icons.person_outline,
                          color: Color(0xFF1e5a8e),
                        ),
                        contentPadding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
                      ),
                      items: ['Male', 'Female', 'Other'].map((String value) {
                        return DropdownMenuItem<String>(
                          value: value,
                          child: Text(value),
                        );
                      }).toList(),
                      onChanged: (String? newValue) {
                        setState(() {
                          selectedSex = newValue;
                        });
                      },
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                // Country of Residence Picker
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: GestureDetector(
                    onTap: () {
                      showCountryPicker(
                        context: context,
                        showPhoneCode: false,
                        onSelect: (Country country) {
                          setState(() {
                            selectedCountry = country.name;
                          });
                        },
                        countryListTheme: CountryListThemeData(
                          borderRadius: BorderRadius.circular(12),
                          inputDecoration: InputDecoration(
                            labelText: 'Search',
                            hintText: 'Start typing to search',
                            prefixIcon: const Icon(Icons.search),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                        ),
                      );
                    },
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
                      padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
                      child: Row(
                        children: [
                          const Icon(
                            Icons.flag,
                            color: Color(0xFF1e5a8e),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Text(
                              selectedCountry ?? 'Select Country of Residence',
                              style: TextStyle(
                                color: selectedCountry == null
                                    ? Colors.grey[600]
                                    : Colors.black87,
                                fontSize: 16,
                              ),
                            ),
                          ),
                          const Icon(
                            Icons.arrow_drop_down,
                            color: Color(0xFF1e5a8e),
                          ),
                        ],
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
                // Create Password Input
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
                        hintText: 'create_new_password'.tr,
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
                const SizedBox(height: 16),
                // Confirm Password Input
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
                      controller: confirmPasswordController,
                      obscureText: true,
                      decoration: InputDecoration(
                        hintText: 'confirm_new_password'.tr,
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
                const SizedBox(height: 16),
                // Terms & Conditions Checkbox
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    children: [
                      Checkbox(
                        value: agreeToTerms,
                        onChanged: (value) {
                          setState(() {
                            agreeToTerms = value ?? false;
                          });
                        },
                      ),
                      Expanded(
                        child: GestureDetector(
                          onTap: () {},
                          child: RichText(
                            text: TextSpan(
                              style: const TextStyle(
                                color: Colors.black87,
                                fontSize: 13,
                              ),
                              children: [
                                TextSpan(text: '${'i_agree_to_the'.tr} '),
                                WidgetSpan(
                                  child: GestureDetector(
                                    onTap: () {
                                      PolicyDialogs.showTermsAndConditions(context);
                                    },
                                    child: Text(
                                      'terms_conditions'.tr,
                                      style: const TextStyle(
                                        color: Color(0xFFD32F2F),
                                        decoration: TextDecoration.underline,
                                        fontSize: 13,
                                      ),
                                    ),
                                  ),
                                ),
                                const TextSpan(text: ' & '),
                                WidgetSpan(
                                  child: GestureDetector(
                                    onTap: () {
                                      PolicyDialogs.showPrivacyPolicy(context);
                                    },
                                    child: Text(
                                      'privacy_policy'.tr,
                                      style: const TextStyle(
                                        color: Color(0xFFD32F2F),
                                        decoration: TextDecoration.underline,
                                        fontSize: 13,
                                      ),
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                SizedBox(height: size.height * 0.06),
                // Sign Up Button
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
                        // Validate inputs
                        if (nameController.text.isEmpty) {
                          Helpers.showSnackbar('Required Field', 'Please enter your name', isError: true);
                          return;
                        }
                        if (emailController.text.isEmpty) {
                          Helpers.showSnackbar('Required Field', 'Please enter your email', isError: true);
                          return;
                        }
                        if (phoneController.text.isEmpty) {
                          Helpers.showSnackbar('Required Field', 'Please enter your phone number', isError: true);
                          return;
                        }
                        if (selectedSex == null) {
                          Helpers.showSnackbar('Required Field', 'Please select your sex', isError: true);
                          return;
                        }
                        if (selectedCountry == null) {
                          Helpers.showSnackbar('Required Field', 'Please select your country', isError: true);
                          return;
                        }
                        if (passwordController.text.isEmpty) {
                          Helpers.showSnackbar('Required Field', 'Please enter your password', isError: true);
                          return;
                        }
                        if (passwordController.text.length < 6) {
                          Helpers.showSnackbar('Invalid Password', 'Password must be at least 6 characters', isError: true);
                          return;
                        }
                        if (confirmPasswordController.text.isEmpty) {
                          Helpers.showSnackbar('Required Field', 'Please confirm your password', isError: true);
                          return;
                        }
                        if (passwordController.text != confirmPasswordController.text) {
                          Helpers.showSnackbar('Password Mismatch', 'Passwords do not match', isError: true);
                          return;
                        }
                        if (!agreeToTerms) {
                          Helpers.showSnackbar('Agreement Required', 'Please agree to the Terms and Conditions', isError: true);
                          return;
                        }

                        // Call the register function
                        final authController = Get.find<AuthController>();
                        authController.register(
                          name: nameController.text,
                          email: emailController.text,
                          password: passwordController.text,
                          confirmPassword: confirmPasswordController.text,
                          phoneNumber: phoneController.text,
                          sex: selectedSex!,
                          country: selectedCountry!,
                        );
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1e5a8e),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(25),
                        ),
                      ),
                      child: Text(
                        'signup'.tr,
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
                // Continue with Apple Button - Only show on iOS/macOS
                if (Platform.isIOS || Platform.isMacOS)
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 24),
                    child: Container(
                      width: double.infinity,
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
                        onPressed: _isAppleLoading ? null : () async {
                          setState(() => _isAppleLoading = true);
                          try {
                            final result = await _authService.signInWithApple();
                            
                            if (!mounted) return;
                            
                            if (result == null) {
                              Helpers.showSnackbar(
                                'Cancelled',
                                'Apple sign-in was cancelled',
                                isError: false,
                              );
                            } else {
                              final user = result['user'];
                              final isNewUser = result['isNewUser'] as bool;
                              
                              Helpers.showSnackbar(
                                'Welcome!',
                                'Signed in as ${user.email ?? "Apple User"}',
                                isError: false,
                              );
                              
                              // Navigate based on user status
                              if (isNewUser) {
                                // New user - go to complete profile
                                Get.offAllNamed(AppRoutes.COMPLETE_PROFILE);
                              } else {
                                // Existing user - go to home
                                Get.offAllNamed(AppRoutes.HOME);
                              }
                            }
                          } catch (e) {
                            if (!mounted) return;
                            Helpers.showSnackbar(
                              'Error',
                              'Apple sign-in failed: $e',
                              isError: true,
                            );
                          } finally {
                            if (mounted) setState(() => _isAppleLoading = false);
                          }
                        },
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
                        child: _isAppleLoading
                            ? const SizedBox(
                                height: 24,
                                width: 24,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                  valueColor: AlwaysStoppedAnimation<Color>(Colors.black),
                                ),
                              )
                            : Row(
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
                if (Platform.isIOS || Platform.isMacOS)
                  const SizedBox(height: 12),
                // Continue with Google Button
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Container(
                    width: double.infinity,
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
                      onPressed: _isGoogleLoading ? null : () async {
                        setState(() => _isGoogleLoading = true);
                        try {
                          final result = await _authService.signInWithGoogle();
                          
                          if (!mounted) return;
                          
                          if (result == null) {
                            Helpers.showSnackbar(
                              'Cancelled',
                              'Google sign-in was cancelled',
                              isError: false,
                            );
                          } else {
                            final user = result['user'];
                            final isNewUser = result['isNewUser'] as bool;
                            
                            Helpers.showSnackbar(
                              'Welcome!',
                              'Signed in as ${user.email}',
                              isError: false,
                            );
                            
                            // Navigate based on user status
                            if (isNewUser) {
                              // New user - go to complete profile
                              Get.offAllNamed(AppRoutes.COMPLETE_PROFILE);
                            } else {
                              // Existing user - go to home
                              Get.offAllNamed(AppRoutes.HOME);
                            }
                          }
                        } catch (e) {
                          if (!mounted) return;
                          Helpers.showSnackbar(
                            'Error',
                            'Google sign-in failed: $e',
                            isError: true,
                          );
                        } finally {
                          if (mounted) setState(() => _isGoogleLoading = false);
                        }
                      },
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
                      child: _isGoogleLoading
                          ? const SizedBox(
                              height: 24,
                              width: 24,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : Row(
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
                // Already have an account
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        '${'already_have_an_account'.tr} ',
                        style: const TextStyle(
                          color: Colors.black87,
                          fontSize: 14,
                        ),
                      ),
                      GestureDetector(
                        onTap: () => Get.back(),
                        child: Text(
                          'login'.tr,
                          style: const TextStyle(
                            color: Color(0xFFD32F2F),
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ],
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
