import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:flutter_svg/flutter_svg.dart';
import '../../../core/widgets/custom_button.dart';
import '../../../app/routes/app_routes.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> with TickerProviderStateMixin {
  String selectedLanguage = 'en'; // Default to English
  
  // Animation dated: 2026-02-01
  late AnimationController _planeController;
  late AnimationController _exitController;
  
  late Animation<Offset> _planeAnimation;
  late Animation<Offset> _exitPlaneAnimation;
  late Animation<double> _buttonsFadeAnimation;

  @override
  void initState() {
    super.initState();
    _setupAnimations();
    // Ensure animations are reset to initial state
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _resetAnimationsState();
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // Reset animations when returning to this page
    final route = ModalRoute.of(context);
    if (route != null && route.isCurrent && _exitController.value > 0) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _resetAnimationsState();
      });
    }
  }

  void _resetAnimationsState() {
    // Reset animations to initial state
    _exitController.reset();
    _planeController.reset();
    // Start entrance animation after delay
    Future.delayed(const Duration(milliseconds: 600), () {
      if (mounted) _planeController.forward();
    });
  }

  void _setupAnimations() {
    // Animation for plane - from below screen to top (entrance)
    _planeController = AnimationController(
      duration: const Duration(milliseconds: 3800),
      vsync: this,
    );
    _planeAnimation = Tween<Offset>(
      begin: const Offset(0, 3), // Start from well below the screen
      end: const Offset(0, 0),   // Move to final position
    ).animate(CurvedAnimation(parent: _planeController, curve: Curves.easeOut));

    // Animation for exit - plane moves up and out
    _exitController = AnimationController(
      duration: const Duration(milliseconds: 3800),
      vsync: this,
    );
    _exitPlaneAnimation = Tween<Offset>(
      begin: const Offset(0, 0),   // Start from current position
      end: const Offset(0, -3),    // Move upward off screen
    ).animate(CurvedAnimation(parent: _exitController, curve: Curves.easeIn));

    // Fade animation for buttons
    _buttonsFadeAnimation = Tween<double>(
      begin: 1.0,
      end: 0.0,
    ).animate(CurvedAnimation(parent: _exitController, curve: Curves.easeIn));

    // Start animation for plane on entrance
    Future.delayed(const Duration(milliseconds: 600), () {
      if (mounted) _planeController.forward();
    });
  }

  @override
  void dispose() {
    _planeController.dispose();
    _exitController.dispose();
    super.dispose();
  }

  void _animateAndNavigate(String routeName) async {
    // Start exit animation (plane moves up, buttons fade)
    await _exitController.forward();
    // Navigate after animation completes
    Get.toNamed(routeName);
  }

  Future<void> _resetAnimations() async {
    // Reset to initial state
    _resetAnimationsState();
  }

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

    return WillPopScope(
      onWillPop: () async {
        _resetAnimations();
        return true;
      },
      child: Scaffold(
        body: Column(
          children: [
          // 🔹 TOP PART: background + logo + plane + clouds
          Expanded(
            child: Stack(
              children: [
                // Background left - STATIC
                Positioned.fill(
                  child: Image.asset(
                    'assets/images/leftback.png',
                    fit: BoxFit.cover,
                  ),
                ),
                // Background right overlay - STATIC
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
                            const Text(
                              'Travel',
                              style: TextStyle(
                                fontSize: 36,
                                fontWeight: FontWeight.bold,
                                color: Colors.white,
                                letterSpacing: -0.5,
                              ),
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
      AnimatedBuilder(
        animation: Listenable.merge([_planeController, _exitController]),
        builder: (context, child) {
          Offset offset;
          if (_exitController.value > 0) {
            offset = _exitPlaneAnimation.value;
          } else {
            offset = _planeAnimation.value;
          }
          return Transform.translate(
            offset: Offset(offset.dx * 0, offset.dy * MediaQuery.of(context).size.height),
            child: Align(
              alignment: const Alignment(0, -0.15),
              child: Image.asset(
                'assets/images/planee.png',
                width: size.width * 0.80,
              ),
            ),
          );
        },
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

          // 🔹 BOTTOM PART: buttons move with plane on entrance, fade on exit
          AnimatedBuilder(
            animation: Listenable.merge([_planeController, _exitController]),
            builder: (context, child) {
              // For entrance, use plane animation offset
              Offset offset = _planeAnimation.value;
              // For exit, use fade animation instead of movement
              double opacity = 1.0 - _exitController.value;
              
              return Transform.translate(
                offset: _exitController.value > 0 
                    ? Offset.zero 
                    : Offset(0, offset.dy * MediaQuery.of(context).size.height),
                child: Opacity(
                  opacity: opacity,
                  child: Container(
                    width: double.infinity,
                    color: Colors.transparent,
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
                        child: InkWell(
                          onTap: () => _animateAndNavigate(AppRoutes.LOGIN),
                          child: Container(
                            width: double.infinity,
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            decoration: BoxDecoration(
                              color: const Color(0xFF1e5a8e),
                              borderRadius: BorderRadius.circular(25),
                            ),
                            child: Text(
                              'login'.tr,
                              textAlign: TextAlign.center,
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
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
                          onTap: () => _animateAndNavigate(AppRoutes.REGISTER),
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
                        onPressed: () => _animateAndNavigate(AppRoutes.HOME),
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
          ),
          );
        },
      ),
        ],
      ),
      ),
    );
  }
}
