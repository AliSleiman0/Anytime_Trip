import 'package:flutter/material.dart';
import 'car_checkout.dart';

class CarDetails extends StatefulWidget {
  final Map<String, dynamic> car;
  final String pickupLocation;
  final String dropoffLocation;
  final DateTime? pickupDate;
  final DateTime? dropoffDate;
  final String pickupTime;
  final String dropoffTime;

  const CarDetails({
    super.key,
    required this.car,
    required this.pickupLocation,
    required this.dropoffLocation,
    this.pickupDate,
    this.dropoffDate,
    required this.pickupTime,
    required this.dropoffTime,
  });

  @override
  State<CarDetails> createState() => _CarDetailsState();
}

class _CarDetailsState extends State<CarDetails> {
  String _selectedExtra = 'No Extras';

  final List<Map<String, String>> _extras = [
    {'name': 'Booster Seat', 'price': '\$65'},
    {'name': 'Toddler Seat', 'price': '\$65'},
    {'name': 'Infant Seat', 'price': '\$65'},
    {'name': 'Satellite Radio', 'price': '\$65'},
    {'name': 'No Extras', 'price': ''},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Colors.black, size: 20),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          widget.car['category'] ?? 'Midsize SUV',
          style: const TextStyle(
            color: Colors.black,
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Car image and price section
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              widget.car['category'] ?? 'Midsize SUV',
                              style: const TextStyle(
                                color: Color(0xFFD32F2F),
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              widget.car['name'] ?? 'Toyota RAV 4 or similar',
                              style: const TextStyle(
                                color: Colors.black87,
                                fontSize: 13,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              widget.car['passengers'] ?? '5 Passengers - Automatic',
                              style: const TextStyle(
                                color: Colors.black87,
                                fontSize: 13,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              widget.car['mileage'] ?? 'Unlimited mileage',
                              style: const TextStyle(
                                color: Colors.black87,
                                fontSize: 13,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              widget.car['shuttle'] ?? 'Shuttle to counter and car',
                              style: const TextStyle(
                                color: Colors.black87,
                                fontSize: 13,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      Column(
                        children: [
                          ClipRRect(
                            borderRadius: BorderRadius.circular(8),
                            child: Image.asset(
                              'assets/images/urus.jpeg',
                              width: 100,
                              height: 70,
                              fit: BoxFit.cover,
                              errorBuilder: (context, error, stackTrace) {
                                return Container(
                                  width: 100,
                                  height: 70,
                                  color: Colors.grey.shade300,
                                  child: const Icon(Icons.directions_car, size: 30),
                                );
                              },
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            widget.car['price'] ?? '\$3,200',
                            style: const TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w700,
                              color: Colors.black,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),

                  // Car Rental Location section
                  const Text(
                    'Car Rental Location',
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
                      color: Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Pick-up & Drop-off',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                            color: Colors.black54,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          _getFormattedDateRange(),
                          style: const TextStyle(
                            fontSize: 13,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _getLocationName(widget.pickupLocation),
                          style: const TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.w600,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 4),
                        const Text(
                          'Shuttle to the car is located in the airport',
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.black54,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Extras section
                  const Text(
                    'Extras',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w700,
                      color: Colors.black87,
                    ),
                  ),
                  const SizedBox(height: 12),
                  ..._extras.map((extra) {
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: Row(
                        children: [
                          Radio<String>(
                            value: extra['name']!,
                            groupValue: _selectedExtra,
                            onChanged: (value) {
                              setState(() {
                                _selectedExtra = value!;
                              });
                            },
                            activeColor: const Color(0xFF1e5a8e),
                          ),
                          Expanded(
                            child: Text(
                              extra['name']!,
                              style: const TextStyle(
                                fontSize: 14,
                                color: Colors.black87,
                              ),
                            ),
                          ),
                          if (extra['price']!.isNotEmpty)
                            Text(
                              extra['price']!,
                              style: const TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Color(0xFF1e5a8e),
                              ),
                            ),
                        ],
                      ),
                    );
                  }).toList(),
                ],
              ),
            ),
          ),

          // Reserve button
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
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => CarCheckout(
                        car: widget.car,
                        pickupLocation: widget.pickupLocation,
                        dropoffLocation: widget.dropoffLocation,
                        pickupDate: widget.pickupDate,
                        dropoffDate: widget.dropoffDate,
                        pickupTime: widget.pickupTime,
                        dropoffTime: widget.dropoffTime,
                        selectedExtra: _selectedExtra,
                      ),
                    ),
                  );
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF1e5a8e),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
                child: const Text(
                  'Reserve',
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

  String _getFormattedDateRange() {
    if (widget.pickupDate == null || widget.dropoffDate == null) {
      return 'Sun, Oct 19, 1:38 am - Sat, Oct 19, 4:30 pm';
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
    // Extract airport code from location like "Bey-Lebanon" -> "BEY Airport"
    if (location.contains('-')) {
      final code = location.split('-')[0].toUpperCase();
      return '$code Airport';
    }
    return '$location Airport';
  }
}
