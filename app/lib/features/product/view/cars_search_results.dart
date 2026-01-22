import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'car_details.dart';
import '../../../core/widgets/unified_ui_components.dart';

class CarsSearchResults extends StatefulWidget {
  final String pickupLocation;
  final String dropoffLocation;
  final DateTime? pickupDate;
  final DateTime? dropoffDate;
  final String pickupTime;
  final String dropoffTime;

  const CarsSearchResults({
    super.key,
    required this.pickupLocation,
    required this.dropoffLocation,
    this.pickupDate,
    this.dropoffDate,
    required this.pickupTime,
    required this.dropoffTime,
  });

  @override
  State<CarsSearchResults> createState() => _CarsSearchResultsState();
}

class _CarsSearchResultsState extends State<CarsSearchResults> {
  // Filter state
  String _sortBy = 'Cheapest';
  String _selectedCarType = '';
  String _selectedRentalCompany = '';
  Map<String, bool> _features = {
    'Automatic Transmission': false,
    'Manual Transmission': false,
    'Air Conditioning': false,
    'Unlimited Mileage': false,
    'Free Cancellation': false,
    'Child Seat Available': false,
    'GPS Included': false,
    'Additional Driver Included': false,
  };
  int _travelTimeHours = 30;
  double _pickupTimeValue = 1.0; // 0=Early Morning, 1=Morning, 2=Afternoon, 3=Evening
  double _dropoffTimeValue = 1.0;

