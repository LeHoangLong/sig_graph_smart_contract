package utility

type SettingsI interface {
	GraphName() string
	MaxTimeDifference_ms() uint64
}

type settings struct {
	graphName string
}

func NewSettings(
	graphName string,
) SettingsI {
	return &settings{
		graphName: graphName,
	}
}

func (s *settings) GraphName() string {
	return s.graphName
}

func (s *settings) MaxTimeDifference_ms() uint64 {
	/// maximum 5 minutes by default
	return 300000
}
