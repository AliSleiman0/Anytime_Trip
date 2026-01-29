import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:country_picker/country_picker.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import '../../../app/routes/app_routes.dart';
import '../../../core/utils/helpers.dart';
import '../../../core/storage/storage_service.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/endpoints.dart';

class CompleteProfilePage extends StatefulWidget {
  const CompleteProfilePage({super.key});

  @override
  State<CompleteProfilePage> createState() => _CompleteProfilePageState();
}

class _CompleteProfilePageState extends State<CompleteProfilePage> {
  final phoneController = TextEditingController();
  final _storage = StorageService();
  final _apiClient = ApiClient();
  
  String? selectedSex;
  String? selectedCountry;
  bool isLoading = false;
  String userName = '';
  String userEmail = '';

  @override
  void initState() {
    super.initState();
    _loadUserData();
  }

  Future<void> _loadUserData() async {
    final user = await _storage.getUser();
    if (user != null) {
      setState(() {
        userName = user['name'] ?? '';
        userEmail = user['email'] ?? '';
      });
    }
  }

  Future<void> _completeProfile() async {
    // Validation
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

    setState(() => isLoading = true);

    try {
      await _apiClient.put(
        Endpoints.updateProfile,
        data: {
          'phone_number': phoneController.text,
          'sex': selectedSex,
          'country': selectedCountry,
        },
      );

      if (!mounted) return;
      
      Helpers.showSnackbar('Success', 'Profile completed successfully!', isError: false);
      Get.offAllNamed(AppRoutes.HOME);
    } catch (e) {
      if (!mounted) return;
      Helpers.showSnackbar('Error', 'Failed to update profile: $e', isError: true);
    } finally {
      setState(() => isLoading = false);
    }
  }

  @override
  void dispose() {
    phoneController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;

    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: const Color(0xFF1e5a8e),
        elevation: 0,
        title: const Text(
          'Complete Your Profile',
          style: TextStyle(color: Colors.white),
        ),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 20),
            const Text(
              'Welcome!',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.bold,
                color: Color(0xFF1e5a8e),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Please complete your profile to continue',
              style: TextStyle(
                fontSize: 16,
                color: Colors.grey[600],
              ),
            ),
            const SizedBox(height: 32),
            
            // Display Name (read-only)
            TextField(
              enabled: false,
              decoration: InputDecoration(
                labelText: 'Name',
                hintText: userName,
                filled: true,
                fillColor: Colors.grey[100],
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
            ),
            const SizedBox(height: 16),
            
            // Display Email (read-only)
            TextField(
              enabled: false,
              decoration: InputDecoration(
                labelText: 'Email',
                hintText: userEmail,
                filled: true,
                fillColor: Colors.grey[100],
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
            ),
            const SizedBox(height: 16),
            
            // Phone Number
            IntlPhoneField(
              controller: phoneController,
              decoration: InputDecoration(
                labelText: 'Phone Number',
                hintText: 'Enter phone number',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
                fillColor: Colors.white,
              ),
              initialCountryCode: 'US',
              onChanged: (phone) {
                // Handle phone number change
              },
            ),
            const SizedBox(height: 16),
            
            // Sex Dropdown
            DropdownButtonFormField<String>(
              value: selectedSex,
              decoration: InputDecoration(
                labelText: 'Sex',
                hintText: 'Select your sex',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
                fillColor: Colors.white,
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
            const SizedBox(height: 16),
            
            // Country Picker
            GestureDetector(
              onTap: () {
                showCountryPicker(
                  context: context,
                  showPhoneCode: false,
                  onSelect: (Country country) {
                    setState(() {
                      selectedCountry = country.name;
                    });
                  },
                );
              },
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 18),
                decoration: BoxDecoration(
                  border: Border.all(color: Colors.grey[300]!),
                  borderRadius: BorderRadius.circular(12),
                  color: Colors.white,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      selectedCountry ?? 'Select your country',
                      style: TextStyle(
                        fontSize: 16,
                        color: selectedCountry == null ? Colors.grey[600] : Colors.black87,
                      ),
                    ),
                    const Icon(Icons.arrow_drop_down),
                  ],
                ),
              ),
            ),
            
            SizedBox(height: size.height * 0.08),
            
            // Complete Profile Button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: isLoading ? null : _completeProfile,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: isLoading
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                          color: Colors.white,
                          strokeWidth: 2,
                        ),
                      )
                    : const Text(
                        'Complete Profile',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
