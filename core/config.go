package core

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	App    app       `mapstructure:"app" json:"app" yaml:"app"`
	Dbs    []dbsConf `mapstructure:"dbs" json:"dbs" yaml:"dbs"`
	Redis  redisConf `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mongo  mongoConf `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	Logger logger    `mapstructure:"logger" json:"logger" yaml:"logger"`
	Casbin casbin    `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	Jwt    jwt       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Oss    oss       `mapstructure:"oss" json:"oss" yaml:"oss"`
	Cores  cores     `mapstructure:"cores" json:"cores" yaml:"cores"`
}

type app struct {
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	Env          string `mapstructure:"env" json:"env" yaml:"env"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Timeout      int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	RouterPrefix string `mapstructure:"router_prefix" json:"router_prefix" yaml:"router_prefix"`
}

type dbsConf struct {
	Name          string `mapstructure:"name" json:"name" yaml:"name"`
	Driver        string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Dsn           string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	UseGorm       bool   `mapstructure:"use_gorm" json:"use_gorm" yaml:"use_gorm"`
	LogLevel      string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	MaxIdleConns  int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns  int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	SlowThreshold int    `mapstructure:"slow_threshold" json:"slow_threshold" yaml:"slow_threshold"`
}

type redisConf struct {
	Addr      string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Pwd       string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`
	Db        int    `mapstructure:"db" json:"db" yaml:"db"`
	IsCluster bool   `mapstructure:"is_cluster" json:"is_cluster" yaml:"is_cluster"`
}

type mongoConf struct {
	URL string `mapstructure:"url" json:"url" yaml:"url"`
}

type logger struct {
	Path       string `mapstructure:"path" json:"path" yaml:"path"`
	Mode       string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Logrotate  bool   `mapstructure:"logrotate" json:"logrotate" yaml:"logrotate"`
	Recover    bool   `mapstructure:"recover" json:"recover" yaml:"recover"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type casbin struct {
	ModePath string `mapstructure:"mode_path" json:"mode_path" yaml:"mode_path"`
}

type jwt struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Expire    int    `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type oss struct {
	Type      string `mapstructure:"type" json:"type" yaml:"type"`
	SavePath  string `mapstructure:"save_path" json:"save_path" yaml:"save_path"`
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}

type cores struct {
	Mode      string   `mapstructure:"mode" json:"mode" yaml:"mode"`
	Whitelist []string `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`
}

func LoadConfig[T any](file string, env string, target *T) error {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件错误: %s", err)
	}

	if env != "" {
		if err := godotenv.Load(env); err != nil {
			return fmt.Errorf("读取环境变量错误: %s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 强制初始化必要的配置
	if viper.GetInt("App.CacheCap") == 0 {
		viper.Set("App.CacheCap", 100000)
	}
	if viper.GetInt("App.CacheShard") == 0 {
		viper.Set("App.CacheShard", 64)
	}
	if viper.GetString("Casbin.DbName") == "" {
		viper.Set("Casbin.DbName", "default")
	}

	return viper.Unmarshal(target)
}
