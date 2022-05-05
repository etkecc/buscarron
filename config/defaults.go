package config

var defaultConfig = &Config{
	Port:     "8080",
	LogLevel: "INFO",
	DB: DB{
		DSN:     "/tmp/buscarron.db",
		Dialect: "sqlite3",
	},
	Ban: &Ban{
		Duration: 24,
		Size:     10000,
	},
}
