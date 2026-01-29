import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:dio/dio.dart' as dio;
import 'package:image_picker/image_picker.dart';
import 'dart:io';
import '../../../core/widgets/unified_ui_components.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';
import '../../../core/utils/helpers.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _lastNameController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final _storage = StorageService();
  final _apiClient = ApiClient();
  final _imagePicker = ImagePicker();
  
  String _selectedNationality = 'Lebanese';
  String _selectedSex = 'Male';
  String _detectedResidency = 'Detecting...';
  bool _residencyLoaded = false;
  bool _isLoading = true;
  bool _isSaving = false;
  bool _isUploadingImage = false;
  String? _profileImageUrl;
  
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
    _loadUserProfile();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _detectResidencyFromIP();
    });
  }

  Future<void> _loadUserProfile() async {
    setState(() => _isLoading = true);
    try {
      final user = await _storage.getUser();
      if (user != null) {
        final fullName = user['name'] ?? '';
        final nameParts = fullName.split(' ');
        final country = user['country'] ?? '';
        
        setState(() {
          _firstNameController.text = nameParts.isNotEmpty ? nameParts[0] : '';
          _profileImageUrl = user['profile_image'];
          _lastNameController.text = nameParts.length > 1 ? nameParts.sublist(1).join(' ') : '';
          _emailController.text = user['email'] ?? '';
          _phoneController.text = user['phone_number'] ?? '';
          
          // If country is not in the list, add it
          if (country.isNotEmpty && !_nationalities.contains(country)) {
            _nationalities.insert(0, country);
          }
          _selectedNationality = country.isNotEmpty ? country : 'Lebanese';
          
          _selectedSex = user['sex'] ?? 'Male';
        });
      }
    } catch (e) {
      print('Error loading user profile: $e');
    } finally {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _saveProfile() async {
    setState(() => _isSaving = true);
    try {
      final fullName = '${_firstNameController.text.trim()} ${_lastNameController.text.trim()}'.trim();
      
      await _apiClient.put(
        Endpoints.updateProfile,
        data: {
          'phone_number': _phoneController.text,
          'sex': _selectedSex,
          'country': _selectedNationality,
        },
      );
      
      // Update local storage
      final user = await _storage.getUser();
      if (user != null) {
        user['name'] = fullName;
        user['phone_number'] = _phoneController.text;
        user['sex'] = _selectedSex;
        user['country'] = _selectedNationality;
        await _storage.saveUser(user);
      }
      
      if (!mounted) return;
      Helpers.showSnackbar('Success', 'Profile updated successfully!', isError: false);
    } catch (e) {
      if (!mounted) return;
      Helpers.showSnackbar('Error', 'Failed to update profile: $e', isError: true);
    } finally {
      setState(() => _isSaving = false);
    }
  }

  Future<void> _pickAndUploadImage(ImageSource source) async {
    try {
      final XFile? pickedFile = await _imagePicker.pickImage(
        source: source,
        maxWidth: 1024,
        maxHeight: 1024,
        imageQuality: 85,
      );

      if (pickedFile == null) return;

      setState(() => _isUploadingImage = true);

      // Create multipart form data
      final formData = dio.FormData.fromMap({
        'image': await dio.MultipartFile.fromFile(
          pickedFile.path,
          filename: pickedFile.name,
        ),
      });

      // Upload image
      final response = await _apiClient.post(
        Endpoints.uploadProfileImage,
        data: formData,
      );

      if (response.data['profile_image'] != null) {
        setState(() {
          _profileImageUrl = response.data['profile_image'];
        });

        // Update local storage
        final user = await _storage.getUser();
        if (user != null) {
          user['profile_image'] = _profileImageUrl;
          await _storage.saveUser(user);
        }

        if (!mounted) return;
        Helpers.showSnackbar('Success', 'Profile image updated!', isError: false);
      }
    } catch (e) {
      if (!mounted) return;
      Helpers.showSnackbar('Error', 'Failed to upload image: $e', isError: true);
    } finally {
      setState(() => _isUploadingImage = false);
    }
  }

  void _showImageSourceDialog() {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (BuildContext context) {
        return SafeArea(
          child: Wrap(
            children: [
              ListTile(
                leading: const Icon(Icons.photo_camera, color: Color(0xFF1e5a8e)),
                title: const Text('Take Photo'),
                onTap: () {
                  Navigator.pop(context);
                  _pickAndUploadImage(ImageSource.camera);
                },
              ),
              ListTile(
                leading: const Icon(Icons.photo_library, color: Color(0xFF1e5a8e)),
                title: const Text('Choose from Gallery'),
                onTap: () {
                  Navigator.pop(context);
                  _pickAndUploadImage(ImageSource.gallery);
                },
              ),
            ],
          ),
        );
      },
    );
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

      final dioClient = dio.Dio();
      
      // Set timeout to 10 seconds
      dioClient.options.connectTimeout = const Duration(seconds: 10);
      dioClient.options.receiveTimeout = const Duration(seconds: 10);
      
      try {
        // Try primary API: ip-api.com (45 requests per minute free tier)
        final response = await dioClient.get('http://ip-api.com/json/');
        
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
        final response = await dioClient.get('https://ipinfo.io/json');
        
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
    if (_isLoading) {
      return const Scaffold(
        backgroundColor: Colors.white,
        body: Center(
          child: CircularProgressIndicator(
            color: Color(0xFF1e5a8e),
          ),
        ),
      );
    }
    
    final userName = _firstNameController.text.isNotEmpty 
        ? _firstNameController.text 
        : 'User';
    
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
            Text(
              'Hello, $userName',
              style: const TextStyle(
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
                child: _profileImageUrl != null && _profileImageUrl!.isNotEmpty
                    ? Image.network(
                        '${Endpoints.baseUrl.replaceAll('/api', '')}$_profileImageUrl',
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return Image.asset(
                            'assets/images/defaultp.png',
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) {
                              return const Icon(Icons.person, size: 20, color: Color(0xFF1e5a8e));
                            },
                          );
                        },
                      )
                    : Image.asset(
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
                      child: _isUploadingImage
                          ? const Center(
                              child: CircularProgressIndicator(
                                color: Color(0xFF1e5a8e),
                              ),
                            )
                          : _profileImageUrl != null && _profileImageUrl!.isNotEmpty
                              ? Image.network(
                                  '${Endpoints.baseUrl.replaceAll('/api', '')}$_profileImageUrl',
                                  fit: BoxFit.cover,
                                  errorBuilder: (context, error, stackTrace) {
                                    return Image.asset(
                                      'assets/images/defaultp.png',
                                      fit: BoxFit.cover,
                                      errorBuilder: (context, error, stackTrace) {
                                        return const Icon(Icons.person, size: 60, color: Color(0xFF1e5a8e));
                                      },
                                    );
                                  },
                                )
                              : Image.asset(
                                  'assets/images/defaultp.png',
                                  fit: BoxFit.cover,
                                  errorBuilder: (context, error, stackTrace) {
                                    return const Icon(Icons.person, size: 60, color: Color(0xFF1e5a8e));
                                  },
                                ),
                    ),
                  ),
                  GestureDetector(
                    onTap: _isUploadingImage ? null : _showImageSourceDialog,
                    child: Container(
                      decoration: BoxDecoration(
                        color: _isUploadingImage ? Colors.grey : const Color(0xFFD32F2F),
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

              // Save Button
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isSaving ? null : _saveProfile,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1e5a8e),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: _isSaving
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            color: Colors.white,
                            strokeWidth: 2,
                          ),
                        )
                      : const Text(
                          'Save Changes',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                ),
              ),
              const SizedBox(height: 16),

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
