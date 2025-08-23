package core

type Config struct {
	Flags Flags
}

type Flags struct {
	Config   string
	Investor bool
}

func NewDefaultConfig() Config {
	return Config{
		Flags: Flags{
			Config:   "./kodo",
			Investor: false,
		},
	}
}
