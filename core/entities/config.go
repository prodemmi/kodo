package entities

type Config struct {
	Flags Flags
}

type Flags struct {
	Port     int
	Silent   bool
	Config   string
	Investor bool
}

func NewDefaultConfig() *Config {
	return &Config{
		Flags: Flags{
			Port:     3519,
			Config:   "./.kodo",
			Silent:   false,
			Investor: false,
		},
	}
}
