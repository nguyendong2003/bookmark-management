package redis

import "github.com/kelseyhightower/envconfig"

type config struct {
	Address  string `default:"localhost:6379" envConfig:"REDIS_ADDRESS"`
	Password string `default:"" envConfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envConfig:"REDIS_DB"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
