package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Env string
		Port int
	}

	Database struct {
		Host            string
        Port            int
        User            string
        Password        string
        Name            string
        MaxIdleConns    int
        MaxOpenConns    int
        ConnMaxLifetime string
	}

	Log struct {
		Level string
	}
}

var AppConfig Config

func Load() {
	viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AutomaticEnv()

	viper.SetDefault("app.port", 8080)
    viper.SetDefault("database.max_idle_conns", 10)
    viper.SetDefault("database.max_open_conns", 100)

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("No config file found, using env/default: %v", err)
    } else {
        log.Println("Loaded config:", viper.ConfigFileUsed())
    }

    viper.WatchConfig() // hot reload (optional)

    if err := viper.Unmarshal(&AppConfig); err != nil {
        log.Fatal("Unable to decode config:", err)
    }
}

