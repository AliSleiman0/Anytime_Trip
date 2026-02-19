import 'dart:async';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../core/widgets/global_chatbot_overlay.dart';
import '../../../core/network/endpoints.dart';
import '../service/banner_service.dart';
import '../model/banner_model.dart';
import '../model/popular_location_model.dart';
import 'hotels_search_results.dart';

class HotelsTab extends StatefulWidget {
  final List<PopularLocationModel> popularLocations;
  final bool popularLocationsLoading;
  
  const HotelsTab({
    super.key,
    this.popularLocations = const [],
    this.popularLocationsLoading = false,
  });

  @override
  State<HotelsTab> createState() => _HotelsTabState();
}

class _HotelsTabState extends State<HotelsTab> {
  // Location
  String _destination = 'Bey-Lebanon';
  
  // Date controllers
  DateTime? _checkinDate;
  DateTime? _checkoutDate;
  
  // Travelers
  int _adultsCount = 1;
  int _childrenCount = 2;
  int _roomsCount = 1;
  
  // Nationality
  String _guestNationality = 'Lebanese';
  
  // Carousel
  late PageController _carouselController;
  int _currentCarouselIndex = 0;
  Timer? _carouselTimer;
  final BannerService _bannerService = BannerService();
  List<BannerModel> _banners = [];
  bool _bannersLoading = true;
  
  // Homepage banners (ads)
  List<BannerModel> _homepageBanners = [];
  bool _homepageBannersLoading = true;
  
  // Fallback images if no banners from backend
  final List<String> _fallbackCarouselImages = [
    'assets/images/barcelona.jpg',
    'assets/images/paris.jpg',
    'assets/images/london.jpg',
  ];
  final List<String> _fallbackCarouselTitles = [
    'Barcelona, Spain',
    'Paris, France',
    'London, UK',
  ];

  final List<Map<String, String>> _destinations = [
    {'name': 'Beirut, Lebanon', 'airport': 'Rafic Hariri Intl.'},
    {'name': 'Atlanta, USA', 'airport': 'Hartsfield–Jackson Atlanta Intl.'},
    {'name': 'Madrid, Spain', 'airport': 'Madrid-Barajas Airport'},
    {'name': 'New York, USA', 'airport': 'JFK International'},
    {'name': 'Dubai, UAE', 'airport': 'Dubai International'},
    {'name': 'London, UK', 'airport': 'Heathrow Airport'},
    {'name': 'Paris, France', 'airport': 'Charles de Gaulle Airport'},
    {'name': 'Tokyo, Japan', 'airport': 'Narita International'},
  ];

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

  String get _formattedDates {
    if (_checkinDate == null && _checkoutDate == null) {
      return '07 Nov 22 - 13 Nov 22';
    }
    final checkinStr = _checkinDate != null
        ? '${_checkinDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_checkinDate!.month)} ${_checkinDate!.year.toString().substring(2)}'
        : 'Select';
    final checkoutStr = _checkoutDate != null
        ? '${_checkoutDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_checkoutDate!.month)} ${_checkoutDate!.year.toString().substring(2)}'
        : 'Select';
    return '$checkinStr - $checkoutStr';
  }

  String get _travelersSummary {
    List<String> parts = [];
    if (_adultsCount > 0) parts.add('$_adultsCount Adult${_adultsCount > 1 ? 's' : ''}');
    if (_childrenCount > 0) parts.add('$_childrenCount Child${_childrenCount > 1 ? 'ren' : ''}');
    parts.add('$_roomsCount room${_roomsCount > 1 ? 's' : ''}');
    return parts.join(', ');
  }

