import 'package:flutter/material.dart';

/// Unified AppBar component for checkout and detail screens
/// Provides consistent styling across all screens
class UnifiedAppBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final VoidCallback? onBackPressed;
  final List<Widget>? actions;
  final bool centerTitle;

  const UnifiedAppBar({
    Key? key,
    required this.title,
    this.onBackPressed,
    this.actions,
    this.centerTitle = true,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      backgroundColor: Colors.white,
      elevation: 0,
      leading: IconButton(
        icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e), size: 24),
        onPressed: onBackPressed ?? () => Navigator.pop(context),
      ),
      title: Text(
        title,
        style: const TextStyle(
          color: Color(0xFF1e5a8e),
          fontSize: 16,
          fontWeight: FontWeight.w600,
        ),
      ),
      centerTitle: centerTitle,
      actions: actions,
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}

/// Unified AppBar for product listings with title and subtitle
class UnifiedAppBarWithSubtitle extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final String subtitle;
  final VoidCallback? onBackPressed;
  final List<Widget>? actions;

  const UnifiedAppBarWithSubtitle({
    Key? key,
    required this.title,
    required this.subtitle,
    this.onBackPressed,
    this.actions,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      backgroundColor: Colors.white,
      elevation: 0,
      leading: IconButton(
        icon: const Icon(Icons.arrow_back, color: Color(0xFF1e5a8e), size: 24),
        onPressed: onBackPressed ?? () => Navigator.pop(context),
      ),
      title: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Text(
            title,
            style: const TextStyle(
              color: Color(0xFF1e5a8e),
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          Text(
            subtitle,
            style: const TextStyle(
              color: Colors.grey,
              fontSize: 11,
            ),
          ),
        ],
      ),
      centerTitle: true,
      actions: actions,
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}

/// Unified section header component for consistent styling
class SectionHeader extends StatelessWidget {
  final String title;
  final EdgeInsets padding;

  const SectionHeader({
    Key? key,
    required this.title,
    this.padding = const EdgeInsets.all(0),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w700,
          color: Color(0xFF1e5a8e),
        ),
      ),
    );
  }
}

/// Unified field label component
class FieldLabel extends StatelessWidget {
  final String label;
  final EdgeInsets padding;

  const FieldLabel({
    Key? key,
    required this.label,
    this.padding = const EdgeInsets.only(bottom: 8),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: Text(
        label,
        style: const TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: Color(0xFFD32F2F),
        ),
      ),
    );
  }
}

/// Unified info header component (used in traveler info sections)
class InfoSectionHeader extends StatelessWidget {
  final String title;
  final String subtitle;
  final EdgeInsets padding;

  const InfoSectionHeader({
    Key? key,
    required this.title,
    required this.subtitle,
    this.padding = const EdgeInsets.all(0),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: Color(0xFF1e5a8e),
            ),
          ),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: const TextStyle(
              fontSize: 12,
              color: Colors.grey,
              fontStyle: FontStyle.italic,
            ),
          ),
        ],
      ),
    );
  }
}

/// Unified close button component (for modals and dialogs)
class UnifiedCloseButton extends StatelessWidget {
  final VoidCallback? onPressed;
  final Color color;
  final double size;

  const UnifiedCloseButton({
    Key? key,
    this.onPressed,
    this.color = const Color(0xFF1e5a8e),
    this.size = 24,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Icon(Icons.close, color: color, size: size),
      onPressed: onPressed ?? () => Navigator.pop(context),
    );
  }
}

/// Unified filter button component
class UnifiedFilterButton extends StatelessWidget {
  final VoidCallback onPressed;
  final int? badgeCount;
  final String tooltip;

  const UnifiedFilterButton({
    Key? key,
    required this.onPressed,
    this.badgeCount,
    this.tooltip = 'Filter',
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        IconButton(
          icon: const Icon(Icons.filter_list, color: Color(0xFF1e5a8e), size: 24),
          onPressed: onPressed,
          tooltip: tooltip,
        ),
        if (badgeCount != null && badgeCount! > 0)
          Positioned(
            right: 0,
            top: 0,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 1),
              decoration: BoxDecoration(
                color: const Color(0xFFD32F2F),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Text(
                '$badgeCount',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 10,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ),
      ],
    );
  }
}
