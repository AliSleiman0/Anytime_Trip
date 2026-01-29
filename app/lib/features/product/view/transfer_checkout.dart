import 'package:flutter/material.dart';
import '../../../core/widgets/unified_ui_components.dart';
import '../../account/view/bookings_page.dart';
import 'home_page.dart';

class TransferCheckout extends StatefulWidget {
  final Map<String, dynamic> transfer;
  final String pickupLocation;
  final String dropoffLocation;
  final String pickupTime;
  final String transferType;
  final int adultsCount;
  final int childrenCount;

  const TransferCheckout({
    super.key,
    required this.transfer,
    required this.pickupLocation,
    required this.dropoffLocation,
    required this.pickupTime,
    required this.transferType,
    this.adultsCount = 1,
    this.childrenCount = 0,
  });

  @override
  State<TransferCheckout> createState() => _TransferCheckoutState();
}

class _TransferCheckoutState extends State<TransferCheckout> {
  String _selectedPaymentMethod = 'Whish Money';
  String _selectedExpiryMonth = 'Month';
  String _selectedExpiryYear = 'Year';

  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _lastNameController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  final TextEditingController _nameOnCardController = TextEditingController();
  final TextEditingController _cardNumberController = TextEditingController();
  final TextEditingController _securityCodeController = TextEditingController();
  final TextEditingController _billingAddress1Controller = TextEditingController();
  final TextEditingController _billingAddress2Controller = TextEditingController();
  final TextEditingController _cityController = TextEditingController();
  final TextEditingController _zipCodeController = TextEditingController();

