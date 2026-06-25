import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:firebase_core/firebase_core.dart';
import 'app/app.dart';
import 'firebase_options.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Firebase with timeout and error handling
  // Don't block app startup if Firebase fails
  _initializeFirebase();
  
  // Immersive mode - hide system UI, appear on swipe
  SystemChrome.setEnabledSystemUIMode(SystemUiMode.immersiveSticky);
  
  runApp(const MyApp());
}

Future<void> _initializeFirebase() async {
  try {
    // Check if Firebase is already initialized (e.g., from iOS native config)
    if (Firebase.apps.isEmpty) {
      // Initialize with timeout to prevent hanging
      await Firebase.initializeApp(
        options: DefaultFirebaseOptions.currentPlatform,
      ).timeout(
        const Duration(seconds: 10),
        onTimeout: () {
          debugPrint('Firebase initialization timed out');
          throw TimeoutException('Firebase initialization timed out');
        },
      );
      debugPrint('Firebase initialized successfully');
    } else {
      debugPrint('Firebase already initialized');
    }
  } catch (e) {
    // Log error but don't block app startup
    debugPrint('Firebase initialization error: $e');
  }
}
