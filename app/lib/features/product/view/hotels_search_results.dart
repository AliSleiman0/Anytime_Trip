import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'hotel_details.dart';
import '../../../core/widgets/unified_ui_components.dart';

class HotelsSearchResults extends StatefulWidget {
  final String destination;
  final String checkinDate;
  final String checkoutDate;
  final int adults;
  final int children;
  final int rooms;
  final String nationality;

  const HotelsSearchResults({
    super.key,
    required this.destination,
    required this.checkinDate,
    required this.checkoutDate,
    required this.adults,
    required this.children,
    required this.rooms,
    required this.nationality,
  });

  @override
  State<HotelsSearchResults> createState() => _HotelsSearchResultsState();
}

class _HotelsSearchResultsState extends State<HotelsSearchResults> {
  // Filter state
  String _sortBy = 'Cheapest';
  String _selectedVehicleType = '';
  Map<String, bool> _propertyFeatures = {
    'Indoor Pool': false,
    'Valet Parking Available': false,
    'American Restaurant': false,
    'Vegetarian Breakfast': false,
    'Dogs & Cats Allowed': false,
    'Free Wi-Fi': false,
  };
  Map<String, bool> _roomFeatures = {
    'Ocean View': false,
    'Free Breakfast': false,
    'Balcony': false,
    'Mini Fridge': false,
    'King Bed / Queen Bed': false,
    'Family Room': false,
  };

  final List<Map<String, dynamic>> _hotelResults = [
    {
      'name': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'image': 'assets/images/hotel.jpg',
      'price': 603,
      'features': ['Breakfast + Pool Included', 'Fully Refundable'],
      'rating': 4.5,
    },
    {
      'name': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'image': 'assets/images/hotel.jpg',
      'price': 603,
      'features': ['Breakfast + Pool Included', 'Fully Refundable'],
      'rating': 4.5,
    },
    {
      'name': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'image': 'assets/images/hotel.jpg',
      'price': 603,
      'features': ['Breakfast + Pool Included', 'Fully Refundable'],
      'rating': 4.5,
    },
    {
      'name': 'Beirut Hotel',
      'location': 'Beirut, Lebanon',
      'image': 'assets/images/hotel.jpg',
      'price': 603,
      'features': ['Breakfast + Pool Included', 'Fully Refundable'],
      'rating': 4.5,
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
              '${widget.destination} (${widget.adults + widget.children} traveler${widget.adults + widget.children > 1 ? 's' : ''})',
              style: const TextStyle(
                color: Color(0xFF1e5a8e),
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            Text(
              '${widget.checkinDate} - ${widget.checkoutDate}',
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
        itemCount: _hotelResults.length + 1, // +1 for the "Results" header
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

          final hotelIndex = index - 1;

          // Add ad banner after 2nd hotel
          if (hotelIndex == 2) {
            return Column(
              children: [
                _buildHotelCard(_hotelResults[hotelIndex]),
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
              _buildHotelCard(_hotelResults[hotelIndex]),
              const SizedBox(height: 16),
            ],
          );
        },
      ),
    );
  }

  Widget _buildHotelCard(Map<String, dynamic> hotel) {
    return GestureDetector(
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => HotelDetails(
              hotel: hotel,
              checkinDate: widget.checkinDate,
              checkoutDate: widget.checkoutDate,
              adults: widget.adults,
              children: widget.children,
              rooms: widget.rooms,
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
        child: IntrinsicHeight(
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          hotel['name'],
                          style: const TextStyle(
                            color: Color(0xFFD32F2F),
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          hotel['location'],
                          style: const TextStyle(
                            color: Colors.black87,
                            fontSize: 13,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          hotel['features'][0],
                          style: const TextStyle(
                            color: Colors.black87,
                            fontSize: 13,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          hotel['features'][1],
                          style: const TextStyle(
                            color: Colors.black87,
                            fontSize: 13,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text(
                      '\$${hotel['price']} total',
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        color: Colors.black,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 12),
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Image.asset(
                  hotel['image'],
                  width: 80,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 80,
                      color: Colors.grey.shade300,
                      child: const Icon(Icons.hotel, size: 30),
                    );
                  },
                ),
              ),
            ],
          ),
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

                      // Vehicle Type (Hotel Type)
                      const Text(
                        'Vehicle Type',
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
                            child: _buildRadioOption('Economy', _selectedVehicleType, (val) {
                              setModalState(() {
                                _selectedVehicleType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('SUV', _selectedVehicleType, (val) {
                              setModalState(() {
                                _selectedVehicleType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildRadioOption('Van', _selectedVehicleType, (val) {
                              setModalState(() {
                                _selectedVehicleType = val ?? '';
                              });
                            }),
                          ),
                          Expanded(
                            child: _buildRadioOption('Luxury', _selectedVehicleType, (val) {
                              setModalState(() {
                                _selectedVehicleType = val ?? '';
                              });
                            }),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Property Features
                      const Text(
                        'Property Features',
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
                            child: _buildCheckboxOption('Indoor Pool', _propertyFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Valet Parking Available', _propertyFeatures, setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('American Restaurant', _propertyFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Vegetarian Breakfast', _propertyFeatures, setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Dogs & Cats Allowed', _propertyFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Free Wi-Fi', _propertyFeatures, setModalState),
                          ),
                        ],
                      ),
                      const SizedBox(height: 24),

                      // Room Features
                      const Text(
                        'Room Features',
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
                            child: _buildCheckboxOption('Ocean View', _roomFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Free Breakfast', _roomFeatures, setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('Balcony', _roomFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Mini Fridge', _roomFeatures, setModalState),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Expanded(
                            child: _buildCheckboxOption('King Bed / Queen Bed', _roomFeatures, setModalState),
                          ),
                          Expanded(
                            child: _buildCheckboxOption('Family Room', _roomFeatures, setModalState),
                          ),
                        ],
                      ),
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

  Widget _buildCheckboxOption(String label, Map<String, bool> features, StateSetter setModalState) {
    return Row(
      children: [
        Checkbox(
          value: features[label] ?? false,
          onChanged: (value) {
            setModalState(() {
              features[label] = value ?? false;
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
}
