import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:dio/dio.dart';
import '../../../core/widgets/unified_ui_components.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final TextEditingController _firstNameController = TextEditingController(text: 'Michael');
  final TextEditingController _lastNameController = TextEditingController(text: 'Doe');
  final TextEditingController _emailController = TextEditingController(text: 'MichaelDoe123@mail.com');
  final TextEditingController _phoneController = TextEditingController(text: '00 123 456');
  
  String _selectedNationality = 'Lebanese';
  String _detectedResidency = 'Detecting...';
  bool _residencyLoaded = false;
  
  // Debug info
  String _debugInfo = '';
  String _countryCode = '';
  String _ipAddress = '';

  final List<String> _nationalities = [
    'Lebanese',
    'American',
    'British',
    'French',
    'Spanish',
    'German',
    'Italian',
    'Canadian',
    'Australian',
    'Japanese',
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _detectResidencyFromIP();
    });
  }

  Future<void> _detectResidencyFromIP() async {
    try {
      // Map of country codes to residency names based on common Middle East countries
      final Map<String, String> countryToResidency = {
        'OM': 'Oman',
        'AE': 'UAE',
        'SA': 'Saudi Arabia',
        'KW': 'Kuwait',
        'QA': 'Qatar',
        'BH': 'Bahrain',
        'JO': 'Jordan',
        'LB': 'Lebanon',
        'EG': 'Egypt',
        'IQ': 'Iraq',
        'US': 'United States',
        'GB': 'United Kingdom',
        'FR': 'France',
        'DE': 'Germany',
        'IT': 'Italy',
        'ES': 'Spain',
        'CA': 'Canada',
        'AU': 'Australia',
        'JP': 'Japan',
      };

      final dio = Dio();
      
      // Set timeout to 10 seconds
      dio.options.connectTimeout = const Duration(seconds: 10);
      dio.options.receiveTimeout = const Duration(seconds: 10);
      
      try {
        // Try primary API: ip-api.com (45 requests per minute free tier)
        final response = await dio.get('http://ip-api.com/json/');
        
        if (response.statusCode == 200) {
          final countryCode = response.data['countryCode'] as String;
          final ipAddress = response.data['query'] as String? ?? 'Unknown';
          final detectedCountry = countryToResidency[countryCode] ?? response.data['country'] ?? 'Lebanon';
          
          if (mounted) {
            setState(() {
              _detectedResidency = detectedCountry;
              _residencyLoaded = true;
              _countryCode = countryCode;
              _ipAddress = ipAddress;
              _debugInfo = 'IP: $ipAddress\nCountry Code: $countryCode\nDetected: $detectedCountry\n\n✅ API: ip-api.com';
            });
          }
          print('✅ Detected country from IP: $countryCode -> $detectedCountry');
          return;
        }
      } catch (e) {
        print('⚠️ Primary API failed, trying fallback: $e');
        if (mounted) {
          setState(() {
            _debugInfo = 'Primary API failed, trying backup...';
          });
        }
      }
      
      // Fallback to ipinfo.io
      try {
        final response = await dio.get('https://ipinfo.io/json');
        
        if (response.statusCode == 200) {
          final countryCode = response.data['country'] as String;
          final ipAddress = response.data['ip'] as String? ?? 'Unknown';
          final detectedCountry = countryToResidency[countryCode] ?? 'Lebanon';
          
          if (mounted) {
            setState(() {
              _detectedResidency = detectedCountry;
              _residencyLoaded = true;
              _countryCode = countryCode;
              _ipAddress = ipAddress;
              _debugInfo = 'IP: $ipAddress\nCountry Code: $countryCode\nDetected: $detectedCountry\n\n✅ API: ipinfo.io';
            });
          }
          print('✅ Detected country from fallback API: $countryCode -> $detectedCountry');
          return;
        }
      } catch (e) {
        print('⚠️ Fallback API also failed: $e');
      }
      
      // If all APIs fail, use Lebanon as default
      if (mounted) {
        setState(() {
          _detectedResidency = 'Lebanon';
          _residencyLoaded = true;
          _debugInfo = '⚠️ Both APIs failed\nUsing default: Lebanon';
        });
      }
    } catch (e) {
      print('❌ Error detecting residency: $e');
      if (mounted) {
        setState(() {
          _detectedResidency = 'Lebanon';
          _residencyLoaded = true;
          _debugInfo = '❌ Error: ${e.toString()}\nUsing default: Lebanon';
        });
      }
    }
  }

  void _retestDetection() {
    setState(() {
      _detectedResidency = 'Detecting...';
      _debugInfo = '';
    });
    _detectResidencyFromIP();
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e), size: 24),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              'Hello, User123',
              style: TextStyle(
                color: Colors.black87,
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(width: 12),
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: const Color(0xFFE0E0E0),
              ),
              child: ClipOval(
                child: Image.asset(
                  'assets/images/defaultp.png',
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                  },
                ),
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const SizedBox(height: 16),
              
              // Profile Picture with Edit Button
              Stack(
                alignment: Alignment.bottomRight,
                children: [
                  Container(
                    width: 140,
                    height: 140,
                    decoration: BoxDecoration(
                      color: const Color(0xFFE0E0E0),
                      borderRadius: BorderRadius.circular(100),
                    ),
                    child: ClipOval(
                      child: Image.asset(
                        'assets/images/defaultp.png',
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return const Icon(Icons.person, size: 60, color: Color(0xFF1e5a8e));
                        },
                      ),
                    ),
                  ),
                  Container(
                    decoration: BoxDecoration(
                      color: const Color(0xFFD32F2F),
                      shape: BoxShape.circle,
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.15),
                          blurRadius: 4,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    padding: const EdgeInsets.all(8),
                    child: const Icon(
                      Icons.camera_alt,
                      color: Colors.white,
                      size: 20,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 28),

              // First Name and Last Name Row
              Row(
                children: [
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'First Name',
                      controller: _firstNameController,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildLabeledTextField(
                      label: 'Last Name',
                      controller: _lastNameController,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Email Address
              _buildLabeledTextField(
                label: 'Email Address',
                controller: _emailController,
              ),
              const SizedBox(height: 16),

              // Nationality and Residency Row
              Row(
                children: [
                  Expanded(
                    child: _buildLabeledDropdown(
                      label: 'Nationality',
                      value: _selectedNationality,
                      items: _nationalities,
                      onChanged: (value) {
                        setState(() {
                          _selectedNationality = value!;
                        });
                      },
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildReadOnlyTextField(
                      label: 'Residency (Auto-Detected)',
                      value: _detectedResidency,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Phone Number
              _buildPhoneNumberField(),
              const SizedBox(height: 32),

              // Delete Account Button
              SizedBox(
                width: double.infinity,
                child: OutlinedButton(
                  onPressed: () {},
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    side: const BorderSide(color: Color(0xFFD32F2F), width: 2),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const SizedBox(width: 24),
                      Row(
                        children: [
                          const Icon(Icons.delete_outline, color: Color(0xFFD32F2F), size: 20),
                          const SizedBox(width: 8),
                          const Text(
                            'Delete Account',
                            style: TextStyle(
                              color: Color(0xFFD32F2F),
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                      const Icon(Icons.arrow_forward_ios, color: Color(0xFFD32F2F), size: 18),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLabeledTextField({
    required String label,
    required TextEditingController controller,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: Color(0xFFD32F2F),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: const Color(0xFF336891), width: 2),
            borderRadius: BorderRadius.circular(8),
            color: Colors.white,
          ),
          child: TextField(
            controller: controller,
            decoration: const InputDecoration(
              border: InputBorder.none,
              contentPadding: EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            ),
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: Colors.black87,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildLabeledDropdown({
    required String label,
    required String value,
    required List<String> items,
    required Function(String?) onChanged,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: Color(0xFFD32F2F),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: const Color(0xFF336891), width: 2),
            borderRadius: BorderRadius.circular(8),
            color: Colors.white,
          ),
          padding: const EdgeInsets.only(left: 8, right: 4),
          child: DropdownButton<String>(
            value: value,
            isExpanded: true,
            underline: const SizedBox(),
            icon: const Icon(Icons.check_circle, color: Color(0xFF1e5a8e), size: 20),
            items: items.map((String item) {
              return DropdownMenuItem<String>(
                value: item,
                child: Text(
                  item,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Colors.black87,
                  ),
                ),
              );
            }).toList(),
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }

  Widget _buildReadOnlyTextField({
    required String label,
    required String value,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: Color(0xFFD32F2F),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: const Color(0xFF336891), width: 2),
            borderRadius: BorderRadius.circular(8),
            color: const Color(0xFFF5F5F5),
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  value,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Colors.black87,
                  ),
                ),
                const Icon(
                  Icons.location_on,
                  color: Color(0xFF1e5a8e),
                  size: 18,
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildPhoneNumberField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Phone Number',
          style: TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: Color(0xFFD32F2F),
          ),
        ),
        const SizedBox(height: 8),
        IntlPhoneField(
          controller: _phoneController,
          initialCountryCode: 'LB',
          decoration: InputDecoration(
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(
                color: Color(0xFF336891),
                width: 2,
              ),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(
                color: Color(0xFF336891),
                width: 2,
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: const BorderSide(
                color: Color(0xFF336891),
                width: 2,
              ),
            ),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
            counterText: '',
          ),
          onChanged: (phone) {
            // Handle phone change
          },
        ),
      ],
    );
  }
}
