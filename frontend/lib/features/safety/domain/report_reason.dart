/// Why someone is being reported. [apiValue] must match the reasons the
/// backend accepts on POST /reports.
enum ReportReason {
  fakeProfile('fake_profile', 'Profil palsu atau penipuan'),
  inappropriate('inappropriate', 'Foto atau pesan tidak pantas'),
  harassment('harassment', 'Pelecehan atau perilaku kasar'),
  spam('spam', 'Spam atau promosi'),
  underage('underage', 'Sepertinya di bawah 18 tahun'),
  other('other', 'Lainnya');

  const ReportReason(this.apiValue, this.label);

  final String apiValue;
  final String label;
}
