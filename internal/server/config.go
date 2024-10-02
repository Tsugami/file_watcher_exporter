package server

type Config struct {
	Host    string
	Port    int
	Dir     string
	Exclude []string
	Include []string
}

func NewConfig() *Config {
	return &Config{
		Host:    "localhost",
		Port:    5428,
		Dir:     "/tmp",
		Exclude: []string{".*"},
		Include: []string{"*.*"},
	}
}
func (c *Config) SetHost(host string) {
	c.Host = host
}

func (c *Config) SetPort(port int) {
	c.Port = port
}

func (c *Config) SetDir(dir string) {
	c.Dir = dir
}

func (c *Config) SetExclude(exclude []string) {
	c.Exclude = exclude
}

func (c *Config) SetInclude(include []string) {
	c.Include = include
}
