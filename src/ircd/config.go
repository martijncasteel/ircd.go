package ircd

type config struct {
	motd string
}

func LoadConfiguration(path string) *config {
	return &config{}
}
