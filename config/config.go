package config

var Configuration CryptoConfig

type CryptoConfig struct {
	Mysql MysqlConfig
	Mongo MongoConfig
}

type MysqlConfig struct {
	Dsn      string
	Password string
}

type MongoConfig struct {
	Url string
}
