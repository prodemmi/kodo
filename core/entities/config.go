package entities

type Config struct {
	Flags Flags
}

type Flags struct {
	Silent   bool
	Config   string
	Investor bool
}

func NewDefaultConfig() *Config {
	return &Config{
		Flags: Flags{
			Config:   "./.kodo",
			Silent:   false,
			Investor: false,
		},
	}
}
