package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	MySQL MySQLConfig     `yaml:"mysql"` // 嵌套 MySQL 配置（yaml 键名 mysql 对应）
	JWT   JWTConfig       `yaml:"jwt"`
	App   AppPublicConfig `yaml:"app"` // 可选：公网主站与注入脚本中的注册页地址
}

// AppPublicConfig 用于区分本地与线上（注入到用户 HTML 的脚本无法读前端 .env）。
type AppPublicConfig struct {
	// PortalOrigin 主站前端根地址，不要末尾斜杠，例如 https://lyyxy.top 或 http://localhost:5173
	PortalOrigin string `yaml:"portal_origin"`
	// HtmlPublicHost 用户 HTML 访问用的主机后缀（不含端口），形如 htmlhub.example.com。
	// 完整访问域名为四级：{slug}.HtmlPublicHost（例如 todo.htmlhub.lyyxy.top），三级 htmlhub.* 专用于本项目，根域其它三级可留给其它项目。
	HtmlPublicHost string `yaml:"html_public_host"`
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

	// 依次尝试这些目录（每个目录下应有 config.yaml）
	// 解决宝塔/Supervisor 启动时工作目录不是项目根目录导致找不到配置的问题
	searchDirs := collectConfigDirs()
	for _, dir := range searchDirs {
		viper.AddConfigPath(dir)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(
			"读取配置文件失败: %v\n"+
				"请任选其一：\n"+
				"1) 将 config.yaml 放在可执行文件同级的 config/ 目录下；\n"+
				"2) 设置环境变量 CONFIG_PATH 为包含 config.yaml 的目录（例如 /www/wwwroot/htmlhub/service/config）；\n"+
				"3) 在进程守护里把工作目录设为 service 目录。",
			err,
		)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct:%v", err)
	}
}

func collectConfigDirs() []string {
	var dirs []string
	seen := map[string]bool{}

	add := func(p string) {
		if p == "" {
			return
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return
		}
		if seen[abs] {
			return
		}
		seen[abs] = true
		dirs = append(dirs, abs)
	}

	if env := os.Getenv("CONFIG_PATH"); env != "" {
		add(env)
	}

	if wd, err := os.Getwd(); err == nil {
		add(filepath.Join(wd, "config"))
		add(wd)
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		add(filepath.Join(exeDir, "config"))
		add(filepath.Join(exeDir, "..", "config"))
		add(exeDir)
	}

	add("./config")
	add(".")

	return dirs
}
