import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:flutter_svg/flutter_svg.dart';
import '../../../core/widgets/custom_button.dart';
import '../../../app/routes/app_routes.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  String selectedLanguage = 'en'; // Default to English

  void _changeLanguage(String languageCode) {
    setState(() {
      selectedLanguage = languageCode;
    });
    // Change the app locale
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

    // Responsive sizes based on screen width
    final planeWidth = size.width * 0.72;
    final cloudLarge = size.width * 0.55;
    final cloudMedium = size.width * 0.42;
    final cloudSmall = size.width * 0.34;

    return Scaffold(
      body: Column(
        children: [
          // 🔹 TOP PART: background + logo + plane + clouds
          Expanded(
            child: Stack(
              children: [
                // Background left
                Positioned.fill(
                  child: Image.asset(
                    'assets/images/leftback.png',
                    fit: BoxFit.cover,
                  ),
                ),
                // Background right overlay
                Positioned.fill(
                  child: Image.asset(
                    'assets/images/leftback2.png',
                    fit: BoxFit.cover,
                  ),
                ),

                // Foreground content
                SafeArea(
                  child: Column(
                    children: [
                      // Top section with branding
                      Padding(
                        padding: const EdgeInsets.all(20),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            // Logo on left
                            Image.asset(
                              'assets/images/main-logo.png',
                              width: 90,
                              height: 90,
                            ),
                            // Language button on right
                            GestureDetector(
                              onTap: _showLanguageMenu,
                              child: Container(
                                padding: const EdgeInsets.all(8),
                                decoration: BoxDecoration(
                                  color: Colors.white.withOpacity(0.2),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                child: Image.asset(
                                  'assets/images/Language.png',
                                  width: 30,
                                  height: 30,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),

                      // Middle section - Airplane and Clouds
                     Expanded(
  child: Stack(
    clipBehavior: Clip.none,
    children: [
      // ---- LEFT CLUSTER ----

      // cloud22 (middle-left at edge, under cloud11)
      Positioned(
        left: 0,
        bottom: size.height * 0.25,
        child: Image.asset(
          'assets/images/cloud22.png',
          width: size.width * 0.32,
        ),
      ),

      // cloud11 (bottom-left at edge, overlays on cloud22)
      Positioned(
        left: 0,
        bottom: size.height * 0.18,
        child: Image.asset(
          'assets/images/cloud11.png',
          width: size.width * 0.48,
        ),
      ),

      // cloud33 (upper-left at edge)
      Positioned(
        left: 0,
        bottom: size.height * 0.42,
        child: Image.asset(
          'assets/images/cloud33.png',
          width: size.width * 0.30,
        ),
      ),

      // cloud44 (left-center)
      Positioned(
        left: 0,
        bottom: size.height * 0.40,
        child: Image.asset(
          'assets/images/cloud44.png',
          width: size.width * 0.26,
        ),
      ),

      // ---- RIGHT CLUSTER ----

      // cloud555 (top-right at edge)
      Positioned(
        right: 0,
        top: size.height * 0.08,
        child: Image.asset(
          'assets/images/cloud555.png',
          width: size.width * 0.28,
        ),
      ),

      // cloud66 (middle-right at edge)
      Positioned(
        right: 0,
        bottom: size.height * 0.22,
        child: Image.asset(
          'assets/images/cloud66.png',
          width: size.width * 0.42,
        ),
      ),

      // cloud77 (lower-right at edge)
      Positioned(
        right: 0,
        bottom: size.height * 0.15,
        child: Image.asset(
          'assets/images/cloud77.png',
          width: size.width * 0.38,
        ),
      ),

      // ---- PLANE ----
      Align(
        alignment: const Alignment(0, -0.15),
        child: Image.asset(
          'assets/images/planee.png',
          width: size.width * 0.80,
        ),
      ),
    ],
  ),
),
                    ],
                  ),
                ),
              ],
            ),
          ),

          // 🔹 BOTTOM PART: buttons, no overlay with background
          Container(
            width: double.infinity,
            color: Colors.white,
            child: SafeArea(
              top: false,
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 24,
                  vertical: 24,
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Login button
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
                      child: CustomButton(
                        text: 'login'.tr,
                        onPressed: () => Get.toNamed(AppRoutes.LOGIN),
                        backgroundColor: const Color(0xFF1e5a8e),
                        textColor: Colors.white,
                      ),
                    ),

                    const SizedBox(height: 16),

                    // Sign up button (outline style)
                    Container(
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
                      child: InkWell(
                        onTap: () => Get.toNamed(AppRoutes.REGISTER),
                        child: Container(
                          width: double.infinity,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          decoration: BoxDecoration(
                            border: Border.all(
                              color: const Color(0xFF1e5a8e),
                              width: 2,
                            ),
                            borderRadius: BorderRadius.circular(25),
                          ),
                          child: Text(
                            'signup'.tr,
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                              color: Color(0xFF1e5a8e),
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 16),

                    // Continue as guest
                    TextButton(
                      onPressed: () => Get.toNamed(AppRoutes.HOME),
                      child: Text(
                        'continue_guest'.tr,
                        style: const TextStyle(
                          color: Color(0xFF666666),
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
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
    );
  }
}
