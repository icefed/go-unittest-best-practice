package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBHost     string `yaml:"dbhost"`
	DBPort     int    `yaml:"dbport"`
	DBUser     string `yaml:"dbuser"`
	DBPassword string `json:"-" yaml:"dbpassword"`
	DBName     string `yaml:"dbname"`
	ListenPort int    `yaml:"listenPort"`

	PprofAddr string `yaml:"pprofAddr"`
}

func (c *Config) AddFlags(flags *pflag.FlagSet) {
	flags.SortFlags = false
	flags.StringVar(&c.DBHost, "dbhost", "", "The database host/ip.")
	flags.IntVar(&c.DBPort, "dbport", 3306, "The database port.")
	flags.StringVar(&c.DBUser, "dbuser", "root", "The database user.")
	flags.StringVar(&c.DBPassword, "dbpassword", "", "The database password of user.")
	flags.StringVar(&c.DBName, "dbname", "user_manage", "The database name.")

	flags.IntVar(&c.ListenPort, "listen-port", 8000, "HTTP server listen port.")
	flags.StringVar(&c.PprofAddr, "pprof-addr", ":8090", "The address the pprof endpoint binds to.")
}

func LoadConfig(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("parse yaml failed: %v", err)
	}
	return &conf, err
}