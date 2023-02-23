package config

import "time"

type Config struct {
	DbDriver             string        `mapstrucutre:"DBDRIVER"`
	DbSource             string        `mapstrucutre:"DBSOURCE"`
	ServerAddress        string        `mapstrucutre:"SERVERADDRESS"`
	SecreteKey           string        `mapstrucutre:"SECRETEKEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}
