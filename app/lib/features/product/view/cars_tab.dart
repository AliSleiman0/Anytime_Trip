import 'dart:async';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../core/widgets/global_chatbot_overlay.dart';
import '../../../core/network/endpoints.dart';
import '../../../core/services/car_service.dart';
import '../service/banner_service.dart';
import '../model/banner_model.dart';
import '../model/popular_location_model.dart';
import 'cars_search_results.dart';

class CarsTab extends StatefulWidget {
  final List<PopularLocationModel> popularLocations;
  final bool popularLocationsLoading;
  
  const CarsTab({
    super.key,
    this.popularLocations = const [],
    this.popularLocationsLoading = false,
  });

  @override
  State<CarsTab> createState() => _CarsTabState();
}

class _CarsTabState extends State<CarsTab> {
  // Location controllers
  String _pickupLocation = '';
  String _dropoffLocation = '';
  
  // Available locations from database
  List<String> _availablePickupLocations = [];
  List<String> _availableDropoffLocations = [];
  bool _locationsLoading = true;
  
  // Date controllers
  DateTime? _pickupDate;
  DateTime? _dropoffDate;
  
  // Time controllers
  String _pickupTime = '10:00am';
  String _dropoffTime = '12:00pm';
  
  // Nationality
  String _guestNationality = 'Lebanese';
  
