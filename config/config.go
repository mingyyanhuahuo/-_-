package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Mysql    MysqlConfig    `yaml:"mysql"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
	Redis    RedisConfig    `yaml:"redis"`
	DeepSeek DeepSeekConfig `yaml:"deepseek"`
	Upload   UploadConfig   `yaml:"upload"`
}
type DeepSeekConfig struct {
	ApiKey              string `mapstructure:"api_key"`
	BaseUrl             string `mapstructure:"base_url"`
	ModelFlash          string `mapstructure:"model-flash"`
	ModelPro            string `mapstructure:"model-pro"`
	ModelFlashVisionExp string `mapstructure:"model-flash-vision-exp"`
}
type MysqlConfig struct {
	Dsn      string `yaml:"dsn"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type JWTConfig struct {
	Secret string `yaml:"secret"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type UploadConfig struct {
	Dir     string `mapstructure:"dir"`
	BaseUrl string `mapstructure:"base_url"`
}

var config *Config

func init() {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}
	if erro := viper.Unmarshal(&config); erro != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", erro))
	}
}
func GetConfig() *Config {
	return config
}
