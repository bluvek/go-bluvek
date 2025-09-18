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
	RouterPrefix string `mapstructure:"routerPrefix" json:"routerPrefix" yaml:"routerPrefix"`
}

type dbsConf struct {
	Name          string `mapstructure:"name" json:"name" yaml:"name"`
	Driver        string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Dsn           string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	UseGorm       bool   `mapstructure:"useGorm" json:"useGorm" yaml:"useGorm"`
	LogLevel      string `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`
	MaxIdleConn   int    `mapstructure:"maxIdleConn" json:"maxIdleConn" yaml:"maxIdleConn"`
	MaxOpenConn   int    `mapstructure:"maxOpenConn" json:"maxOpenConn" yaml:"maxOpenConn"`
	SlowThreshold int    `mapstructure:"slowThreshold" json:"slowThreshold" yaml:"slowThreshold"`
}

type redisConf struct {
	Addr      string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Pwd       string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`
	Db        int    `mapstructure:"db" json:"db" yaml:"db"`
	IsCluster bool   `mapstructure:"isCluster" json:"isCluster" yaml:"isCluster"`
}

type mongoConf struct {
	URL string `mapstructure:"url" json:"url" yaml:"url"`
}

type logger struct {
	Path       string `mapstructure:"path" json:"path" yaml:"path"`
	Mode       string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Logrotate  bool   `mapstructure:"logrotate" json:"logrotate" yaml:"logrotate"`
	Recover    bool   `mapstructure:"recover" json:"recover" yaml:"recover"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type casbin struct {
	ModePath string `mapstructure:"modePath" json:"modePath" yaml:"modePath"`
}

type jwt struct {
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
	Expire    int    `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type oss struct {
	Type      string `mapstructure:"type" json:"type" yaml:"type"`
	SavePath  string `mapstructure:"savePath" json:"savePath" yaml:"savePath"`
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	AccessKey string `mapstructure:"accessKey" json:"accessKey" yaml:"accessKey"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
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