  // Carousel
  late PageController _carouselController;
  int _currentCarouselIndex = 0;
  Timer? _carouselTimer;
  final BannerService _bannerService = BannerService();
  final CarService _carService = CarService();
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
    if (_pickupDate == null && _dropoffDate == null) {
      return '07 Nov 22 - 13 Nov 22';
    }
    final pickupStr = _pickupDate != null
        ? '${_pickupDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_pickupDate!.month)} ${_pickupDate!.year.toString().substring(2)}'
        : 'Select';
    final dropoffStr = _dropoffDate != null
        ? '${_dropoffDate!.day.toString().padLeft(2, '0')} ${_getMonthAbbr(_dropoffDate!.month)} ${_dropoffDate!.year.toString().substring(2)}'
        : 'Select';
    return '$pickupStr - $dropoffStr';
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
    _loadAvailableLocations();
    _startCarouselTimer();
  }

  Future<void> _loadAvailableLocations() async {
    try {
      final locations = await _carService.getAvailableLocations();
      if (mounted) {
        setState(() {
          _availablePickupLocations = locations['pickup_locations'] ?? [];
          _availableDropoffLocations = locations['dropoff_locations'] ?? [];
          // Set default values from available locations
          if (_pickupLocation.isEmpty && _availablePickupLocations.isNotEmpty) {
            _pickupLocation = _availablePickupLocations.first;
          }
          if (_dropoffLocation.isEmpty && _availableDropoffLocations.isNotEmpty) {
            _dropoffLocation = _availableDropoffLocations.first;
          }
          _locationsLoading = false;
        });
      }
    } catch (e) {
      print('[CARS_TAB] Failed to load available locations: $e');
      if (mounted) {
        setState(() {
          _locationsLoading = false;
        });
      }
    }
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
      print('[CARS_TAB] Failed to load banners: $e');
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
      print('[CARS_TAB] Failed to load homepage banners: $e');
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
            // Pickup - Dropoff Row
            Row(
              children: [
                // Pickup Field
                Expanded(
                  child: GestureDetector(
                    onTap: () => _showLocationBottomSheet(context, isPickup: true),
                    child: Container(
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
                            'Pickup',
                            style: TextStyle(
                              fontSize: 11,
                              color: const Color(0xFFD32F2F),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.location_on, color: Color(0xFFD32F2F), size: 20),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  _pickupLocation,
                                  style: const TextStyle(
                                    fontSize: 14,
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
                ),
                const SizedBox(width: 12),
                // Dropoff Field
                Expanded(
                  child: GestureDetector(
                    onTap: () => _showLocationBottomSheet(context, isPickup: false),
                    child: Container(
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
                            'Drop-off',
                            style: TextStyle(
                              fontSize: 11,
                              color: const Color(0xFFD32F2F),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.location_on, color: Color(0xFFD32F2F), size: 20),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  _dropoffLocation,
                                  style: const TextStyle(
                                    fontSize: 14,
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
                ),
              ],
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

            // Pickup Time - Dropoff Time Row
            Row(
              children: [
                // Pickup Time
                Expanded(
                  child: GestureDetector(
                    onTap: () => _showTimeBottomSheet(context, isPickup: true),
                    child: Container(
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
                            'Pickup time',
                            style: TextStyle(
                              fontSize: 11,
                              color: const Color(0xFFD32F2F),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.access_time, color: Color(0xFFD32F2F), size: 20),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  _pickupTime,
                                  style: const TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                // Dropoff Time
                Expanded(
                  child: GestureDetector(
                    onTap: () => _showTimeBottomSheet(context, isPickup: false),
                    child: Container(
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
                            'Drop-off time',
                            style: TextStyle(
                              fontSize: 11,
                              color: const Color(0xFFD32F2F),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              const Icon(Icons.access_time, color: Color(0xFFD32F2F), size: 20),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  _dropoffTime,
                                  style: const TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.black,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
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
                      builder: (context) => CarsSearchResults(
                        pickupLocation: _pickupLocation,
                        dropoffLocation: _dropoffLocation,
                        pickupDate: _pickupDate,
                        dropoffDate: _dropoffDate,
                        pickupTime: _pickupTime,
                        dropoffTime: _dropoffTime,
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

  // Bottom sheet methods and widget builders...
  void _showLocationBottomSheet(BuildContext context, {required bool isPickup}) {
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
            final availableLocations = isPickup ? _availablePickupLocations : _availableDropoffLocations;
            final filteredDestinations = availableLocations.where((location) {
              final query = searchController.text.toLowerCase();
              return location.toLowerCase().contains(query);
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
                            Text(
                              isPickup ? 'Select Pickup' : 'Select Drop-off',
                              style: const TextStyle(
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
                            final location = filteredDestinations[index];
                            final currentLocation = isPickup ? _pickupLocation : _dropoffLocation;
                            final isSelected = currentLocation == location;

                            return GestureDetector(
                              onTap: () {
                                setState(() {
                                  if (isPickup) {
                                    _pickupLocation = location;
                                  } else {
                                    _dropoffLocation = location;
                                  }
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
                                      child: Text(
                                        location,
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w600,
                                          color: isSelected
                                              ? const Color(0xFF1e5a8e)
                                              : Colors.black,
                                        ),
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
    DateTime? tempPickupDate = _pickupDate;
    DateTime? tempDropoffDate = _dropoffDate;
    DateTime currentMonth = tempPickupDate ?? DateTime.now();
    bool isSelectingDropoff = false;
    
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
                      // Pickup / Drop-off header display
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 20),
                          decoration: BoxDecoration(
                            color: Colors.grey[100],
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Row(
                            children: [
                              Expanded(
                                child: Column(
                                  children: [
                                    Text(
                                      'Pickup',
                                      style: TextStyle(
                                        fontSize: 12,
                                        color: Colors.grey[600],
                                        fontWeight: FontWeight.w500,
                                      ),
                                    ),
                                    const SizedBox(height: 4),
                                    Text(
                                      tempPickupDate != null
                                          ? '${tempPickupDate!.day} ${_getMonthAbbr(tempPickupDate!.month)}'
                                          : 'Select',
                                      style: const TextStyle(
                                        fontSize: 16,
                                        fontWeight: FontWeight.w700,
                                        color: Color(0xFFD32F2F),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              const Icon(Icons.arrow_forward, color: Color(0xFF1e5a8e)),
                              Expanded(
                                child: Column(
                                  children: [
                                    Text(
                                      'Drop-off',
                                      style: TextStyle(
                                        fontSize: 12,
                                        color: Colors.grey[600],
                                        fontWeight: FontWeight.w500,
                                      ),
                                    ),
                                    const SizedBox(height: 4),
                                    Text(
                                      tempDropoffDate != null
                                          ? '${tempDropoffDate!.day} ${_getMonthAbbr(tempDropoffDate!.month)}'
                                          : 'Select',
                                      style: const TextStyle(
                                        fontSize: 16,
                                        fontWeight: FontWeight.w700,
                                        color: Color(0xFFD32F2F),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
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
                                const SizedBox(height: 12),
                                _buildDateRangeCalendar(
                                  context,
                                  currentMonth,
                                  tempPickupDate,
                                  tempDropoffDate,
                                  (newMonth) {
                                    setModalState(() {
                                      currentMonth = newMonth;
                                    });
                                  },
                                  (date) {
                                    setModalState(() {
                                      if (tempPickupDate == null || (tempDropoffDate != null)) {
                                        // First selection or reset
                                        tempPickupDate = date;
                                        tempDropoffDate = null;
                                        isSelectingDropoff = true;
                                      } else if (date.isBefore(tempPickupDate!) || date.isAtSameMomentAs(tempPickupDate!)) {
                                        // Selected before or same as pickup, reset
                                        tempPickupDate = date;
                                        tempDropoffDate = null;
                                        isSelectingDropoff = true;
                                      } else {
                                        // Second selection, set as drop-off
                                        tempDropoffDate = date;
                                        isSelectingDropoff = false;
                                      }
                                    });
                                  },
                                  isSelectingDropoff: isSelectingDropoff,
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
                            onPressed: (tempPickupDate == null || tempDropoffDate == null)
                                ? null
                                : () {
                              setState(() {
                                _pickupDate = tempPickupDate;
                                _dropoffDate = tempDropoffDate;
                              });
                              Navigator.pop(context);
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF1e5a8e),
                              disabledBackgroundColor: Colors.grey[300],
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

  void _showTimeBottomSheet(BuildContext context, {required bool isPickup}) {
    final times = [
      '12:00am', '1:00am', '2:00am', '3:00am', '4:00am', '5:00am',
      '6:00am', '7:00am', '8:00am', '9:00am', '10:00am', '11:00am',
      '12:00pm', '1:00pm', '2:00pm', '3:00pm', '4:00pm', '5:00pm',
      '6:00pm', '7:00pm', '8:00pm', '9:00pm', '10:00pm', '11:00pm',
    ];

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
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    isPickup ? 'Select Pickup Time' : 'Select Drop-off Time',
                    style: const TextStyle(
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
              const SizedBox(height: 16),
              SizedBox(
                height: 300,
                child: GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 4,
                    mainAxisSpacing: 8,
                    crossAxisSpacing: 8,
                    childAspectRatio: 2,
                  ),
                  itemCount: times.length,
                  itemBuilder: (context, index) {
                    final time = times[index];
                    final currentTime = isPickup ? _pickupTime : _dropoffTime;
                    final isSelected = currentTime == time;

                    return GestureDetector(
                      onTap: () {
                        setState(() {
                          if (isPickup) {
                            _pickupTime = time;
                          } else {
                            _dropoffTime = time;
                          }
                        });
                        Navigator.pop(context);
                      },
                      child: Container(
                        decoration: BoxDecoration(
                          color: isSelected
                              ? const Color(0xFF1e5a8e)
                              : Colors.grey[200],
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Center(
                          child: Text(
                            time,
                            style: TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                              color: isSelected ? Colors.white : Colors.black,
                            ),
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ),
            ],
          ),
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
              ..._nationalities.map((nationality) {
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
              const SizedBox(height: 16),
            ],
          ),
        );
      },
    ).then((_) {
      isModalOpenNotifier.value = false;
    });
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

  // Date range calendar for selecting pickup and drop-off dates
  Widget _buildDateRangeCalendar(
    BuildContext context,
    DateTime currentMonth,
    DateTime? pickupDate,
    DateTime? dropoffDate,
    Function(DateTime) onMonthChanged,
    Function(DateTime) onDateTap, {
    bool isSelectingDropoff = false,
  }) {
    final now = DateTime.now();
    final daysInMonth = DateTime(currentMonth.year, currentMonth.month + 1, 0).day;
    final firstDayOfWeek = DateTime(currentMonth.year, currentMonth.month, 1).weekday;
    final currentMonthNormalized = DateTime(now.year, now.month);
    final isPreviousMonthDisabled = currentMonth.year == currentMonthNormalized.year && 
        currentMonth.month == currentMonthNormalized.month;
    
    return Column(
      children: [
        // Month/Year header with navigation
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            IconButton(
              icon: Icon(
                Icons.chevron_left,
                color: isPreviousMonthDisabled ? Colors.grey[300] : Colors.black,
              ),
              onPressed: isPreviousMonthDisabled ? null : () {
                final newDate = DateTime(currentMonth.year, currentMonth.month - 1);
                onMonthChanged(newDate);
              },
            ),
            Text(
              '${_getMonthName(currentMonth.month)} ${currentMonth.year}',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            IconButton(
              icon: const Icon(Icons.chevron_right),
              onPressed: () {
                final newDate = DateTime(currentMonth.year, currentMonth.month + 1);
                onMonthChanged(newDate);
              },
            ),
          ],
        ),
        const SizedBox(height: 8),
        // Weekday headers
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
        // Calendar grid
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
                
                final date = DateTime(currentMonth.year, currentMonth.month, dayNumber);
                final normalizedDate = DateTime(date.year, date.month, date.day);
                final normalizedPickup = pickupDate != null ? DateTime(pickupDate.year, pickupDate.month, pickupDate.day) : null;
                final normalizedDropoff = dropoffDate != null ? DateTime(dropoffDate.year, dropoffDate.month, dropoffDate.day) : null;
                
                final isPickupDate = normalizedPickup != null && normalizedDate.isAtSameMomentAs(normalizedPickup);
                final isDropoffDate = normalizedDropoff != null && normalizedDate.isAtSameMomentAs(normalizedDropoff);
                final isInRange = normalizedPickup != null && normalizedDropoff != null &&
                    normalizedDate.isAfter(normalizedPickup) && normalizedDate.isBefore(normalizedDropoff);
                final isDisabled = date.isBefore(DateTime(now.year, now.month, now.day));
                
                Color? bgColor;
                Color? textColor;
                BorderRadius? borderRadius;
                
                if (isPickupDate || isDropoffDate) {
                  bgColor = const Color(0xFFD32F2F);
                  textColor = Colors.white;
                  if (isPickupDate && isDropoffDate) {
                    borderRadius = BorderRadius.circular(8);
                  } else if (isPickupDate) {
                    borderRadius = const BorderRadius.only(
                      topLeft: Radius.circular(8),
                      bottomLeft: Radius.circular(8),
                    );
                  } else {
                    borderRadius = const BorderRadius.only(
                      topRight: Radius.circular(8),
                      bottomRight: Radius.circular(8),
                    );
                  }
                } else if (isInRange) {
                  bgColor = const Color(0xFFD32F2F).withOpacity(0.15);
                  textColor = Colors.black;
                  borderRadius = BorderRadius.zero;
                } else {
                  bgColor = Colors.transparent;
                  textColor = isDisabled ? Colors.grey[400] : Colors.black;
                  borderRadius = BorderRadius.circular(8);
                }
                
                return GestureDetector(
                  onTap: isDisabled ? null : () => onDateTap(date),
                  child: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: bgColor,
                      borderRadius: borderRadius,
                    ),
                    child: Center(
                      child: Text(
                        dayNumber.toString(),
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: (isPickupDate || isDropoffDate) ? FontWeight.w700 : FontWeight.w500,
                          color: textColor,
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
