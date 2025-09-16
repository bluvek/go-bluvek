package core

import (
	"fmt"
	"os"
	"time"

	"github.com/bluvek/go-bluvek/console"
	"github.com/bluvek/go-bluvek/pkg/bvcache"
	"github.com/bluvek/go-bluvek/pkg/bvutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var configFile string
var env string
var show bool

func init() {
	console.CoreCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	console.CoreCmd.PersistentFlags().StringVar(&env, "env", "", "env file")
	console.CoreCmd.PersistentFlags().BoolVar(&show, "show", true, "Whether to display startup information")
	console.CoreCmd.CompletionOptions.DisableDefaultCmd = true
	console.RegisterTask(-1, serviceCmd)
	console.CoreCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			console.Show(getCommands(), getGlobalFlags())
			os.Exit(-1)
			return nil
		}

		if show {
			console.Show(getCommands(), getGlobalFlags())
		}

		// 1. 初始化配置文件
		if err := LoadConfig(configFile, env, &Config); err != nil {
			console.Echo = initSugaredLogger("")
			return err
		}

		// 2. 初始化 Echo 输出
		console.Echo = initSugaredLogger(Config.App.Env)

		// 3. 初始化日志
		initILog()

		// 4. 初始化缓存模块
		Cache = bvcache.New(viper.GetInt("App.CacheCap"), viper.GetInt("App.CacheShard"), time.Duration(viper.GetInt("App.CacheClear")))

		return nil
	}
}

func initSugaredLogger(env string) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	if !bvutils.InArray(env, []string{"dev", "local", "debug", "test"}) {
		config.OutputPaths = []string{}
	}
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = nil
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.DisableStacktrace = true
	logger, _ := config.Build()

	return logger.Sugar()
}

func getCommands() []console.CommandInfo {
	var commands = []console.CommandInfo{}
	for _, command := range console.CoreCmd.Commands() {
		if !command.Hidden {
			commands = append(commands, console.CommandInfo{
				Name: command.Name(),
				Desc: command.Long,
			})
		}
	}
	return commands
}

func getGlobalFlags() []console.CommandInfo {
	var globalFlags = []console.CommandInfo{}
	console.CoreCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		line := "--" + f.Name
		if f.Shorthand != "" {
			line += ", -" + f.Shorthand
		}
		desc := f.Usage
		if f.DefValue != "" {
			desc += fmt.Sprintf(" (default: %s)", f.DefValue)
		}

		globalFlags = append(globalFlags, console.CommandInfo{
			Name: line,
			Desc: desc,
		})
	})

	return globalFlags
}
