class BannerModel {
  final String id;
  final String title;
  final String imagePath;

  BannerModel({
    required this.id,
    required this.title,
    required this.imagePath,
  });

  factory BannerModel.fromJson(Map<String, dynamic> json) {
    return BannerModel(
      id: json['id'] as String,
      title: json['title'] as String,
      imagePath: json['image_path'] as String,
    );
  }
}
