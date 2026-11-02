import 'package:flutter/material.dart';

import '../domain/report_reason.dart';

/// What the user filled in on the report sheet.
typedef ReportDraft = ({ReportReason reason, String details});

/// Asks why [name] is being reported. Resolves to null when dismissed.
Future<ReportDraft?> showReportSheet(
  BuildContext context, {
  required String name,
}) => showModalBottomSheet<ReportDraft>(
  context: context,
  isScrollControlled: true,
  showDragHandle: true,
  builder: (_) => _ReportSheet(name: name),
);

class _ReportSheet extends StatefulWidget {
  const _ReportSheet({required this.name});

  final String name;

  @override
  State<_ReportSheet> createState() => _ReportSheetState();
}

class _ReportSheetState extends State<_ReportSheet> {
  /// Must match the backend's limit.
  static const _maxDetails = 500;

  final _details = TextEditingController();
  ReportReason? _reason;

  @override
  void dispose() {
    _details.dispose();
    super.dispose();
  }

  void _submit(ReportReason reason) =>
      Navigator.pop(context, (reason: reason, details: _details.text.trim()));

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final reason = _reason;

    return Padding(
      // Keeps the details field above the keyboard.
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              'Laporkan ${widget.name}',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              'Laporanmu bersifat rahasia. ${widget.name} juga akan diblokir, '
              'jadi kalian tidak akan saling melihat lagi.',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 8),
            RadioGroup<ReportReason>(
              groupValue: reason,
              onChanged: (value) => setState(() => _reason = value),
              child: Column(
                children: [
                  for (final r in ReportReason.values)
                    RadioListTile<ReportReason>(
                      value: r,
                      title: Text(r.label),
                      contentPadding: EdgeInsets.zero,
                    ),
                ],
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: _details,
              minLines: 2,
              maxLines: 4,
              maxLength: _maxDetails,
              textCapitalization: TextCapitalization.sentences,
              decoration: const InputDecoration(
                hintText: 'Ceritakan apa yang terjadi (opsional)',
              ),
            ),
            const SizedBox(height: 12),
            FilledButton(
              onPressed: reason == null ? null : () => _submit(reason),
              child: const Text('Kirim laporan'),
            ),
          ],
        ),
      ),
    );
  }
}
