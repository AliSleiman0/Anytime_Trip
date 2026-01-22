import 'package:flutter/material.dart';
import '../../../core/widgets/unified_ui_components.dart';

class CarCheckout extends StatefulWidget {
  final Map<String, dynamic> car;
  final String pickupLocation;
  final String dropoffLocation;
  final DateTime? pickupDate;
  final DateTime? dropoffDate;
  final String pickupTime;
  final String dropoffTime;
  final String selectedExtra;

  const CarCheckout({
    super.key,
    required this.car,
    required this.pickupLocation,
    required this.dropoffLocation,
    this.pickupDate,
    this.dropoffDate,
    required this.pickupTime,
    required this.dropoffTime,
    required this.selectedExtra,
  });

  @override
  State<CarCheckout> createState() => _CarCheckoutState();
}

class _CarCheckoutState extends State<CarCheckout> {
  String _selectedPaymentMethod = 'Whish Money';
  String _selectedCountryCode = '';
  String _selectedExpiryMonth = 'Month';
  String _selectedExpiryYear = 'Year';
  String _selectedBillingCountry = '';
  String _selectedState = '';

  final TextEditingController _driverNameController = TextEditingController();
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _nameOnCardController = TextEditingController();
  final TextEditingController _cardNumberController = TextEditingController();
  final TextEditingController _securityCodeController = TextEditingController();
  final TextEditingController _billingAddress1Controller = TextEditingController();
  final TextEditingController _billingAddress2Controller = TextEditingController();
  final TextEditingController _cityController = TextEditingController();
  final TextEditingController _zipCodeController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: UnifiedAppBarWithSubtitle(
        title: widget.car['category'] ?? 'Midsize SUV',
        subtitle: _getFormattedDateRange(),
        onBackPressed: () => Navigator.pop(context),
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Car section
                  const Text(
                    'Car',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Colors.black87,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey.shade300),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              _getLocationName(widget.pickupLocation),
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black87,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              '${widget.pickupTime} - ${widget.dropoffTime}',
                              style: const TextStyle(
                                fontSize: 12,
                                color: Colors.black54,
                              ),
                            ),
                          ],
                        ),
                        const Text(
                          '5 Passengers',
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.black54,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Price Summary
                  const Text(
                    'Price Summary',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildPriceLine('Due today', '\$40.99'),
                  const Divider(height: 20),
                  _buildPriceLine('Due today', '\$35.27'),
                  const Divider(height: 20),
                  _buildPriceLine('Total', '\$76.26', isTotal: true),
                  const SizedBox(height: 24),

                  // Who's Driving?
                  const Text(
                    'Who\'s Driving?',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildTextField('Driver Name', _driverNameController),
                  const SizedBox(height: 12),
                  _buildTextField('Name', _nameController),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildDropdown('Country/ Territory Code', _selectedCountryCode, (val) {
                          setState(() {
                            _selectedCountryCode = val ?? '';
                          });
                        }),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('Phone Number', _phoneController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // Pay With
                  const Text(
                    'Pay With',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildPaymentOption('Whish Money', 'whish.png'),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildPaymentOption('OMT Pay', 'omt.png'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildPaymentOption('Crypto Pay', 'crypto.png'),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildPaymentOption('Visa', 'visa.png'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // Card Details
                  const Text(
                    'CREDIT CARD INFO',
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                      color: Colors.black54,
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildTextField('Name on Card', _nameOnCardController),
                  const SizedBox(height: 12),
                  _buildTextField('0000 0000 0000 0000', _cardNumberController, keyboardType: TextInputType.number),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        flex: 2,
                        child: Row(
                          children: [
                            Expanded(
                              child: _buildDropdown('Month', _selectedExpiryMonth, (val) {
                                setState(() {
                                  _selectedExpiryMonth = val ?? 'Month';
                                });
                              }, items: ['Month', '01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12']),
                            ),
                            const Padding(
                              padding: EdgeInsets.symmetric(horizontal: 8),
                              child: Text('-', style: TextStyle(fontSize: 18)),
                            ),
                            Expanded(
                              child: _buildDropdown('Year', _selectedExpiryYear, (val) {
                                setState(() {
                                  _selectedExpiryYear = val ?? 'Year';
                                });
                              }, items: ['Year', '2024', '2025', '2026', '2027', '2028', '2029', '2030']),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('Security Code', _securityCodeController, keyboardType: TextInputType.number),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  _buildDropdown('Country/ Territory Code', _selectedBillingCountry, (val) {
                    setState(() {
                      _selectedBillingCountry = val ?? '';
                    });
                  }),
                  const SizedBox(height: 12),
                  _buildTextField('Billing Address 1', _billingAddress1Controller),
                  const SizedBox(height: 12),
                  _buildTextField('Billing Address 2', _billingAddress2Controller),
                  const SizedBox(height: 12),
                  _buildTextField('City', _cityController),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildDropdown('State', _selectedState, (val) {
                          setState(() {
                            _selectedState = val ?? '';
                          });
                        }),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildTextField('ZIP code', _zipCodeController, keyboardType: TextInputType.number),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                ],
              ),
            ),
          ),

          // Complete Reserving button
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white,
              boxShadow: [
                BoxShadow(
                  color: Colors.grey.shade200,
                  blurRadius: 4,
                  offset: const Offset(0, -2),
                ),
              ],
            ),
            child: SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: () {
                  _showConfirmationDialog();
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: const Text(
                  'Complete Reserving',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPriceLine(String label, String amount, {bool isTotal = false}) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 14,
            color: isTotal ? const Color(0xFFD32F2F) : Colors.black87,
            fontWeight: isTotal ? FontWeight.w600 : FontWeight.w400,
          ),
        ),
        Text(
          amount,
          style: TextStyle(
            fontSize: 14,
            color: isTotal ? const Color(0xFFD32F2F) : Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }

  Widget _buildTextField(String hint, TextEditingController controller, {TextInputType? keyboardType}) {
    return SizedBox(
      height: 45,
      child: Container(
        decoration: const BoxDecoration(
          border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
        ),
        child: TextField(
          controller: controller,
          keyboardType: keyboardType,
          decoration: InputDecoration(
            hintText: hint,
            hintStyle: const TextStyle(fontSize: 13, color: Colors.grey),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
            border: InputBorder.none,
            enabledBorder: InputBorder.none,
            focusedBorder: InputBorder.none,
          ),
        ),
      ),
    );
  }

  Widget _buildDropdown(String hint, String value, Function(String?) onChanged, {List<String>? items}) {
    return Container(
      height: 45,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: const BoxDecoration(
        border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
      ),
      child: Row(
        children: [
          Expanded(
            child: Text(
              value.isEmpty ? hint : value,
              style: TextStyle(
                fontSize: 13,
                color: value.isEmpty ? Colors.grey : Colors.black87,
              ),
            ),
          ),
          const Icon(Icons.arrow_drop_down, color: Color(0xFF1e5a8e), size: 20),
        ],
      ),
    );
  }

  Widget _buildPaymentOption(String name, String iconPath) {
    final isSelected = _selectedPaymentMethod == name;
    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedPaymentMethod = name;
        });
      },
      child: Container(
        height: 50,
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 10),
        decoration: BoxDecoration(
          border: Border.all(
            color: isSelected ? const Color(0xFF1e5a8e) : Colors.grey.shade300,
            width: isSelected ? 2 : 1,
          ),
          borderRadius: BorderRadius.circular(6),
          color: Colors.white,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset(
              'assets/images/$iconPath',
              width: 20,
              height: 20,
              errorBuilder: (context, error, stackTrace) {
                return const Icon(Icons.payment, size: 20, color: Colors.grey);
              },
            ),
            const SizedBox(width: 6),
            Flexible(
              child: Text(
                name,
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                  color: isSelected ? const Color(0xFF1e5a8e) : Colors.black87,
                ),
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _getFormattedDateRange() {
    if (widget.pickupDate == null || widget.dropoffDate == null) {
      return 'Sun, Oct 19, 1:30 am - Sun, Oct 19, 4:30 am';
    }
    
    final pickupStr = _formatDateTime(widget.pickupDate!, widget.pickupTime);
    final dropoffStr = _formatDateTime(widget.dropoffDate!, widget.dropoffTime);
    
    return '$pickupStr - $dropoffStr';
  }

  String _formatDateTime(DateTime date, String time) {
    const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    
    final dayName = days[date.weekday % 7];
    final monthName = months[date.month - 1];
    final day = date.day;
    
    return '$dayName, $monthName $day, $time';
  }

  String _getLocationName(String location) {
    if (location.contains('-')) {
      final code = location.split('-')[0].toUpperCase();
      return '$code Airport';
    }
    return '$location Airport';
  }

  void _showConfirmationDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Container(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Close button
                Align(
                  alignment: Alignment.topRight,
                  child: GestureDetector(
                    onTap: () {
                      Navigator.of(context).pop();
                      Navigator.of(context).pop();
                      Navigator.of(context).pop();
                    },
                    child: Container(
                      padding: const EdgeInsets.all(4),
                      decoration: const BoxDecoration(
                        color: Color(0xFFD32F2F),
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(
                        Icons.close,
                        color: Colors.white,
                        size: 20,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                
                // Car image
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.grey.shade300, width: 2),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Image.asset(
                    'assets/images/car.png',
                    width: 150,
                    height: 100,
                    fit: BoxFit.contain,
                    errorBuilder: (context, error, stackTrace) {
                      return const Icon(Icons.directions_car, size: 100, color: Colors.blue);
                    },
                  ),
                ),
                const SizedBox(height: 20),
                
                // Success message
                const Text(
                  'Your rental has been completed successfully.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Your rental reference is:',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 12),
                
                // Reference number
                const Text(
                  '1872306517801',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w700,
                    color: Colors.black,
                  ),
                ),
                const SizedBox(height: 16),
                
                // Confirmation email text
                const Text(
                  'You will receive a confirmation email including\nyour itinerary details',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.black87,
                    height: 1.4,
                  ),
                ),
                const SizedBox(height: 24),
                
                // Done button
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.of(context).pop();
                      Navigator.of(context).pop();
                      Navigator.of(context).pop();
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF1e5a8e),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(25),
                      ),
                    ),
                    child: const Text(
                      'Done',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  @override
  void dispose() {
    _driverNameController.dispose();
    _nameController.dispose();
    _phoneController.dispose();
    _nameOnCardController.dispose();
    _cardNumberController.dispose();
    _securityCodeController.dispose();
    _billingAddress1Controller.dispose();
    _billingAddress2Controller.dispose();
    _cityController.dispose();
    _zipCodeController.dispose();
    super.dispose();
  }
}
