package application

import "github.com/spf13/viper"

func LoadConfig(configPath string, configMapping interface{}) *viper.Viper {
	//viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(configMapping); err != nil {
		panic(err)
	}
	return viper.GetViper()
}
