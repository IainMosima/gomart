package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource                string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress       string `mapstructure:"HTTP_SERVER_ADDRESS"`
	AWSRegion               string `mapstructure:"AWS_REGION"`
	CognitoClientID         string `mapstructure:"COGNITO_CLIENT_ID"`
	CognitoClientSecret     string `mapstructure:"COGNITO_CLIENT_SECRET"`
	CognitoRedirectURI      string `mapstructure:"COGNITO_REDIRECT_URI"`
	CognitoDomain           string `mapstructure:"COGNITO_DOMAIN"`
	CognitoUserPoolID       string `mapstructure:"COGNITO_USER_POOL_ID"`
	AfricasTalkingAPIKey    string `mapstructure:"atApiKeys"`
	AfricasTalkingShortCode string `mapstructure:"atShortCode"`
	AfricasTalkingUsername  string `mapstructure:"atUsername"`
	AfricasTalkingSandbox   string `mapstructure:"atSandbox"`
	EmailHost               string `mapstructure:"EMAIL_HOST"`
	EmailPort               string `mapstructure:"EMAIL_PORT"`
	EmailUsername           string `mapstructure:"EMAIL_USERNAME"`
	EmailPassword           string `mapstructure:"EMAIL_PASSWORD"`
	EmailFrom               string `mapstructure:"EMAIL_FROM"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
		return
	} else {
		log.Printf("Using local app.env configuration")
		err = viper.Unmarshal(&config)
		if err != nil {
			return
		}
	}

	return
}