  // Mock car results
  final List<Map<String, dynamic>> _carResults = [
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
    {
      'category': 'Midsize SUV',
      'name': 'Toyota RAV 4 or similar',
      'passengers': '5 Passengers',
      'transmission': 'Automatic',
      'mileage': 'Unlimited mileage',
      'shuttle': 'Shuttle to counter and car',
      'price': '\$3,200',
    },
  ];

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
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${widget.pickupLocation} (1 traveler)',
              style: const TextStyle(
                color: Color(0xFF1e5a8e),
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            Text(
              '${widget.pickupTime} - ${widget.dropoffTime}',
              style: const TextStyle(
                color: Colors.grey,
                fontSize: 12,
                fontWeight: FontWeight.w400,
              ),
            ),
          ],
        ),
        actions: [
          UnifiedFilterButton(
            onPressed: _showFilterBottomSheet,
          ),
        ],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _carResults.length + 1, // +1 for the "Results" header
        itemBuilder: (context, index) {
          if (index == 0) {
            // Results header
            return const Padding(
              padding: EdgeInsets.only(bottom: 16),
              child: Text(
                'Results',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w700,
                  color: Colors.black,
                ),
              ),
            );
          }

          final carIndex = index - 1;

          // Add ad banner after 2nd car
          if (carIndex == 2) {
            return Column(
              children: [
                _buildCarCard(_carResults[carIndex]),
                const SizedBox(height: 16),
                // Ad Banner
                ClipRRect(
                  borderRadius: BorderRadius.circular(12),
                  child: Image.asset(
                    'assets/images/Ad.png',
                    width: double.infinity,
                    height: 100,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return Container(
                        height: 100,
                        decoration: BoxDecoration(
                          color: Colors.grey.shade200,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Center(
                          child: Text('Advertisement'),
                        ),
                      );
                    },
                  ),
                ),
                const SizedBox(height: 16),
              ],
            );
          }

          return Column(
            children: [
              _buildCarCard(_carResults[carIndex]),
              const SizedBox(height: 16),
            ],
          );
        },
      ),
    );
  }

  Widget _buildCarCard(Map<String, dynamic> car) {
    return GestureDetector(
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => CarDetails(
              car: car,
              pickupLocation: widget.pickupLocation,
              dropoffLocation: widget.dropoffLocation,
              pickupDate: widget.pickupDate,
              dropoffDate: widget.dropoffDate,
              pickupTime: widget.pickupTime,
              dropoffTime: widget.dropoffTime,
            ),
          ),
        );
      },
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: const Color(0xFFE8F0F7),
          borderRadius: BorderRadius.circular(8),
        ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  car['category'],
                  style: const TextStyle(
                    color: Color(0xFFD32F2F),
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  car['name'],
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '${car['passengers']} - ${car['transmission']}',
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  car['mileage'],
                  style: const TextStyle(
                    color: Colors.black87,
                    fontSize: 13,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  car['shuttle'],
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
                  width: 80,
                  height: 60,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 80,
                      height: 60,
                      color: Colors.grey.shade300,
                      child: const Icon(Icons.directions_car, size: 30),
                    );
                  },
                ),
              ),
              const SizedBox(height: 8),
              Text(
                car['price'],
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
    ),
    );
  }

  void _showFilterBottomSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => StatefulBuilder(
        builder: (context, setModalState) {
          return Container(
            height: MediaQuery.of(context).size.height * 0.9,
            decoration: const BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
            ),
            child: Column(
              children: [
                // Header
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      IconButton(
                        icon: const Icon(Icons.arrow_back_ios, size: 20),
                        onPressed: () => Navigator.pop(context),
                        padding: EdgeInsets.zero,
                        constraints: const BoxConstraints(),
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        'Filter By:',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),

                const Divider(height: 1),

                // Filters content
                Expanded(
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Sort By
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'Sort By',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFFD32F2F),
                            ),
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                            decoration: BoxDecoration(
                              border: Border.all(color: Colors.grey.shade300),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Row(
                              children: [
                                Text(
                                  _sortBy,
                                  style: const TextStyle(fontSize: 14),
                                ),
                                const SizedBox(width: 8),
                                const Icon(Icons.arrow_drop_down, size: 20),
                              ],
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Car Type
                      const Text(
                        'Car Type',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Economy', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('SUV', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Van', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Luxury', _selectedCarType, (val) {
                              setModalState(() {
                                _selectedCarType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Rental Companies
                      const Text(
                        'Rental Companies',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Hertz', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Sixt', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Avis', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Budget', _selectedRentalCompany, (val) {
                              setModalState(() {
                                _selectedRentalCompany = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Features & Inclusions
                      const Text(
                        'Features & Inclusions',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Automatic Transmission', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Manual Transmission', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Air Conditioning', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Unlimited Mileage', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Free Cancellation', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Child Seat Available', setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('GPS Included', setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Additional Driver Included', setModalState),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Travel time
                      const Text(
                        'Travel time',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          IconButton(
                            onPressed: () {
                              setModalState(() {
                                if (_travelTimeHours > 1) _travelTimeHours--;
                              });
                            },
                            icon: const Icon(Icons.remove_circle_outline),
                            color: const Color(0xFF1e5a8e),
                          ),
                          const SizedBox(width: 16),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                            decoration: BoxDecoration(
                              color: const Color(0xFF1e5a8e),
                              borderRadius: BorderRadius.circular(20),
                            ),
                            child: Text(
                              '${_travelTimeHours}h',
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                          const SizedBox(width: 16),
                          IconButton(
                            onPressed: () {
                              setModalState(() {
                                _travelTimeHours++;
                              });
                            },
                            icon: const Icon(Icons.add_circle_outline),
                            color: const Color(0xFF1e5a8e),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Pickup Time
                      const Text(
                        'Pickup Time',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      _buildTimeSlider(_pickupTimeValue, (value) {
                        setModalState(() {
                          _pickupTimeValue = value;
                        });
                      }),
                      const SizedBox(height: 24),

                      // Drop-off Time
                      const Text(
                        'Drop-off Time',
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFFD32F2F),
                        ),
                      ),
                      const SizedBox(height: 12),
                      _buildTimeSlider(_dropoffTimeValue, (value) {
                        setModalState(() {
                          _dropoffTimeValue = value;
                        });
                      }),
                      const SizedBox(height: 24),
                    ],
                  ),
                ),

                // Done button
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
                        Navigator.pop(context);
                        setState(() {});
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
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildRadioOption(String label, String groupValue, Function(String?) onChanged) {
    return Row(
      children: [
        Radio<String>(
          value: label,
          groupValue: groupValue,
          onChanged: onChanged,
          activeColor: const Color(0xFF1e5a8e),
        ),
        Expanded(
          child: Text(
            label,
            style: const TextStyle(fontSize: 14),
          ),
        ),
      ],
    );
  }

  Widget _buildCheckboxOption(String label, StateSetter setModalState) {
    return Row(
      children: [
        Checkbox(
          value: _features[label] ?? false,
          onChanged: (value) {
            setModalState(() {
              _features[label] = value ?? false;
            });
          },
          activeColor: const Color(0xFF1e5a8e),
        ),
        Expanded(
          child: Text(
            label,
            style: const TextStyle(fontSize: 12),
          ),
        ),
      ],
    );
  }

  Widget _buildTimeSlider(double value, Function(double) onChanged) {
    final labels = ['Early Morning', 'Morning', 'Afternoon', 'Evening'];
    
    return Column(
      children: [
        SliderTheme(
          data: SliderThemeData(
            activeTrackColor: const Color(0xFF1e5a8e),
            inactiveTrackColor: Colors.grey.shade300,
            thumbColor: const Color(0xFF1e5a8e),
            overlayColor: const Color(0xFF1e5a8e).withOpacity(0.2),
            thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 8),
            trackHeight: 3,
          ),
          child: Slider(
            value: value,
            min: 0,
            max: 3,
            divisions: 3,
            onChanged: onChanged,
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 8),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: labels.map((label) {
              return Flexible(
                child: Text(
                  label,
                  style: const TextStyle(fontSize: 10),
                  textAlign: TextAlign.center,
                ),
              );
            }).toList(),
          ),
        ),
      ],
    );
  }
}

