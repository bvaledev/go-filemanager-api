package configs

import "github.com/spf13/viper"

type Conf struct {
	Port       string `mapstructure:"PORT"`
	S3Key      string `mapstructure:"S3_KEY"`
	S3Secret   string `mapstructure:"S3_SECRET"`
	S3Bucket   string `mapstructure:"S3_BUCKET"`
	S3Endpoint string `mapstructure:"S3_ENDPOINT"`
	S3Region   string `mapstructure:"S3_REGION"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
