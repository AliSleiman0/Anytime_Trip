class ProfileModel {
  final String id;
  final String name;
  final String email;
  final String? phone;
  final String? avatar;

  ProfileModel({
    required this.id,
    required this.name,
    required this.email,
    this.phone,
    this.avatar,
  });

  factory ProfileModel.fromJson(Map<String, dynamic> json) {
    return ProfileModel(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'],
      avatar: json['avatar'],
    );
  }
}

class PaymentMethod {
  final String id;
  final String cardHolderName;
  final String cardNumberLast4;
  final String cardBrand;
  final int expiryMonth;
  final int expiryYear;
  final String billingAddress;
  final String city;
  final String state;
  final String postalCode;
  final String country;
  final String phoneNumber;
  final bool isDefault;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  PaymentMethod({
    required this.id,
    required this.cardHolderName,
    required this.cardNumberLast4,
    required this.cardBrand,
    required this.expiryMonth,
    required this.expiryYear,
    required this.billingAddress,
    required this.city,
    required this.state,
    required this.postalCode,
    required this.country,
    required this.phoneNumber,
    required this.isDefault,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory PaymentMethod.fromJson(Map<String, dynamic> json) {
    return PaymentMethod(
      id: json['id'] ?? '',
      cardHolderName: json['card_holder_name'] ?? '',
      cardNumberLast4: json['card_number_last4'] ?? '',
      cardBrand: json['card_brand'] ?? '',
      expiryMonth: json['expiry_month'] ?? 0,
      expiryYear: json['expiry_year'] ?? 0,
      billingAddress: json['billing_address'] ?? '',
      city: json['city'] ?? '',
      state: json['state'] ?? '',
      postalCode: json['postal_code'] ?? '',
      country: json['country'] ?? '',
      phoneNumber: json['phone_number'] ?? '',
      isDefault: json['is_default'] ?? false,
      isActive: json['is_active'] ?? true,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class PaymentMethodRequest {
  final String cardHolderName;
  final String cardNumber;
  final String cardBrand;
  final int expiryMonth;
  final int expiryYear;
  final String cvv;
  final String billingAddress;
  final String city;
  final String state;
  final String postalCode;
  final String country;
  final String phoneNumber;
  final bool isDefault;

  PaymentMethodRequest({
    required this.cardHolderName,
    required this.cardNumber,
    required this.cardBrand,
    required this.expiryMonth,
    required this.expiryYear,
    required this.cvv,
    required this.billingAddress,
    required this.city,
    required this.state,
    required this.postalCode,
    required this.country,
    required this.phoneNumber,
    required this.isDefault,
  });

  Map<String, dynamic> toJson() {
    return {
      'card_holder_name': cardHolderName,
      'card_number': cardNumber,
      'card_brand': cardBrand,
      'expiry_month': expiryMonth,
      'expiry_year': expiryYear,
      'cvv': cvv,
      'billing_address': billingAddress,
      'city': city,
      'state': state,
      'postal_code': postalCode,
      'country': country,
      'phone_number': phoneNumber,
      'is_default': isDefault,
    };
  }
}

class UpdatePaymentMethodRequest {
  final String? cardHolderName;
  final String? billingAddress;
  final String? city;
  final String? state;
  final String? postalCode;
  final String? country;
  final String? phoneNumber;
  final bool? isDefault;

  UpdatePaymentMethodRequest({
    this.cardHolderName,
    this.billingAddress,
    this.city,
    this.state,
    this.postalCode,
    this.country,
    this.phoneNumber,
    this.isDefault,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (cardHolderName != null) json['card_holder_name'] = cardHolderName;
    if (billingAddress != null) json['billing_address'] = billingAddress;
    if (city != null) json['city'] = city;
    if (state != null) json['state'] = state;
    if (postalCode != null) json['postal_code'] = postalCode;
    if (country != null) json['country'] = country;
    if (phoneNumber != null) json['phone_number'] = phoneNumber;
    if (isDefault != null) json['is_default'] = isDefault;
    return json;
  }
}
