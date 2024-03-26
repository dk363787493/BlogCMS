package config

var Configuration CryptoConfig

type CryptoConfig struct {
	Mysql MysqlConfig
}

type MysqlConfig struct {
	Dsn      string
	Password string
}
