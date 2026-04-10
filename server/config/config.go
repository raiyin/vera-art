package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name string
	Port int
}

type DirectoriesConfig struct {
	WorksDirSave     string `mapstructure:"works_dir_save"`
	WorksDbDirPrefix string `mapstructure:"works_db_dir_prefix"`
	SaleDirSave      string `mapstructure:"sale_dir_save"`
	SaleDbDirPrefix  string `mapstructure:"sale_db_dir_prefix"`
	NewsDirSave      string `mapstructure:"news_dir_save"`
	NewsDbDirPrefix  string `mapstructure:"news_db_dir_prefix"`
}

type Config struct {
	App         AppConfig
	Directories DirectoriesConfig
}

var AppConfigInstance Config

func LoadConfig(configPath string) error {
	// Set the file name of the configurations file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // path to look for the config file in
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("app.port", 8000)
	viper.SetDefault("app.name", "artserver")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(&AppConfigInstance); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	log.Println("Configuration loaded successfully")
	return nil
}
