package settings

type Settings struct {
	JwtSecret string `json:"jwtSecret"`
}

type MemorySettingsExporter struct{}

func (s *MemorySettingsExporter) Load() (Settings, error) {
	return Settings{JwtSecret: "secret"}, nil
}
