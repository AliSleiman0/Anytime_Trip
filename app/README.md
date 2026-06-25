# Travel - Flutter App

A Flutter travel application built with **GetX** architecture following clean code principles and best practices.

## 📁 Project Structure

```
lib/
├── main.dart                    # App entry point
├── app/                         # App-level configuration
│   ├── app.dart                # Main app widget
│   ├── routes/                 # Navigation management
│   │   ├── app_routes.dart    # Route constants
│   │   └── app_pages.dart     # Route-page bindings
│   └── theme/                  # App theming
│       └── app_theme.dart     # Light/Dark themes
│
├── core/                        # Shared resources
│   ├── network/                # API communication
│   │   ├── api_client.dart    # Dio HTTP client
│   │   └── endpoints.dart     # API endpoint URLs
│   ├── storage/                # Local data persistence
│   │   └── storage_service.dart
│   ├── utils/                  # Utility functions
│   │   ├── validators.dart    # Form validators
│   │   └── helpers.dart       # Helper functions
│   └── widgets/                # Reusable UI components
│       ├── custom_button.dart
│       └── custom_text_field.dart
│
└── features/                    # Feature modules
    ├── auth/                   # Authentication feature
    ├── product/                # Product browsing feature
    ├── account/                # User account feature
    └── order/                  # Order management feature
```

## 🎯 Feature Structure

Each feature follows a consistent structure for maintainability:

```
features/product/
├── binding/
│   └── product_binding.dart      # Dependency injection
│
├── controller/
│   └── product_controller.dart   # Business logic & state
│
├── view/
│   ├── product_page.dart         # UI screens
│   └── product_details_page.dart
│
├── widget/
│   └── product_card_widget.dart  # Feature-specific widgets
│
├── model/
│   └── product_ui_model.dart     # Data models
│
└── service/
    └── product_api.dart          # API calls
```

## 🏗️ Architecture Pattern

**GetX MVC Pattern:**

1. **Model** - Data structure and business entities
2. **View** - UI components (Pages & Widgets)
3. **Controller** - Business logic and state management
4. **Service** - External data sources (API, Database)
5. **Binding** - Dependency injection

### Flow Example:

```
User Action (View)
    ↓
Controller Method
    ↓
Service API Call
    ↓
Update State (Controller)
    ↓
UI Rebuild (View)
```

## 📦 Core Components

### Routes vs Endpoints

**Routes** (`app/routes/`) - Frontend Navigation
- Define app screen navigation paths
- Example: `/login`, `/product/:id`
- Used with: `Get.toNamed(AppRoutes.LOGIN)`

**Endpoints** (`core/network/`) - Backend API URLs
- Define server API endpoints
- Example: `POST /api/app/login`
- Used with: `_apiClient.post(Endpoints.login)`

### API Client

Centralized HTTP client using Dio:
- Base URL configuration
- Request/Response interceptors
- Authentication token injection
- Error handling

### Storage Service

Local data persistence using GetStorage:
- User authentication tokens
- User profile data
- App preferences
- Theme settings

## 🎨 Theming

Supports light and dark themes:
- Material Design 3
- Centralized color scheme
- Consistent typography
- Custom component themes

## 📱 Features

### 1. Authentication (`features/auth/`)
- Splash screen
- Login
- Registration
- Logout

### 2. Product (`features/product/`)
- Product listing
- Product details
- Search & filtering
- Add to cart

### 3. Account (`features/account/`)
- User profile
- Profile editing
- Settings
- Logout

### 4. Order (`features/order/`)
- Order history
- Order details
- Order tracking

## 🚀 Getting Started

### Prerequisites
- Flutter SDK (>=3.0.0)
- Dart SDK
- VS Code or Android Studio

### Installation

1. **Install dependencies:**
```bash
flutter pub get
```

2. **Run the app:**
```bash
flutter run
```

3. **Build for production:**
```bash
# Android
flutter build apk --release

# iOS
flutter build ios --release
```

## 📚 Dependencies

- **get** - State management & navigation
- **dio** - HTTP client
- **get_storage** - Local storage

## 🔧 Configuration

Update the base URL in `lib/core/network/endpoints.dart`:

```dart
static const String baseUrl = 'http://your-api-url.com/api';
```

## 📖 Code Guidelines

### Naming Conventions
- **Files**: `snake_case.dart`
- **Classes**: `PascalCase`
- **Variables**: `camelCase`
- **Constants**: `UPPER_CASE`

### Feature Development

1. Create feature folder in `features/`
2. Add required subfolders: `binding/`, `controller/`, `view/`, `model/`, `service/`
3. Implement binding for dependency injection
4. Create controller for business logic
5. Add service for API calls
6. Build views (pages)
7. Register routes in `app/routes/app_pages.dart`

### Example - Adding New Feature:

```dart
// 1. Create binding
class SettingsBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => SettingsController());
  }
}

// 2. Create controller
class SettingsController extends GetxController {
  // State and logic
}

// 3. Create view
class SettingsPage extends GetView<SettingsController> {
  // UI
}

// 4. Register route
GetPage(
  name: AppRoutes.SETTINGS,
  page: () => const SettingsPage(),
  binding: SettingsBinding(),
)
```

## 🧪 Testing

```bash
# Run tests
flutter test

# Run with coverage
flutter test --coverage
```

## 📝 Best Practices

✅ Use GetX for state management  
✅ Separate business logic from UI  
✅ Keep controllers thin, services fat  
✅ Use proper error handling  
✅ Follow DRY principle  
✅ Use const constructors where possible  
✅ Implement proper null safety  
✅ Add proper documentation  

## 🤝 Contributing

1. Create feature branch
2. Follow the established structure
3. Write clean, documented code
4. Test your changes
5. Submit pull request

## 📄 License

This project is licensed under the MIT License.
