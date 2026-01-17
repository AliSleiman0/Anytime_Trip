import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'hotel_payment.dart';

class HotelDetails extends StatefulWidget {
  final Map<String, dynamic> hotel;
  final String checkinDate;
  final String checkoutDate;
  final int adults;
  final int children;
  final int rooms;

  const HotelDetails({
    super.key,
    required this.hotel,
    required this.checkinDate,
    required this.checkoutDate,
    required this.adults,
    required this.children,
    required this.rooms,
  });

  @override
  State<HotelDetails> createState() => _HotelDetailsState();
}

class _HotelDetailsState extends State<HotelDetails> {
  late PageController _carouselController;
  int _currentCarouselIndex = 0;
  
  final List<String> _carouselImages = [
    'assets/images/hotel1.jpg',
    'assets/images/hotel2.jpg',
    'assets/images/hotel3.jpg',
  ];

  final List<Map<String, dynamic>> _rooms = [
    {
      'name': 'Traditional Room, 2 beds',
      'size': '330 sq ft',
      'sleeps': '3',
      'features': ['Free WiFi', 'Ocean View'],
      'price': 603,
      'image': 'assets/images/room.jpeg',
    },
    {
      'name': 'Deluxe Room, 1 King Bed',
      'size': '400 sq ft',
      'sleeps': '2',
      'features': ['Free WiFi', 'Ocean view'],
      'price': 785,
      'image': 'assets/images/room.jpeg',
    },
    {
      'name': 'Suite, 2 Queen Beds',
      'size': '500 sq ft',
      'sleeps': '5',
      'features': ['Free WiFi', 'Mini fridge'],
      'price': 950,
      'image': 'assets/images/room.jpeg',
    },
  ];

  int _adultsCount = 1;
  int _childrenCount = 2;
  int _roomsCount = 1;
  DateTime? _checkinDate;
  DateTime? _checkoutDate;

  @override
  void initState() {
    super.initState();
    _carouselController = PageController();
    _adultsCount = widget.adults;
    _childrenCount = widget.children;
    _roomsCount = widget.rooms;
  }

  @override
  void dispose() {
    _carouselController.dispose();
    super.dispose();
  }

  String get _travelersSummary {
    List<String> parts = [];
    if (_adultsCount > 0) parts.add('$_adultsCount Adult${_adultsCount > 1 ? 's' : ''}');
    if (_childrenCount > 0) parts.add('$_childrenCount child${_childrenCount > 1 ? 'ren' : ''}');
    parts.add('$_roomsCount room${_roomsCount > 1 ? 's' : ''}');
    return parts.join(', ');
  }

  String get _formattedCheckinDate {
    if (_checkinDate == null) return widget.checkinDate;
    return '${_checkinDate!.day.toString().padLeft(2, '0')} Nov ${_checkinDate!.year.toString().substring(2)}';
  }

