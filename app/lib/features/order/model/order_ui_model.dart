class OrderModel {
  final String id;
  final String status;
  final double total;
  final DateTime date;

  OrderModel({
    required this.id,
    required this.status,
    required this.total,
    required this.date,
  });

  factory OrderModel.fromJson(Map<String, dynamic> json) {
    return OrderModel(
      id: json['id'] ?? '',
      status: json['status'] ?? '',
      total: (json['total'] ?? 0).toDouble(),
      date: DateTime.parse(json['date'] ?? DateTime.now().toString()),
    );
  }
}
