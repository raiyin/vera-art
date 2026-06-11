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

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	Debug          bool     `mapstructure:"debug"`
}

type DirectoriesConfig struct {
	AbsWorksDir   string `mapstructure:"abs_works_dir"`
	RelWorksDir   string `mapstructure:"rel_works_dir"`
	AbsSalesDir   string `mapstructure:"abs_sales_dir"`
	RelSalesDir   string `mapstructure:"rel_sales_dir"`
	AbsNewsDir    string `mapstructure:"abs_news_dir"`
	RelNewsDir    string `mapstructure:"rel_news_dir"`
	AbsAvatarsDir string `mapstructure:"abs_avatars_dir"`
}

type Config struct {
	App         AppConfig
	CORS        CORSConfig
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
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://127.0.0.1:3000"})
	viper.SetDefault("cors.debug", false)

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