  String get _formattedCheckoutDate {
    if (_checkoutDate == null) return widget.checkoutDate;
    return '${_checkoutDate!.day.toString().padLeft(2, '0')} Nov ${_checkoutDate!.year.toString().substring(2)}';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: Column(
        children: [
          // Custom App Bar with Carousel
          SizedBox(
            height: 280,
            child: Stack(
              children: [
                // Carousel
                PageView.builder(
                  controller: _carouselController,
                  onPageChanged: (index) {
                    setState(() {
                      _currentCarouselIndex = index;
                    });
                  },
                  itemCount: _carouselImages.length,
                  itemBuilder: (context, index) {
                    return Image.asset(
                      _carouselImages[index],
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) {
                        return Container(
                          color: Colors.grey[300],
                          child: const Icon(Icons.hotel, size: 80, color: Colors.grey),
                        );
                      },
                    );
                  },
                ),
                
                // Search bar at top
                Positioned(
                  top: 40,
                  left: 16,
                  right: 16,
                  child: Container(
                    height: 50,
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(25),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.1),
                          blurRadius: 8,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: Row(
                      children: [
                        // Back button
                        GestureDetector(
                          onTap: () => Navigator.pop(context),
                          child: Container(
                            padding: const EdgeInsets.all(12),
                            child: const Icon(Icons.arrow_back_ios_new, color: Colors.black, size: 18),
                          ),
                        ),
                        const SizedBox(width: 8),
                        // Hotel name
                        Expanded(
                          child: Text(
                            widget.hotel['name'],
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                
                // Carousel indicators
                Positioned(
                  bottom: 16,
                  left: 0,
                  right: 0,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(_carouselImages.length, (index) {
                      return Container(
                        margin: const EdgeInsets.symmetric(horizontal: 4),
                        width: 8,
                        height: 8,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: _currentCarouselIndex == index
                              ? Colors.white
                              : Colors.white.withOpacity(0.5),
                        ),
                      );
                    }),
                  ),
                ),
              ],
            ),
          ),

          // Scrollable content
          Expanded(
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: 16),

                  const SizedBox(height: 16),

                  // Property details section
                  Container(
                    margin: const EdgeInsets.symmetric(horizontal: 16),
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xFFE3F2FD),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            const Text(
                              'Property details',
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: Colors.black87,
                              ),
                            ),
                            const Spacer(),
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                              decoration: BoxDecoration(
                                color: const Color(0xFF1e5a8e),
                                borderRadius: BorderRadius.circular(4),
                              ),
                              child: const Text(
                                '330 × 417',
                                style: TextStyle(
                                  fontSize: 12,
                                  fontWeight: FontWeight.w600,
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        const Text('Indoor pool', style: TextStyle(fontSize: 12, color: Colors.black87)),
                        const Text('American Restaurant available', style: TextStyle(fontSize: 12, color: Colors.black87)),
                        const Text('American Restaurant', style: TextStyle(fontSize: 12, color: Colors.black87)),
                        const Text('Vegetarian breakfast available', style: TextStyle(fontSize: 12, color: Colors.black87)),
                        const Text('dogs and cats allowed', style: TextStyle(fontSize: 12, color: Colors.black87)),
                        const Text('Free Wifi', style: TextStyle(fontSize: 12, color: Colors.black87)),
                      ],
                    ),
                  ),

                  const SizedBox(height: 24),

                  // Explore Area
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 16),
                    child: Text(
                      'Explore Area',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                        color: Colors.black,
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),

                  // Map
                  Container(
                    margin: const EdgeInsets.symmetric(horizontal: 16),
                    height: 150,
                    decoration: BoxDecoration(
                      color: Colors.grey[200],
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(8),
                      child: Image.asset(
                        'assets/images/location.png',
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return Stack(
                            children: [
                              Center(
                                child: Icon(Icons.map, size: 50, color: Colors.grey[400]),
                              ),
                              Positioned(
                                top: 8,
                                right: 8,
                                child: Container(
                                  padding: const EdgeInsets.all(4),
                                  decoration: BoxDecoration(
                                    color: Colors.white,
                                    borderRadius: BorderRadius.circular(4),
                                    boxShadow: [
                                      BoxShadow(
                                        color: Colors.black.withOpacity(0.1),
                                        blurRadius: 4,
                                      ),
                                    ],
                                  ),
                                  child: const Icon(Icons.my_location, size: 16, color: Colors.pink),
                                ),
                              ),
                            ],
                          );
                        },
                      ),
                    ),
                  ),

                  const SizedBox(height: 12),

                  // Locations list
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: Column(
                      children: [
                        _buildLocationItem('Zeitouna Bay', '20 min'),
                        _buildLocationItem('Zeitouna Bay', '20 min'),
                        _buildLocationItem('Zeitouna Bay', '20 min'),
                      ],
                    ),
                  ),

                  const SizedBox(height: 24),

