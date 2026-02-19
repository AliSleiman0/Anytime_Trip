class PopularLocationModel {
  final String id;
  final String title;
  final String location;
  final String imagePath;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  PopularLocationModel({
    required this.id,
    required this.title,
    required this.location,
    required this.imagePath,
    this.createdAt,
    this.updatedAt,
  });

  factory PopularLocationModel.fromJson(Map<String, dynamic> json) {
    return PopularLocationModel(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      location: json['location'] ?? '',
      imagePath: json['image_path'] ?? '',
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at'])
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at'])
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'location': location,
      'image_path': imagePath,
      'created_at': createdAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }
}
