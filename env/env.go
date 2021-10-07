package env

import "github.com/kelseyhightower/envconfig"

type TwitterConfig struct {
	BearerToken       string `envconfig:"BEARER_TOKEN" required:"true"`
	ConsumerKey       string `envconfig:"CONSUMER_KEY" required:"true"`
	ConsumerSecret    string `envconfig:"CONSUMER_SECRET" required:"true"`
	AccessToken       string `envconfig:"ACCESS_TOKEN" required:"true"`
	AccessTokenSecret string `envconfig:"ACCESS_TOKEN_SECRET" required:"true"`
}

func GetTwitterConfig() TwitterConfig {
	var con TwitterConfig
	err := envconfig.Process("", &con)
	if err != nil {
		panic(err)
	}
	return con
}
