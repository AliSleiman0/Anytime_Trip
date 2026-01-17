import 'package:flutter/material.dart';
import 'package:get/get.dart';

class FlightSearchResults extends StatefulWidget {
  final String from;
  final String to;
  final String departureDate;
  final String? returnDate;
  final int travelers;
  final String travelClass;

  const FlightSearchResults({
    super.key,
    required this.from,
    required this.to,
    required this.departureDate,
    this.returnDate,
    required this.travelers,
    required this.travelClass,
  });

  @override
  State<FlightSearchResults> createState() => _FlightSearchResultsState();
}

class _FlightSearchResultsState extends State<FlightSearchResults> {
  int _selectedDateIndex = 1;
  
  // Filter state
  Map<String, bool> _filterStops = {};
  Map<String, bool> _filterAirlines = {};
  Map<String, bool> _filterBaggage = {};
  int _travelTime = 30;
  double _departureTime = 0.5;
  double _arrivalTime = 0.5;
  
  final List<Map<String, dynamic>> _datePrices = [
    {'date': 'Thu, Sep 25', 'price': '1,200\$'},
    {'date': 'Thu, Sep 25', 'price': '1,200\$'},
    {'date': 'Thu, Sep 25', 'price': '1,200\$'},
    {'date': 'Thu, Sep 25', 'price': '1,200\$'},
  ];

  final List<Map<String, dynamic>> _bestDeals = [
    {
      'from': 'BEY',
      'to': 'DXB',
      'fromCity': 'Beirut',
      'toCity': 'Dubai',
      'date': '06 Feb 2026',
      'time': '04:10',
      'type': 'Direct Flight',
    },
    {
      'from': 'BEY',
      'to': 'DXB',
      'fromCity': 'Beirut',
      'toCity': 'Dubai',
      'date': '12 Feb 2026',
      'time': '06:10',
      'type': 'Direct Flight',
    },
  ];