  String _getMonthAbbr(int month) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return months[month - 1];
  }

  @override
  void initState() {
    super.initState();
    _carouselController = PageController();
    _loadBanners();
    _loadHomepageBanners();
    _startCarouselTimer();
  }

  void _startCarouselTimer() {
    _carouselTimer?.cancel();
    _carouselTimer = Timer.periodic(const Duration(seconds: 3), (timer) {
      final imageCount = _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length;
      if (_currentCarouselIndex < imageCount - 1) {
        _currentCarouselIndex++;
      } else {
        _currentCarouselIndex = 0;
      }
      if (_carouselController.hasClients) {
        _carouselController.animateToPage(
          _currentCarouselIndex,
          duration: const Duration(milliseconds: 350),
          curve: Curves.easeInOut,
        );
      }
    });
  }

  Future<void> _loadBanners() async {
    try {
      final banners = await _bannerService.getBanners();
      if (mounted) {
        setState(() {
          _banners = banners;
          _bannersLoading = false;
        });
      }
    } catch (e) {
      print('[HOTELS_TAB] Failed to load banners: $e');
      if (mounted) {
        setState(() {
          _bannersLoading = false;
        });
      }
    }
  }

  Future<void> _loadHomepageBanners() async {
    try {
      final banners = await _bannerService.getHomepageBanners();
      if (mounted) {
        setState(() {
          _homepageBanners = banners;
          _homepageBannersLoading = false;
        });
      }
    } catch (e) {
      print('[HOTELS_TAB] Failed to load homepage banners: $e');
      if (mounted) {
        setState(() {
          _homepageBannersLoading = false;
        });
      }
    }
  }

  @override
  void dispose() {
    _carouselController.dispose();
    _carouselTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      physics: const AlwaysScrollableScrollPhysics(),
      child: Padding(
        padding: const EdgeInsets.only(bottom: 24),
        child: Column(
          children: [
            // Going to
            GestureDetector(
              onTap: () => _showLocationBottomSheet(context),
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Going to',
                      style: TextStyle(
                        fontSize: 11,
                        color: const Color(0xFFD32F2F),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Expanded(
                          child: Text(
                            _destination,
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                        const Icon(Icons.location_on, color: Color(0xFFD32F2F), size: 20),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Dates
            GestureDetector(
              onTap: () => _showDateBottomSheet(context),
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Dates',
                          style: TextStyle(
                            fontSize: 11,
                            color: const Color(0xFFD32F2F),
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Row(
                          children: [
                            const Icon(Icons.calendar_today, color: Color(0xFFD32F2F), size: 16),
                            const SizedBox(width: 8),
                            Text(
                              _formattedDates,
                              style: const TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Colors.black,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                    const Icon(Icons.calendar_month, color: Color(0xFFD32F2F), size: 24),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Number of Travelers
            GestureDetector(
              onTap: () => _showTravelersBottomSheet(context),
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Number of Travelers',
                      style: TextStyle(
                        fontSize: 11,
                        color: const Color(0xFFD32F2F),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        const Icon(Icons.person, color: Color(0xFFD32F2F), size: 20),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            _travelersSummary,
                            style: const TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w600,
                              color: Colors.black,
                            ),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Guest Nationality
            GestureDetector(
              onTap: () => _showNationalityBottomSheet(context),
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(15),
                  border: Border.all(color: const Color(0xFF1e5a8e)),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Guest Nationality',
                      style: TextStyle(
                        fontSize: 11,
                        color: const Color(0xFFD32F2F),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          _guestNationality,
                          style: const TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                            color: Colors.black,
                          ),
                        ),
                        const Icon(Icons.check_circle, color: Color(0xFF1e5a8e), size: 24),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            // Search Button
            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => HotelsSearchResults(
                        destination: _destination,
                        checkinDate: _checkinDate != null
                            ? '${_checkinDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_checkinDate!.month)} ${_checkinDate!.year.toString().substring(2)}'
                            : '07 Nov 22',
                        checkoutDate: _checkoutDate != null
                            ? '${_checkoutDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_checkoutDate!.month)} ${_checkoutDate!.year.toString().substring(2)}'
                            : '13 Nov 22',
                        adults: _adultsCount,
                        children: _childrenCount,
                        rooms: _roomsCount,
                        nationality: _guestNationality,
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
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.search, color: Colors.white, size: 20),
                    const SizedBox(width: 8),
                    Text(
                      'search'.tr,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 32),

            // Carousel Section
            _buildDestinationCarousel(),
            const SizedBox(height: 24),

            // Advertisement
            _buildAdvertisement(),
            const SizedBox(height: 32),

            // Popular Locations
            _buildPopularLocations(),
          ],
        ),
      ),
    );
  }

  // Bottom sheet methods
  void _showLocationBottomSheet(BuildContext context) {
    final TextEditingController searchController = TextEditingController();
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            final filteredDestinations = _destinations.where((destination) {
              final name = destination['name']!.toLowerCase();
              final airport = destination['airport']!.toLowerCase();
              final query = searchController.text.toLowerCase();
              return name.contains(query) || airport.contains(query);
            }).toList();

            return DraggableScrollableSheet(
              initialChildSize: 0.6,
              minChildSize: 0.4,
              maxChildSize: 0.95,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text(
                              'Select Destination',
                              style: TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: Colors.black,
                              ),
                            ),
                            GestureDetector(
                              onTap: () => Navigator.pop(context),
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: TextField(
                          controller: searchController,
                          autofocus: false,
                          onChanged: (value) => setModalState(() {}),
                          decoration: InputDecoration(
                            hintText: 'Search destinations...',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            enabledBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(color: Colors.grey[300]!),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(
                                color: Color(0xFF1e5a8e),
                                width: 2,
                              ),
                            ),
                            prefixIcon: const Icon(Icons.search, color: Colors.grey),
                            suffixIcon: searchController.text.isNotEmpty
                                ? GestureDetector(
                                    onTap: () {
                                      searchController.clear();
                                      setModalState(() {});
                                    },
                                    child: const Icon(Icons.clear, color: Colors.grey),
                                  )
                                : const Icon(Icons.location_on, color: Color(0xFFD32F2F)),
                            contentPadding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 12,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      Expanded(
                        child: ListView.builder(
                          controller: scrollController,
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          itemCount: filteredDestinations.length,
                          itemBuilder: (context, index) {
                            final destination = filteredDestinations[index];
                            final isSelected = _destination == destination['name'];

                            return GestureDetector(
                              onTap: () {
                                setState(() {
                                  _destination = destination['name']!;
                                });
                                Navigator.pop(context);
                              },
                              child: Container(
                                margin: const EdgeInsets.only(bottom: 12),
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 16,
                                  vertical: 12,
                                ),
                                decoration: BoxDecoration(
                                  color: isSelected
                                      ? const Color(0xFF1e5a8e).withOpacity(0.1)
                                      : Colors.white,
                                  borderRadius: BorderRadius.circular(12),
                                  border: Border.all(
                                    color: isSelected
                                        ? const Color(0xFF1e5a8e)
                                        : Colors.grey[200]!,
                                    width: isSelected ? 2 : 1,
                                  ),
                                ),
                                child: Row(
                                  children: [
                                    Container(
                                      width: 40,
                                      height: 40,
                                      decoration: BoxDecoration(
                                        color: const Color(0xFFD32F2F).withOpacity(0.1),
                                        borderRadius: BorderRadius.circular(8),
                                      ),
                                      child: const Icon(
                                        Icons.location_on,
                                        color: Color(0xFFD32F2F),
                                        size: 20,
                                      ),
                                    ),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Text(
                                            destination['name']!,
                                            style: TextStyle(
                                              fontSize: 14,
                                              fontWeight: FontWeight.w600,
                                              color: isSelected
                                                  ? const Color(0xFF1e5a8e)
                                                  : Colors.black,
                                            ),
                                          ),
                                          const SizedBox(height: 4),
                                          Text(
                                            destination['airport']!,
                                            style: TextStyle(
                                              fontSize: 12,
                                              color: Colors.grey[600],
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                    if (isSelected)
                                      const Icon(
                                        Icons.check_circle,
                                        color: Color(0xFF1e5a8e),
                                        size: 24,
                                      ),
                                  ],
                                ),
                              ),
                            );
                          },
                        ),
                      ),
                      const SizedBox(height: 24),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  void _showDateBottomSheet(BuildContext context) {
    DateTime tempCheckinDate = _checkinDate ?? DateTime.now();
    DateTime tempCheckoutDate = _checkoutDate ?? DateTime.now().add(const Duration(days: 7));
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return DraggableScrollableSheet(
              initialChildSize: 0.75,
              minChildSize: 0.5,
              maxChildSize: 0.9,
              builder: (BuildContext context, ScrollController scrollController) {
                return Container(
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(20),
                      topRight: Radius.circular(20),
                    ),
                  ),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: Colors.grey[400],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.end,
                          children: [
                            GestureDetector(
                              onTap: () => Navigator.pop(context),
                              child: Container(
                                decoration: BoxDecoration(
                                  color: const Color(0xFFD32F2F),
                                  borderRadius: BorderRadius.circular(50),
                                ),
                                padding: const EdgeInsets.all(4),
                                child: const Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 20,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      Expanded(
                        child: SingleChildScrollView(
                          controller: scrollController,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Check-in',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Color(0xFF1e5a8e),
                                  ),
                                ),
                                const SizedBox(height: 12),
                                _buildCalendar(
                                  context,
                                  tempCheckinDate,
                                  (date) {
                                    setModalState(() {
                                      tempCheckinDate = date;
                                      if (tempCheckoutDate.isBefore(tempCheckinDate)) {
                                        tempCheckoutDate = tempCheckinDate.add(const Duration(days: 1));
                                      }
                                    });
                                  },
                                ),
                                const SizedBox(height: 24),
                                const Text(
                                  'Check-out',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Color(0xFF1e5a8e),
                                  ),
                                ),
                                const SizedBox(height: 12),
                                _buildCalendar(
                                  context,
                                  tempCheckoutDate,
                                  (date) {
                                    setModalState(() {
                                      tempCheckoutDate = date;
                                    });
                                  },
                                  minDate: tempCheckinDate,
                                ),
                                const SizedBox(height: 24),
                              ],
                            ),
                          ),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.all(20),
                        child: SizedBox(
                          width: double.infinity,
                          height: 50,
                          child: ElevatedButton(
                            onPressed: () {
                              setState(() {
                                _checkinDate = tempCheckinDate;
                                _checkoutDate = tempCheckoutDate;
                              });
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
                      ),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  void _showTravelersBottomSheet(BuildContext context) {
    int tempAdults = _adultsCount;
    int tempChildren = _childrenCount;
    int tempRooms = _roomsCount;
    
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      isDismissible: true,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      elevation: 10,
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
              padding: const EdgeInsets.all(20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: Colors.grey[400],
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      GestureDetector(
                        onTap: () => Navigator.pop(context),
                        child: Container(
                          decoration: BoxDecoration(
                            color: const Color(0xFFD32F2F),
                            borderRadius: BorderRadius.circular(50),
                          ),
                          padding: const EdgeInsets.all(4),
                          child: const Icon(
                            Icons.close,
                            color: Colors.white,
                            size: 20,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  _buildTravelerRow(
                    title: 'Adults',
                    count: tempAdults,
                    onDecrement: tempAdults > 1
                        ? () => setModalState(() => tempAdults--)
                        : null,
                    onIncrement: () => setModalState(() => tempAdults++),
                  ),
                  const SizedBox(height: 16),
                  _buildTravelerRow(
                    title: 'Children',
                    count: tempChildren,
                    onDecrement: tempChildren > 0
                        ? () => setModalState(() => tempChildren--)
                        : null,
                    onIncrement: () => setModalState(() => tempChildren++),
                  ),
                  const SizedBox(height: 16),
                  _buildTravelerRow(
                    title: 'Rooms',
                    count: tempRooms,
                    onDecrement: tempRooms > 1
                        ? () => setModalState(() => tempRooms--)
                        : null,
                    onIncrement: () => setModalState(() => tempRooms++),
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 50,
                    child: ElevatedButton(
                      onPressed: () {
                        setState(() {
                          _adultsCount = tempAdults;
                          _childrenCount = tempChildren;
                          _roomsCount = tempRooms;
                        });
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
                  const SizedBox(height: 8),
                ],
              ),
            );
          },
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  void _showNationalityBottomSheet(BuildContext context) {
    isModalOpenNotifier.value = true;
    
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      elevation: 10,
      builder: (BuildContext context) {
        return Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.only(
              topLeft: Radius.circular(20),
              topRight: Radius.circular(20),
            ),
          ),
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[400],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 12),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  GestureDetector(
                    onTap: () => Navigator.pop(context),
                    child: Container(
                      decoration: BoxDecoration(
                        color: const Color(0xFFD32F2F),
                        borderRadius: BorderRadius.circular(50),
                      ),
                      padding: const EdgeInsets.all(4),
                      child: const Icon(
                        Icons.close,
                        color: Colors.white,
                        size: 20,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Flexible(
                child: SingleChildScrollView(
                  child: Column(
                    children: _nationalities.map((nationality) {
                      return GestureDetector(
                        onTap: () {
                          setState(() {
                            _guestNationality = nationality;
                          });
                          Navigator.pop(context);
                        },
                        child: Container(
                          width: double.infinity,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          child: Center(
                            child: Text(
                              nationality,
                              style: TextStyle(
                                fontSize: 16,
                                fontWeight: _guestNationality == nationality
                                    ? FontWeight.w700
                                    : FontWeight.w500,
                                color: const Color(0xFF1e5a8e),
                              ),
                            ),
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
  }

  Widget _buildTravelerRow({
    required String title,
    required int count,
    required VoidCallback? onDecrement,
    required VoidCallback? onIncrement,
  }) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          title,
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: Color(0xFF1e5a8e),
          ),
        ),
        Row(
          children: [
            GestureDetector(
              onTap: onDecrement,
              child: Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  color: onDecrement != null
                      ? const Color(0xFF1e5a8e)
                      : Colors.grey[300],
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.remove,
                  color: onDecrement != null ? Colors.white : Colors.grey[500],
                  size: 20,
                ),
              ),
            ),
            const SizedBox(width: 20),
            SizedBox(
              width: 30,
              child: Text(
                count.toString(),
                textAlign: TextAlign.center,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
            ),
            const SizedBox(width: 20),
            GestureDetector(
              onTap: onIncrement,
              child: Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  color: onIncrement != null
                      ? const Color(0xFF1e5a8e)
                      : Colors.grey[300],
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.add,
                  color: onIncrement != null ? Colors.white : Colors.grey[500],
                  size: 20,
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildCalendar(
    BuildContext context,
    DateTime selectedDate,
    Function(DateTime) onDateSelected, {
    DateTime? minDate,
  }) {
    final now = DateTime.now();
    final daysInMonth = DateTime(selectedDate.year, selectedDate.month + 1, 0).day;
    final firstDayOfWeek = DateTime(selectedDate.year, selectedDate.month, 1).weekday;
    
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              icon: const Icon(Icons.chevron_left),
              onPressed: () {
                final newDate = DateTime(selectedDate.year, selectedDate.month - 1);
                onDateSelected(DateTime(newDate.year, newDate.month, selectedDate.day.clamp(1, DateTime(newDate.year, newDate.month + 1, 0).day)));
              },
            ),
            Text(
              '${_getMonthName(selectedDate.month)} ${selectedDate.year}',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            IconButton(
              icon: const Icon(Icons.chevron_right),
              onPressed: () {
                final newDate = DateTime(selectedDate.year, selectedDate.month + 1);
                onDateSelected(DateTime(newDate.year, newDate.month, selectedDate.day.clamp(1, DateTime(newDate.year, newDate.month + 1, 0).day)));
              },
            ),
          ],
        ),
        const SizedBox(height: 8),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'].map((day) {
            return SizedBox(
              width: 40,
              child: Center(
                child: Text(
                  day,
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: Colors.grey[600],
                  ),
                ),
              ),
            );
          }).toList(),
        ),
        const SizedBox(height: 8),
        ...List.generate((daysInMonth + firstDayOfWeek % 7 + 6) ~/ 7, (weekIndex) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: List.generate(7, (dayIndex) {
                final dayNumber = weekIndex * 7 + dayIndex - (firstDayOfWeek % 7) + 1;
                
                if (dayNumber < 1 || dayNumber > daysInMonth) {
                  return const SizedBox(width: 40, height: 40);
                }
                
                final date = DateTime(selectedDate.year, selectedDate.month, dayNumber);
                final isSelected = date.year == selectedDate.year &&
                    date.month == selectedDate.month &&
                    date.day == selectedDate.day;
                final isPast = minDate != null && date.isBefore(minDate);
                final isDisabled = date.isBefore(DateTime(now.year, now.month, now.day));
                
                return GestureDetector(
                  onTap: (isPast || isDisabled) ? null : () => onDateSelected(date),
                  child: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: isSelected
                          ? const Color(0xFFD32F2F)
                          : Colors.transparent,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Text(
                        dayNumber.toString(),
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: isSelected ? FontWeight.w700 : FontWeight.w500,
                          color: (isPast || isDisabled)
                              ? Colors.grey[400]
                              : isSelected
                                  ? Colors.white
                                  : Colors.black,
                        ),
                      ),
                    ),
                  ),
                );
              }),
            ),
          );
        }),
      ],
    );
  }

  String _getMonthName(int month) {
    const months = [
      'January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'
    ];
    return months[month - 1];
  }

  Widget _buildDestinationCarousel() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.only(bottom: 12),
          child: Text(
            'Get the best travel experience',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w700,
              color: Colors.black,
            ),
          ),
        ),
        SizedBox(
          height: 200,
          child: PageView.builder(
            controller: _carouselController,
            onPageChanged: (index) {
              setState(() {
                _currentCarouselIndex = index;
              });
            },
            itemCount: _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length,
            itemBuilder: (context, index) {
              // Use backend banners if available, otherwise use fallback images
              final imageUrl = _banners.isNotEmpty
                  ? '${Endpoints.baseUrl.replaceAll('/api', '')}${_banners[index].imagePath}'
                  : null;
              final assetImage = _banners.isEmpty ? _fallbackCarouselImages[index] : null;
              final title = _banners.isNotEmpty ? _banners[index].title : _fallbackCarouselTitles[index];

              return Padding(
                padding: const EdgeInsets.only(right: 12),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(15),
                  child: Stack(
                    fit: StackFit.expand,
                    children: [
                      // Image - either from network (backend) or assets (fallback)
                      if (imageUrl != null)
                        Image.network(
                          imageUrl,
                          fit: BoxFit.cover,
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return Container(
                              color: Colors.grey[300],
                              child: Center(
                                child: CircularProgressIndicator(
                                  value: loadingProgress.expectedTotalBytes != null
                                      ? loadingProgress.cumulativeBytesLoaded /
                                          loadingProgress.expectedTotalBytes!
                                      : null,
                                ),
                              ),
                            );
                          },
                          errorBuilder: (context, error, stackTrace) {
                            // Fallback to asset if network image fails
                            return Image.asset(
                              _fallbackCarouselImages[index % _fallbackCarouselImages.length],
                              fit: BoxFit.cover,
                              errorBuilder: (context, error, stackTrace) {
                                return Container(
                                  color: Colors.grey[300],
                                  child: Center(
                                    child: Icon(
                                      Icons.image,
                                      size: 60,
                                      color: Colors.grey[500],
                                    ),
                                  ),
                                );
                              },
                            );
                          },
                        )
                      else if (assetImage != null)
                        Image.asset(
                          assetImage,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) {
                            return Container(
                              color: Colors.grey[300],
                              child: Center(
                                child: Icon(
                                  Icons.image,
                                  size: 60,
                                  color: Colors.grey[500],
                                ),
                              ),
                            );
                          },
                        ),
                      // Gradient overlay
                      Container(
                        decoration: BoxDecoration(
                          gradient: LinearGradient(
                            begin: Alignment.bottomCenter,
                            end: Alignment.topCenter,
                            colors: [
                              Colors.black.withOpacity(0.6),
                              Colors.transparent,
                            ],
                          ),
                        ),
                      ),
                      // Title
                      Positioned(
                        bottom: 16,
                        left: 16,
                        child: Text(
                          title,
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
        const SizedBox(height: 12),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(
            _banners.isNotEmpty ? _banners.length : _fallbackCarouselImages.length,
            (index) => Container(
              width: _currentCarouselIndex == index ? 24 : 8,
              height: 8,
              margin: const EdgeInsets.symmetric(horizontal: 4),
              decoration: BoxDecoration(
                color: _currentCarouselIndex == index
                    ? const Color(0xFF1e5a8e)
                    : Colors.grey[300],
                borderRadius: BorderRadius.circular(4),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildAdvertisement() {
    // Display first homepage banner if available, otherwise show fallback ad
    if (_homepageBanners.isNotEmpty) {
      final banner = _homepageBanners[0];
      final imageUrl = '${Endpoints.baseUrl.replaceAll('/api', '')}${banner.imagePath}';
      
      return ClipRRect(
        borderRadius: BorderRadius.circular(15),
        child: Image.network(
          imageUrl,
          width: double.infinity,
          fit: BoxFit.cover,
          loadingBuilder: (context, child, loadingProgress) {
            if (loadingProgress == null) return child;
            return Container(
              width: double.infinity,
              height: 100,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(15),
              ),
              child: Center(
                child: CircularProgressIndicator(
                  value: loadingProgress.expectedTotalBytes != null
                      ? loadingProgress.cumulativeBytesLoaded /
                          loadingProgress.expectedTotalBytes!
                      : null,
                ),
              ),
            );
          },
          errorBuilder: (context, error, stackTrace) {
            // Fallback to static ad if network image fails
            return Image.asset(
              'assets/images/Ad.png',
              width: double.infinity,
              fit: BoxFit.cover,
              errorBuilder: (context, error, stackTrace) {
                return Container(
                  width: double.infinity,
                  height: 100,
                  decoration: BoxDecoration(
                    color: Colors.grey[300],
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Center(
                    child: Icon(
                      Icons.image,
                      size: 40,
                      color: Colors.grey[500],
                    ),
                  ),
                );
              },
            );
          },
        ),
      );
    }
    
    // Fallback to static ad
    return ClipRRect(
      borderRadius: BorderRadius.circular(15),
      child: Image.asset(
        'assets/images/Ad.png',
        width: double.infinity,
        fit: BoxFit.cover,
        errorBuilder: (context, error, stackTrace) {
          return Container(
            width: double.infinity,
            height: 100,
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(15),
            ),
            child: Center(
              child: Icon(
                Icons.image,
                size: 40,
                color: Colors.grey[500],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildPopularLocations() {
    // Show loading indicator while fetching data
    if (widget.popularLocationsLoading) {
      return const SizedBox(
        height: 220,
        child: Center(
          child: CircularProgressIndicator(
            valueColor: AlwaysStoppedAnimation<Color>(Color(0xFF1e5a8e)),
          ),
        ),
      );
    }

    // Show empty state if no locations available
    if (widget.popularLocations.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.only(bottom: 16),
          child: Text(
            'Popular Locations',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w700,
              color: Colors.black,
            ),
          ),
        ),
        SizedBox(
          height: 220,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: widget.popularLocations.length,
            itemBuilder: (context, index) {
              final location = widget.popularLocations[index];
              return Container(
                width: 240,
                margin: EdgeInsets.only(
                  right: index < widget.popularLocations.length - 1 ? 16 : 0,
                ),
                child: Stack(
                  children: [
                    ClipRRect(
                      borderRadius: BorderRadius.circular(15),
                      child: Stack(
                        fit: StackFit.expand,
                        children: [
                          Image.network(
                            '${Endpoints.baseUrl.replaceAll('/api', '')}${location.imagePath}',
                            fit: BoxFit.cover,
                            errorBuilder: (context, error, stackTrace) {
                              return Container(
                                color: Colors.grey[300],
                                child: Icon(
                                  Icons.image,
                                  size: 60,
                                  color: Colors.grey[500],
                                ),
                              );
                            },
                            loadingBuilder: (context, child, loadingProgress) {
                              if (loadingProgress == null) return child;
                              return Container(
                                color: Colors.grey[200],
                                child: Center(
                                  child: CircularProgressIndicator(
                                    value: loadingProgress.expectedTotalBytes != null
                                        ? loadingProgress.cumulativeBytesLoaded /
                                            loadingProgress.expectedTotalBytes!
                                        : null,
                                    valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF1e5a8e)),
                                  ),
                                ),
                              );
                            },
                          ),
                          Container(
                            decoration: BoxDecoration(
                              gradient: LinearGradient(
                                begin: Alignment.topCenter,
                                end: Alignment.bottomCenter,
                                colors: [
                                  Colors.black.withOpacity(0.3),
                                  Colors.black.withOpacity(0.7),
                                ],
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    if (location.title.isNotEmpty)
                      Positioned(
                        top: 12,
                        left: 12,
                        child: Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 6,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.9),
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Text(
                            location.title,
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                              color: Color(0xFF1e5a8e),
                            ),
                          ),
                        ),
                      ),
                    Positioned(
                      bottom: 12,
                      left: 12,
                      right: 12,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(
                            child: Row(
                              children: [
                                const Icon(
                                  Icons.location_on,
                                  color: Colors.white,
                                  size: 18,
                                ),
                                const SizedBox(width: 4),
                                Expanded(
                                  child: Text(
                                    location.location,
                                    style: const TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.w600,
                                      color: Colors.white,
                                    ),
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(width: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 16,
                              vertical: 8,
                            ),
                            decoration: BoxDecoration(
                              color: const Color(0xFF1e5a8e),
                              borderRadius: BorderRadius.circular(20),
                            ),
                            child: const Text(
                              'Book',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: Colors.white,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
