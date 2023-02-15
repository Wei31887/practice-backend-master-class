package config

type Config struct {
	DbDriver string  `mapstrucutre:"DBDRIVER"`
	DbSource string	`mapstrucutre:"DBSOURCE"`
	ServerAddress string `mapstrucutre:"SERVERADDRESS"`
}