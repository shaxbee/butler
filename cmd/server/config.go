package main

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Addr string `koanf:"addr"`
	DB   string `koanf:"db"`
}

func parseConfig() (Config, error) {
	k := koanf.New(".")
	cfg := Config{}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringSlice("config", []string{}, "config files")
	fs.String("addr", ":8080", "listen address")
	fs.String("db", "sqlite://local.db", "database connection string")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return cfg, err
	}

	files, err := fs.GetStringSlice("config")
	if err != nil {
		return cfg, err
	}

	for _, name := range files {
		if err := k.Load(file.Provider(name), yaml.Parser()); err != nil {
			return cfg, err
		}
	}

	if err := k.Load(posflag.Provider(fs, ".", k), nil); err != nil {
		return cfg, err
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
