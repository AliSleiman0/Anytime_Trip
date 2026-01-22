# Flutter Native Splash Setup Guide

## Quick Setup

Your Flutter app now has `flutter_native_splash` properly configured for your preloader logo. Follow these steps:

### 1. Install Dependencies
```bash
cd app
flutter pub get
```

### 2. Generate Splash Screen
Run this command to generate native splash screens for all platforms:
```bash
dart run flutter_native_splash:create
```

If you're using the separate config file:
```bash
dart run flutter_native_splash:create --path=../flutter_native_splash.yaml
```

### 3. (Optional) Preserve Splash During App Initialization

If you want the splash screen to stay visible while your app initializes, update `lib/main.dart`:

```dart
import 'package:flutter_native_splash/flutter_native_splash.dart';

void main() {
  WidgetsBinding widgetsBinding = WidgetsFlutterBinding.ensureInitialized();
  FlutterNativeSplash.preserve(widgetsBinding: widgetsBinding);
  runApp(const MyApp());
}

// Later, when initialization is complete:
void initializeApp() {
  // ... your initialization code ...
  FlutterNativeSplash.remove();
}
```

## Configuration Files

### File 1: `pubspec.yaml`
Contains the basic splash configuration as a dev dependency. This is the recommended approach.

### File 2: `flutter_native_splash.yaml`
Standalone configuration file with detailed comments and additional options. Use this for:
- More complex setups
- Multiple flavors/environments
- Platform-specific customization

## Current Configuration

**Splash Image:** `assets/images/main-logo.png`
**Background Color:** White (#FFFFFF) / Dark (#1a1a1a)
**Fullscreen:** Yes (status bar hidden)
**Supported Platforms:** 
- ✅ Android (all versions including 12+)
- ✅ iOS
- ✅ Web

## Customizing Your Splash Screen

### Change the Logo Image
Edit the `image` field in either `pubspec.yaml` or `flutter_native_splash.yaml`:
```yaml
image: assets/images/your-logo.png
```

### Change Background Color
```yaml
color: "#YOUR_HEX_COLOR"  # Light mode
color_dark: "#YOUR_HEX_COLOR"  # Dark mode
```

### Add a Branding Image (smaller logo at bottom)
```yaml
branding: assets/images/branding.png
branding_mode: bottom
branding_bottom_padding: 24
```

### Platform-Specific Settings
```yaml
color_ios: "#FFFFFF"
color_android: "#FFFFFF"
color_web: "#FFFFFF"
```

## Troubleshooting

### Splash screen doesn't appear
1. Clean your project: `flutter clean`
2. Get dependencies: `flutter pub get`
3. Regenerate splash: `dart run flutter_native_splash:create`
4. Rebuild app: `flutter run`

### Android 12+ specific issues
- The splash screen may not appear when launching from Android Studio on API 31
- It should appear when using the launcher icon
- Non-Google launchers may not display correctly
- Splash screen won't show when launched from notifications

### iOS setup for multiple flavors
See the separate `ios-flavors-setup.md` file for detailed Xcode configuration

## Advanced: Multiple Flavors

If your app has multiple flavors/environments (Development, Staging, Production):

1. Create separate config files:
   - `flutter_native_splash-development.yaml`
   - `flutter_native_splash-staging.yaml`
   - `flutter_native_splash-production.yaml`

2. Generate for specific flavor:
   ```bash
   dart run flutter_native_splash:create --flavor production
   ```

3. Generate for all flavors:
   ```bash
   dart run flutter_native_splash:create --all-flavors
   ```

## Image Requirements

### Standard (All platforms except Android 12+)
- Format: PNG
- Recommended: 1152×1152 px
- Fits within circle: 768px diameter
- 4x pixel density

### Android 12+
- With background: 960×960 px (fits within circle 640px diameter)
- Without background: 1152×1152 px (fits within circle 768px diameter)
- One-third of foreground is masked

### Branding (optional)
- Dimensions: 800×320 px
- Format: PNG

## Useful Commands

| Command | Description |
|---------|-------------|
| `dart run flutter_native_splash:create` | Generate splash screens |
| `dart run flutter_native_splash:remove` | Remove and restore default white splash |
| `flutter clean; flutter pub get` | Reset and reinstall dependencies |
| `dart run flutter_native_splash:create --flavor staging` | Generate for specific flavor |
| `dart run flutter_native_splash:create --all-flavors` | Generate for all flavors |

## Best Practices

✅ **Do:**
- Use PNG format images
- Test on actual devices, not just emulators
- Provide both light and dark mode versions
- Make images 4x pixel density
- Use the branding element for secondary logos

❌ **Don't:**
- Use JPEG or other formats for splash images
- Rely on Android Studio emulator for testing (API 31 bug)
- Use background images on Android 12+
- Place important design elements outside the safe circle area
- Forget to run `flutter pub get` after changes

## References

- [flutter_native_splash Pub.dev](https://pub.dev/packages/flutter_native_splash)
- [Android Splash Screen Documentation](https://developer.android.com/guide/topics/ui/splash-screen)
- [Apple iOS Launch Screen Documentation](https://developer.apple.com/documentation/uikit/uiview/contentmode)