  int get _totalPassengers => widget.adultsCount + widget.childrenCount;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: UnifiedAppBar(
        title: widget.transfer['category'] ?? 'Transfer',
        onBackPressed: () => Navigator.pop(context),
      ),
      body: Column(
        children: [
          const Divider(height: 1, color: Colors.grey),
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Transfer Card
                  _buildTransferCard(),
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
                  _buildPriceLine('$_totalPassengers passengers', '\$50.00'),
                  const SizedBox(height: 8),
                  _buildPriceLine('Taxes', '\$83.73'),
                  const SizedBox(height: 8),
                  _buildPriceLine('Tax on fees', '\$3.53'),
                  const SizedBox(height: 12),
                  const Align(
                    alignment: Alignment.centerRight,
                    child: Text(
                      'Total: \$602.73',
                      style: TextStyle(
                        fontSize: 14,
                        color: Color(0xFFD32F2F),
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Who's booking?
                  const Text(
                    'Who\'s booking?',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Color(0xFF1e5a8e),
                    ),
                  ),
                  const SizedBox(height: 12),
                  _buildSelectTravelerDropdown(),
                  const SizedBox(height: 12),
                  // Title Row
                  Row(
                    children: [
                      SizedBox(
                        width: 100,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Title', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87)),
                            const SizedBox(height: 4),
                            Container(
                              decoration: const BoxDecoration(
                                border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                              ),
                              child: DropdownButtonFormField<String>(
                                decoration: const InputDecoration(
                                  border: InputBorder.none,
                                  enabledBorder: InputBorder.none,
                                  focusedBorder: InputBorder.none,
                                  isDense: true,
                                  contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                                ),
                                value: 'Mr',
                                items: ['Mr', 'Mrs', 'Ms', 'Dr'].map((String value) {
                                  return DropdownMenuItem<String>(
                                    value: value,
                                    child: Text(value, style: const TextStyle(fontSize: 14)),
                                  );
                                }).toList(),
                                onChanged: (String? newValue) {
                                  // Will be populated from DB later
                                },
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildSimpleTextField('First Name', _firstNameController),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('Last Name', _lastNameController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: _buildDropdownWithCheckRight('Country/ Territory Code'),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('Phone Number', _phoneController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),

                  // Pay With
                  const Center(
                    child: Text(
                      'Pay With',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                        color: Colors.black87,
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
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
                  Row(
                    children: [
                      Expanded(
                        child: _buildSimpleTextField('Name on Card', _nameOnCardController),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Container(
                          height: 45,
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: const Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text('Select Card', style: TextStyle(fontSize: 13, color: Colors.grey)),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        flex: 2,
                        child: Row(
                          children: [
                            Expanded(
                              child: Container(
                                height: 45,
                                decoration: const BoxDecoration(
                                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                                ),
                                child: TextField(
                                  decoration: InputDecoration(
                                    hintText: 'MM',
                                    hintStyle: const TextStyle(fontSize: 13, color: Colors.grey),
                                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
                                    border: InputBorder.none,
                                    enabledBorder: InputBorder.none,
                                    focusedBorder: InputBorder.none,
                                  ),
                                ),
                              ),
                            ),
                            const Padding(
                              padding: EdgeInsets.symmetric(horizontal: 4),
                              child: Text('-', style: TextStyle(fontSize: 16)),
                            ),
                            Expanded(
                              child: Container(
                                height: 45,
                                decoration: const BoxDecoration(
                                  border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                                ),
                                child: TextField(
                                  decoration: InputDecoration(
                                    hintText: 'YY',
                                    hintStyle: const TextStyle(fontSize: 13, color: Colors.grey),
                                    contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
                                    border: InputBorder.none,
                                    enabledBorder: InputBorder.none,
                                    focusedBorder: InputBorder.none,
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('Security Code', _securityCodeController, keyboardType: TextInputType.number),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Container(
                          height: 45,
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: const Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text('Select', style: TextStyle(fontSize: 13, color: Colors.grey)),
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('Billing Address 1', _billingAddress1Controller),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: _buildSimpleTextField('Billing Address 2', _billingAddress2Controller),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('City', _cityController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Container(
                          height: 45,
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          decoration: const BoxDecoration(
                            border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
                          ),
                          child: const Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text('Select', style: TextStyle(fontSize: 13, color: Colors.grey)),
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildSimpleTextField('ZIP code', _zipCodeController),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                ],
              ),
            ),
          ),

          // Complete Booking button
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
                  'Complete Booking',
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

  Widget _buildTransferCard() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFFE8F0F7),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: const Color(0xFF1e5a8e).withOpacity(0.3)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.transfer['category'] ?? 'Midsize SUV',
                  style: const TextStyle(
                    color: Color(0xFFD32F2F),
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  widget.transfer['name'] ?? 'Toyota RAV 4 or similar',
                  style: const TextStyle(color: Colors.black87, fontSize: 12),
                ),
                const SizedBox(height: 2),
                Text(
                  widget.transfer['passengers'] ?? '5 Passengers',
                  style: const TextStyle(color: Colors.black87, fontSize: 12),
                ),
                const SizedBox(height: 2),
                Text(
                  widget.transfer['estimatedTime'] ?? 'Estimated time: 40 min',
                  style: const TextStyle(color: Colors.black87, fontSize: 12),
                ),
                const SizedBox(height: 2),
                Text(
                  widget.transfer['meetGreet'] ?? 'Meet & Greet availability',
                  style: const TextStyle(color: Colors.black87, fontSize: 12),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              ClipRRect(
                borderRadius: BorderRadius.circular(6),
                child: Image.asset(
                  'assets/images/urus.jpeg',
                  width: 70,
                  height: 50,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 70,
                      height: 50,
                      color: Colors.grey.shade300,
                      child: const Icon(Icons.directions_car, size: 24),
                    );
                  },
                ),
              ),
              const SizedBox(height: 6),
              Text(
                widget.transfer['price'] ?? '\$50',
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                  color: Colors.black,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildPriceLine(String label, String amount) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(fontSize: 13, color: Colors.black87),
        ),
        Text(
          amount,
          style: const TextStyle(fontSize: 13, color: Colors.black87),
        ),
      ],
    );
  }

  Widget _buildSelectTravelerDropdown() {
    return Container(
      height: 45,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Row(
        children: [
          const Expanded(
            child: Text(
              'Select Traveler',
              style: TextStyle(fontSize: 13, color: Colors.grey),
            ),
          ),
          Icon(Icons.check_circle_outline, color: Colors.grey.shade400, size: 20),
        ],
      ),
    );
  }

  Widget _buildDropdownWithCheckRight(String hint) {
    return Container(
      height: 45,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Row(
        children: [
          Expanded(
            child: Text(
              hint,
              style: const TextStyle(fontSize: 11, color: Colors.grey),
              overflow: TextOverflow.ellipsis,
            ),
          ),
          Icon(Icons.check_circle_outline, color: Colors.grey.shade400, size: 18),
        ],
      ),
    );
  }

  Widget _buildSimpleTextField(String hint, TextEditingController controller, {TextInputType? keyboardType}) {
    return SizedBox(
      height: 45,
      child: Container(
        decoration: const BoxDecoration(
          border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
        ),
        child: TextField(
          controller: controller,
          keyboardType: keyboardType,
          style: const TextStyle(fontSize: 13),
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

  Widget _buildLabeledTextField(String label, String hint, TextEditingController controller, {TextInputType? keyboardType}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(fontSize: 11, color: Colors.grey),
        ),
        const SizedBox(height: 4),
        _buildSimpleTextField(hint, controller, keyboardType: keyboardType),
      ],
    );
  }

  Widget _buildSmallDropdown(String value) {
    return Container(
      height: 45,
      padding: const EdgeInsets.symmetric(horizontal: 8),
      decoration: const BoxDecoration(
        border: Border(left: BorderSide(color: Color(0xFFD32F2F), width: 4)),
      ),
      child: Row(
        children: [
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 12, color: Colors.grey),
            ),
          ),
          const Icon(Icons.arrow_drop_down, color: Color(0xFF1e5a8e), size: 18),
        ],
      ),
    );
  }

  Widget _buildLabeledDropdownWithCheck(String label) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(fontSize: 11, color: Colors.grey),
        ),
        const SizedBox(height: 4),
        Container(
          height: 45,
          padding: const EdgeInsets.symmetric(horizontal: 12),
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey.shade300),
            borderRadius: BorderRadius.circular(6),
          ),
          child: Row(
            children: [
              const Spacer(),
              Icon(Icons.check_circle_outline, color: Colors.grey.shade400, size: 18),
            ],
          ),
        ),
      ],
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
          borderRadius: BorderRadius.circular(25),
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
                      Navigator.of(context).pushAndRemoveUntil(
                        MaterialPageRoute(builder: (context) => HomePage()),
                        (route) => false,
                      );
                      Navigator.push(
                        context,
                        MaterialPageRoute(builder: (context) => BookingsPage()),
                      );
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
                Image.asset(
                  'assets/images/car.png',
                  width: 120,
                  height: 100,
                  fit: BoxFit.contain,
                  errorBuilder: (context, error, stackTrace) {
                    return const Text(
                      '🚙',
                      style: TextStyle(fontSize: 80),
                    );
                  },
                ),
                const SizedBox(height: 20),
                
                // Success message
                const Text(
                  'Your booking has been completed successfully.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Colors.black87,
                  ),
                ),
                const Text(
                  'Your booking reference is:',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 16),
                
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
                    color: Colors.black54,
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
                      Navigator.of(context).pushAndRemoveUntil(
                        MaterialPageRoute(builder: (context) => HomePage()),
                        (route) => false,
                      );
                      Navigator.push(
                        context,
                        MaterialPageRoute(builder: (context) => BookingsPage()),
                      );
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
    _firstNameController.dispose();
    _lastNameController.dispose();
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
