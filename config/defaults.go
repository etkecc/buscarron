package config

var defaultConfig = &Config{
	Port:     "8080",
	LogLevel: "INFO",
	DB: DB{
		DSN:     "/tmp/buscarron.db",
		Dialect: "sqlite3",
	},
}
