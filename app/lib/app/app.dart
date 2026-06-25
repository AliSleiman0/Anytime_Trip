import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'routes/app_pages.dart';
import 'routes/app_routes.dart';
import 'theme/app_theme.dart';
import '../core/localization/app_translations.dart';
import '../core/widgets/global_chatbot_overlay.dart';

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return GetMaterialApp(
      title: 'Travel',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.light,
      initialRoute: AppRoutes.SPLASH_SCREEN,
      getPages: AppPages.pages,
      debugShowCheckedModeBanner: false,
      // Localization setup
      translations: AppTranslations(),
      locale: const Locale('en', 'US'),
      fallbackLocale: const Locale('en', 'US'),
      // Add global chatbot overlay to all screens
      builder: (context, child) {
        return GlobalChatbotOverlay(child: child ?? const SizedBox.shrink());
      },
    );
  }
}
