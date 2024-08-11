package ical

func NewIcalSource(settings *IcalSettings) *IcalSource {
	return &IcalSource{
		settings: settings,
	}
}