  final List<Map<String, dynamic>> _flightResults = [
    {
      'departureTime': '7:55am',
      'arrivalTime': '3:05pm',
      'from': 'BEY',
      'to': 'DXB',
      'airline': 'Middle East Airlines',
      'price': 3200,
      'checkedBag': '1 x 23kg',
      'carryOn': '1 x 23kg',
      'duration': '18h - 10min',
      'stops': 'Direct Flight',
    },
    {
      'departureTime': '7:55am',
      'arrivalTime': '3:05pm',
      'connectTime1': '4:05pm',
      'connectTime2': '8:23pm',
      'from': 'BEY',
      'to': 'DXB',
      'airline': 'Middle East Airlines',
      'price': 3200,
      'checkedBag': '1 x 23kg',
      'carryOn': '1 x 23kg',
      'duration': '18h - 10min',
      'stops': '1 stop',
    },
    {
      'departureTime': '7:55am',
      'arrivalTime': '3:05pm',
      'connectTime1': '4:05pm',
      'connectTime2': '8:23pm',
      'from': 'BEY',
      'to': 'DXB',
      'airline': 'Middle East Airlines',
      'price': 3200,
      'checkedBag': '1 x 23kg',
      'carryOn': '1 x 23kg',
      'duration': '18h - 10min',
      'stops': '1 stop',
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F5F5),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e)),
          onPressed: () => Navigator.pop(context),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${widget.from} (${widget.travelers} traveler${widget.travelers > 1 ? 's' : ''})',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Color(0xFF1e5a8e),
              ),
            ),
            Text(
              widget.departureDate + (widget.returnDate != null ? ' - ${widget.returnDate}' : ''),
              style: const TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w400,
                color: Colors.grey,
              ),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list, color: Color(0xFFD32F2F)),
            onPressed: () {
              _showFilterBottomSheet(context);
            },
          ),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: 16),
                  
                  // Select Departure Flight Section
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: const Text(
                      'Select Departure Flight',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        color: Colors.black,
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  
                  // Date Selection Chips
                  SizedBox(
                    height: 60,
                    child: ListView.builder(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      itemCount: _datePrices.length,
                      itemBuilder: (context, index) {
                        final dateInfo = _datePrices[index];
                        final isSelected = index == _selectedDateIndex;
                        return GestureDetector(
                          onTap: () {
                            setState(() {
                              _selectedDateIndex = index;
                            });
                          },
                          child: Container(
                            margin: const EdgeInsets.only(right: 8),
                            padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 8),
                            decoration: BoxDecoration(
                              color: isSelected ? const Color(0xFF4A7BA7) : Colors.white,
                              borderRadius: BorderRadius.circular(20),
                              border: Border.all(
                                color: isSelected ? const Color(0xFF4A7BA7) : const Color(0xFFE0E0E0),
                                width: 1,
                              ),
                            ),
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Text(
                                  dateInfo['date'],
                                  style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w500,
                                    color: isSelected ? Colors.white : Colors.black87,
                                  ),
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  dateInfo['price'],
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w700,
                                    color: isSelected ? Colors.white : Colors.black87,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Best Deals Section
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: const Text(
                      'Best Deals',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        color: Colors.black,
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  
                  SizedBox(
                    height: 110,
                    child: ListView.builder(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      itemCount: _bestDeals.length,
                      itemBuilder: (context, index) {
                        final deal = _bestDeals[index];
                        return Container(
                          width: 180,
                          margin: const EdgeInsets.only(right: 10),
                          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
                          decoration: BoxDecoration(
                            gradient: const LinearGradient(
                              colors: [Color(0xFF5A9BD5), Color(0xFF2C5F8D)],
                              begin: Alignment.topLeft,
                              end: Alignment.bottomRight,
                            ),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    deal['from'],
                                    style: const TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.w700,
                                      color: Colors.white,
                                    ),
                                  ),
                                  Expanded(
                                    child: Row(
                                      mainAxisAlignment: MainAxisAlignment.center,
                                      children: [
                                        const Expanded(
                                          child: Divider(
                                            color: Colors.white,
                                            thickness: 1,
                                          ),
                                        ),
                                        Padding(
                                          padding: const EdgeInsets.symmetric(horizontal: 4),
                                          child: Image.asset(
                                            'assets/images/plane.png',
                                            width: 14,
                                            height: 14,
                                          ),
                                        ),
                                        const Expanded(
                                          child: Divider(
                                            color: Colors.white,
                                            thickness: 1,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  Text(
                                    deal['to'],
                                    style: const TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.w700,
                                      color: Colors.white,
                                    ),
                                  ),
                                ],
                              ),
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    deal['fromCity'],
                                    style: const TextStyle(
                                      fontSize: 9,
                                      color: Colors.white,
                                    ),
                                  ),
                                  Text(
                                    '${deal['date']}',
                                    style: const TextStyle(
                                      fontSize: 9,
                                      color: Colors.white,
                                    ),
                                  ),
                                  Text(
                                    '${deal['time']}',
                                    style: const TextStyle(
                                      fontSize: 9,
                                      color: Colors.white,
                                    ),
                                  ),
                                ],
                              ),
                              Text(
                                deal['type'],
                                style: const TextStyle(
                                  fontSize: 10,
                                  fontWeight: FontWeight.w600,
                                  color: Colors.white,
                                ),
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Results Section
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: const Text(
                      'Results',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        color: Colors.black,
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  
                  // Flight Results List with Ad Banner
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: Column(
                      children: [
                        // First 2 flights
                        ...List.generate(
                          _flightResults.length > 2 ? 2 : _flightResults.length,
                          (index) => _buildFlightCard(_flightResults[index]),
                        ),
                        
                        // Ad Banner after 2nd flight
                        if (_flightResults.length > 2)
                          Container(
                            height: 100,
                            margin: const EdgeInsets.only(bottom: 16),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(12),
                              image: const DecorationImage(
                                image: AssetImage('assets/images/Ad.png'),
                                fit: BoxFit.cover,
                              ),
                            ),
                          ),
                        
                        // Remaining flights
                        if (_flightResults.length > 2)
                          ...List.generate(
                            _flightResults.length - 2,
                            (index) => _buildFlightCard(_flightResults[index + 2]),
                          ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showFilterBottomSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return Container(
              height: MediaQuery.of(context).size.height * 0.9,
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(20),
                  topRight: Radius.circular(20),
                ),
              ),
              child: Column(
                children: [
                  // Header
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: const BorderRadius.only(
                        topLeft: Radius.circular(20),
                        topRight: Radius.circular(20),
                      ),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 4,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: Row(
                      children: [
                        IconButton(
                          icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e)),
                          onPressed: () => Navigator.pop(context),
                          padding: EdgeInsets.zero,
                          constraints: const BoxConstraints(),
                        ),
                        const SizedBox(width: 12),
                        const Text(
                          'Filter By:',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w600,
                            color: Colors.black,
                          ),
                        ),
                      ],
                    ),
                  ),
                  
                  // Content
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Sort By
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text(
                                'Sort By',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFFD32F2F),
                                ),
                              ),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                                decoration: BoxDecoration(
                                  border: Border.all(color: Colors.grey.shade300),
                                  borderRadius: BorderRadius.circular(6),
                                ),
                                child: Row(
                                  children: const [
                                    Text(
                                      'Cheapest',
                                      style: TextStyle(
                                        fontSize: 13,
                                        color: Colors.black87,
                                      ),
                                    ),
                                    SizedBox(width: 8),
                                    Icon(Icons.keyboard_arrow_down, size: 18),
                                  ],
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 24),
                          
                          // Stops
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: const [
                              Text(
                                'Stops',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFFD32F2F),
                                ),
                              ),
                              Text(
                                'From',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFF1e5a8e),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          _buildFilterCheckbox('1 stop', '\$1,000', setModalState, 'stops'),
                          _buildFilterCheckbox('2 stops', '\$1,000', setModalState, 'stops'),
                          const SizedBox(height: 24),
                          
                          // Airlines
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: const [
                              Text(
                                'Airlines',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFFD32F2F),
                                ),
                              ),
                              Text(
                                'From',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFF1e5a8e),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          _buildFilterCheckbox('Middle East Airlines', '\$1,000', setModalState, 'airlines'),
                          _buildFilterCheckbox('Sun Express', '\$1,000', setModalState, 'airlines'),
                          _buildFilterCheckbox('Etihad Airways', '\$1,000', setModalState, 'airlines'),
                          _buildFilterCheckbox('Air Canada', '\$1,000', setModalState, 'airlines'),
                          const SizedBox(height: 24),
                          
                          // Travel and baggage
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: const [
                              Text(
                                'Travel and baggage',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFFD32F2F),
                                ),
                              ),
                              Text(
                                'From',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFF1e5a8e),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          _buildFilterCheckbox('Seat choice included', '\$1,000', setModalState, 'baggage'),
                          _buildFilterCheckbox('Carry-on bag included', '\$1,000', setModalState, 'baggage'),
                          _buildFilterCheckbox('No cancel fee', '\$1,000', setModalState, 'baggage'),
                          _buildFilterCheckbox('Changes included', '\$1,000', setModalState, 'baggage'),
                          const SizedBox(height: 24),
                          
                          // Travel time
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text(
                                'Travel time',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFFD32F2F),
                                ),
                              ),
                              Row(
                                children: [
                                  Container(
                                    decoration: const BoxDecoration(
                                      color: Color(0xFF1e5a8e),
                                      shape: BoxShape.circle,
                                    ),
                                    child: IconButton(
                                      icon: const Icon(Icons.remove, color: Colors.white, size: 18),
                                      onPressed: () {
                                        setModalState(() {
                                          if (_travelTime > 0) _travelTime--;
                                        });
                                      },
                                      padding: const EdgeInsets.all(4),
                                      constraints: const BoxConstraints(),
                                    ),
                                  ),
                                  const SizedBox(width: 16),
                                  Text(
                                    '${_travelTime}h',
                                    style: const TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.w600,
                                      color: Colors.black87,
                                    ),
                                  ),
                                  const SizedBox(width: 16),
                                  Container(
                                    decoration: const BoxDecoration(
                                      color: Color(0xFF1e5a8e),
                                      shape: BoxShape.circle,
                                    ),
                                    child: IconButton(
                                      icon: const Icon(Icons.add, color: Colors.white, size: 18),
                                      onPressed: () {
                                        setModalState(() {
                                          _travelTime++;
                                        });
                                      },
                                      padding: const EdgeInsets.all(4),
                                      constraints: const BoxConstraints(),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                          const SizedBox(height: 24),
                          
                          // Departure Time
                          const Text(
                            'Departure Time',
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFFD32F2F),
                            ),
                          ),
                          const SizedBox(height: 12),
                          _buildTimeSlider(setModalState, true),
                          const SizedBox(height: 24),
                          
                          // Arrival Time
                          const Text(
                            'Arrival Time',
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFFD32F2F),
                            ),
                          ),
                          const SizedBox(height: 12),
                          _buildTimeSlider(setModalState, false),
                          const SizedBox(height: 80),
                        ],
                      ),
                    ),
                  ),
                  
                  // Done Button
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.1),
                          blurRadius: 8,
                          offset: const Offset(0, -2),
                        ),
                      ],
                    ),
                    child: SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () => Navigator.pop(context),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF2c5f8d),
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(30),
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
        );
      },
    );
  }

  Widget _buildFilterCheckbox(String label, String price, StateSetter setModalState, String category) {
    Map<String, bool> categoryMap;
    switch (category) {
      case 'stops':
        categoryMap = _filterStops;
        break;
      case 'airlines':
        categoryMap = _filterAirlines;
        break;
      case 'baggage':
        categoryMap = _filterBaggage;
        break;
      default:
        categoryMap = {};
    }
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Expanded(
            child: Row(
              children: [
                SizedBox(
                  width: 20,
                  height: 20,
                  child: Checkbox(
                    value: categoryMap[label] ?? false,
                    onChanged: (value) {
                      setModalState(() {
                        categoryMap[label] = value ?? false;
                      });
                    },
                    activeColor: const Color(0xFF1e5a8e),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(4),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    label,
                    style: const TextStyle(
                      fontSize: 13,
                      color: Colors.black87,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Text(
            price,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: Color(0xFF1e5a8e),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTimeSlider(StateSetter setModalState, bool isDeparture) {
    double currentValue = isDeparture ? _departureTime : _arrivalTime;
    
    return Column(
      children: [
        Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.wb_sunny_outlined,
                size: 20,
                color: Colors.grey,
              ),
            ),
            Expanded(
              child: SliderTheme(
                data: SliderThemeData(
                  trackHeight: 3,
                  thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 8),
                  overlayShape: const RoundSliderOverlayShape(overlayRadius: 16),
                  activeTrackColor: const Color(0xFF1e5a8e),
                  inactiveTrackColor: Colors.grey.shade300,
                  thumbColor: const Color(0xFF1e5a8e),
                ),
                child: Slider(
                  value: currentValue,
                  onChanged: (value) {
                    setModalState(() {
                      if (isDeparture) {
                        _departureTime = value;
                      } else {
                        _arrivalTime = value;
                      }
                    });
                  },
                  min: 0,
                  max: 1,
                ),
              ),
            ),
          ],
        ),
        Padding(
          padding: const EdgeInsets.only(left: 40),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: const [
              Text(
                'Early Morning',
                style: TextStyle(fontSize: 11, color: Colors.grey),
              ),
              Text(
                'Morning',
                style: TextStyle(fontSize: 11, color: Colors.grey),
              ),
              Text(
                'Afternoon',
                style: TextStyle(fontSize: 11, color: Colors.grey),
              ),
              Text(
                'Evening',
                style: TextStyle(fontSize: 11, color: Colors.grey),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildFlightCard(Map<String, dynamic> flight) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Time and Price Row
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Time section
                    Expanded(
                      child: Row(
                        children: [
                          Text(
                            flight['departureTime'],
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black87,
                            ),
                          ),
                          const SizedBox(width: 4),
                          const Text(
                            '-----',
                            style: TextStyle(
                              fontSize: 12,
                              color: Colors.grey,
                            ),
                          ),
                          const Icon(Icons.flight, size: 14, color: Colors.grey),
                          const Text(
                            '-----',
                            style: TextStyle(
                              fontSize: 12,
                              color: Colors.grey,
                            ),
                          ),
                          const SizedBox(width: 4),
                          Text(
                            flight['arrivalTime'],
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black87,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 12),
                    // Price
                    Text(
                      '\$${flight['price'].toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (Match m) => '${m[1]},')}',
                      style: const TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.w700,
                        color: Colors.black87,
                      ),
                    ),
                  ],
                ),
                
                // Connection times (if any)
                if (flight.containsKey('connectTime1'))
                  Padding(
                    padding: const EdgeInsets.only(top: 2, left: 60),
                    child: Text(
                      '${flight['connectTime1']} --- ${flight['connectTime2']}',
                      style: const TextStyle(
                        fontSize: 11,
                        color: Colors.grey,
                      ),
                    ),
                  ),
                
                const SizedBox(height: 8),
                
                // Route
                Text(
                  '${flight['from']} - ${flight['to']}',
                  style: const TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.w700,
                    color: Color(0xFFD32F2F),
                  ),
                ),
                
                const SizedBox(height: 6),
                
                // Baggage Info
                Text(
                  'Checked: ${flight['checkedBag']} / Carry on: ${flight['carryOn']}',
                  style: const TextStyle(
                    fontSize: 10,
                    color: Colors.grey,
                  ),
                ),
                
                const SizedBox(height: 8),
                
                // Airline and Duration Row
                Row(
                  children: [
                    Image.asset(
                      'assets/images/middleeast.png',
                      width: 14,
                      height: 14,
                    ),
                    const SizedBox(width: 4),
                    Expanded(
                      child: Text(
                        flight['airline'],
                        style: const TextStyle(
                          fontSize: 11,
                          color: Color(0xFF1e5a8e),
                        ),
                      ),
                    ),
                    Text(
                      '${flight['duration']} - ${flight['stops']}',
                      style: const TextStyle(
                        fontSize: 11,
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          
          // Flight Details Button
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(vertical: 10),
            decoration: const BoxDecoration(
              color: Color(0xFF2c5f8d),
              borderRadius: BorderRadius.only(
                bottomLeft: Radius.circular(8),
                bottomRight: Radius.circular(8),
              ),
            ),
            child: const Text(
              'Flight Details',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Colors.white,
                fontSize: 13,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
