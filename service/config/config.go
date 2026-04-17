package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"` // 嵌套 MySQL 配置（yaml 键名 mysql 对应）
	JWT   JWTConfig   `yaml:"jwt"`
}
type MySQLConfig struct {
	Prefix   string `yaml:"prefix"`   // 表前缀（yaml 键名对应）
	Port     string `yaml:"port"`     // 端口
	Config   string `yaml:"config"`   // 连接参数（字符集等）
	DBName   string `yaml:"dbname"`   // 数据库名（yaml 中的 db-name 对应结构体 DBName）
	Username string `yaml:"username"` // 用户名
	Password string `yaml:"password"` // 密码
	Path     string `yaml:"path"`     // 数据库地址（yaml 中的 path 对应结构体 Path）
}

type JWTConfig struct {
	Secret               string `yaml:"secret"`
	Issuer               string `yaml:"issuer"`
	ExpireHours          int    `yaml:"expire_hours"`
	RefreshBeforeMinutes int    `yaml:"refresh_before_minutes"`
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败:%w", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct:%v", err)
	}

}