                  // Select Room header
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                    color: const Color(0xFFF5F5F5),
                    child: Row(
                      children: [
                        const Text(
                          'Select Room',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w700,
                            color: Colors.black,
                          ),
                        ),
                        const Spacer(),
                        GestureDetector(
                          onTap: () => _showDatePicker(),
                          child: Row(
                            children: [
                              Text(
                                '$_formattedCheckinDate - $_formattedCheckoutDate',
                                style: const TextStyle(
                                  fontSize: 12,
                                  color: Colors.black87,
                                ),
                              ),
                              const SizedBox(width: 4),
                              const Icon(Icons.calendar_today, size: 14, color: Color(0xFFD32F2F)),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 16),

                  // Number of Travelers
                  GestureDetector(
                    onTap: () => _showTravelersBottomSheet(),
                    child: Container(
                      margin: const EdgeInsets.symmetric(horizontal: 16),
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        border: Border.all(color: const Color(0xFFD32F2F)),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Row(
                        children: [
                          const Icon(Icons.people, color: Color(0xFFD32F2F), size: 20),
                          const SizedBox(width: 8),
                          const Text(
                            'Number of Travelers',
                            style: TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFFD32F2F),
                            ),
                          ),
                          const Spacer(),
                          Text(
                            _travelersSummary,
                            style: const TextStyle(
                              fontSize: 12,
                              color: Colors.black87,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),

                  const SizedBox(height: 16),

                  // Rooms list
                  ..._rooms.map((room) => _buildRoomCard(room)).toList(),

                  const SizedBox(height: 24),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLocationItem(String name, String distance) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Text(
            name,
            style: const TextStyle(
              fontSize: 13,
              color: Colors.black87,
            ),
          ),
          const Spacer(),
          Text(
            distance,
            style: TextStyle(
              fontSize: 13,
              color: Colors.blue[700],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRoomCard(Map<String, dynamic> room) {
    return GestureDetector(
      onTap: () => _showPaymentOptionsDialog(room),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: const Color(0xFFE8F0F7),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Room details
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    room['name'],
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFFD32F2F),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    room['size'],
                    style: const TextStyle(fontSize: 13, color: Colors.black87),
                  ),
                  Text(
                    'Sleeps ${room['sleeps']}',
                    style: const TextStyle(fontSize: 13, color: Colors.black87),
                  ),
                  const SizedBox(height: 4),
                  ...room['features'].map<Widget>((feature) => Text(
                    feature,
                    style: const TextStyle(fontSize: 13, color: Colors.black87),
                  )).toList(),
                ],
              ),
            ),
            const SizedBox(width: 12),
            // Room image and price
            Column(
              children: [
                ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: Image.asset(
                    room['image'],
                    width: 80,
                    height: 60,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return Container(
                        width: 80,
                        height: 60,
                        color: Colors.grey[300],
                        child: const Icon(Icons.bed, size: 30, color: Colors.grey),
                      );
                    },
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  '\$${room['price']}',
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

  void _showDatePicker() async {
    final DateTimeRange? picked = await showDateRangePicker(
      context: context,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365)),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: const ColorScheme.light(
              primary: Color(0xFF1e5a8e),
              onPrimary: Colors.white,
              onSurface: Colors.black,
            ),
          ),
          child: child!,
        );
      },
    );

    if (picked != null) {
      setState(() {
        _checkinDate = picked.start;
        _checkoutDate = picked.end;
      });
    }
  }

  void _showTravelersBottomSheet() {
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return Container(
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(20),
                  topRight: Radius.circular(20),
                ),
              ),
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Handle
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey[300],
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Title
                  const Text(
                    'Number of Travelers',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Adults
                  _buildCounterRow('Adults', _adultsCount, (value) {
                    setModalState(() {
                      _adultsCount = value;
                    });
                  }),
                  const SizedBox(height: 16),

                  // Children
                  _buildCounterRow('Children', _childrenCount, (value) {
                    setModalState(() {
                      _childrenCount = value;
                    });
                  }),
                  const SizedBox(height: 16),

                  // Rooms
                  _buildCounterRow('Rooms', _roomsCount, (value) {
                    setModalState(() {
                      _roomsCount = value;
                    });
                  }),
                  const SizedBox(height: 32),

                  // Done button
                  SizedBox(
                    width: double.infinity,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: () {
                        setState(() {});
                        Navigator.pop(context);
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
                  const SizedBox(height: 16),
                ],
              ),
            );
          },
        );
      },
    );
  }

  Widget _buildCounterRow(String label, int value, Function(int) onChanged) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
        Row(
          children: [
            IconButton(
              onPressed: () {
                if (value > 0) onChanged(value - 1);
              },
              icon: const Icon(Icons.remove_circle_outline),
              color: const Color(0xFF1e5a8e),
            ),
            Container(
              width: 40,
              alignment: Alignment.center,
              child: Text(
                value.toString(),
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            IconButton(
              onPressed: () {
                onChanged(value + 1);
              },
              icon: const Icon(Icons.add_circle_outline),
              color: const Color(0xFF1e5a8e),
            ),
          ],
        ),
      ],
    );
  }

  void _showPaymentOptionsDialog(Map<String, dynamic> room) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          child: Container(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Close button
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    GestureDetector(
                      onTap: () => Navigator.pop(context),
                      child: Container(
                        decoration: const BoxDecoration(
                          color: Color(0xFFD32F2F),
                          shape: BoxShape.circle,
                        ),
                        padding: const EdgeInsets.all(6),
                        child: const Icon(
                          Icons.close,
                          color: Colors.white,
                          size: 18,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Pay the total now section
                const Text(
                  'Pay the total now',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                    color: Color(0xFF1e5a8e),
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'You can use a valid Any time Travel coupon',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 16),

                // Pay now button
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pop(context);
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => HotelPayment(
                            hotel: widget.hotel,
                            room: room,
                            checkinDate: widget.checkinDate,
                            checkoutDate: widget.checkoutDate,
                            adults: _adultsCount,
                            children: _childrenCount,
                            rooms: _roomsCount,
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
                      'Pay now',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 24),

                // Pay when you stay section
                const Text(
                  'Pay when you stay',
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                    color: Color(0xFF1e5a8e),
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Pay the property directly in their preferred currency',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 16),

                // Pay at department button
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pop(context);
                      _showBookingConfirmationDialog();
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF1e5a8e),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(25),
                      ),
                    ),
                    child: const Text(
                      'Pay at department',
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

  void _showBookingConfirmationDialog() {
    // Generate random booking reference
    final bookingReference = '1872306517801';
    
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          child: Container(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Close button
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    GestureDetector(
                      onTap: () => Navigator.pop(context),
                      child: Container(
                        decoration: const BoxDecoration(
                          color: Color(0xFFD32F2F),
                          shape: BoxShape.circle,
                        ),
                        padding: const EdgeInsets.all(6),
                        child: const Icon(
                          Icons.close,
                          color: Colors.white,
                          size: 18,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Building/Hotel icon
                Container(
                  width: 100,
                  height: 100,
                  decoration: BoxDecoration(
                    color: const Color(0xFF1e5a8e).withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Icon(
                    Icons.apartment,
                    size: 60,
                    color: Color(0xFF1e5a8e),
                  ),
                ),
                const SizedBox(height: 24),

                // Success message
                const Text(
                  'Your booking has been completed successfully.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  'Your booking reference is:',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 12),

                // Booking reference number
                Text(
                  bookingReference,
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w700,
                    color: Colors.black,
                    letterSpacing: 1.2,
                  ),
                ),
                const SizedBox(height: 16),

                // Email confirmation message
                const Text(
                  'You will receive confirmation email including your itinerary details',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 13,
                    color: Colors.black87,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 24),

                // Done button
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pop(context);
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
}
