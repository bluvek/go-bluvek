package db

import (
	"encoding/json"
	"fmt"

	"github.com/bluvek/go-bluvek/console"
	"github.com/bluvek/go-bluvek/core"
	"github.com/bluvek/go-bluvek/pkg/bvutils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DbTypeMysql      = "mysql"
	DbTypePostgresql = "postgresql"
	DbTypeSqlite     = "sqlite"
	DbTypeSqlserver  = "sqlserver"
	DbTypeOracle     = "oracle"
)

func init() {
	console.RegisterTask(10, dbCmd)
}

var dbCmd = &cobra.Command{
	Use:    "db",
	Short:  "Init DB",
	Long:   `加载DB模块`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initFunc()
	},
}

type dbConfig struct {
	Name            string
	Dsn             string
	Driver          string
	UseGorm         bool
	LogLevel        int
	EnableLogWriter bool
	MaxIdleConn     int
	MaxConn         int
	SlowThreshold   int
}

func initFunc() error {
	conf := viper.Get(`dbs`)
	confMap, ok := conf.([]any)
	if !ok || len(confMap) == 0 {
		return fmt.Errorf("请确保 `dbs` 模块的配置符合要求")
	}

	isDefault := len(confMap) == 1
	for _, v := range confMap {
		dbConfMap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("请确保 `dbs` 模块的配置符合要求")
		}

		jsonData, err := json.Marshal(dbConfMap)
		if err != nil {
			return fmt.Errorf("请确保 `dbs` 模块的配置符合要求")
		}

		var dbConf dbConfig
		if err = json.Unmarshal(jsonData, &dbConf); err != nil {
			return fmt.Errorf("请确保 `dbs` 模块的配置符合要求")
		}

		// 默认值设置
		dbConf.LogLevel = bvutils.Ternary(dbConf.LogLevel <= 0, 3, dbConf.LogLevel)
		dbConf.MaxConn = bvutils.Ternary(dbConf.MaxConn <= 0, 200, dbConf.MaxConn)
		dbConf.MaxIdleConn = bvutils.Ternary(dbConf.MaxIdleConn <= 0, 10, dbConf.MaxIdleConn)
		dbConf.SlowThreshold = bvutils.Ternary(dbConf.SlowThreshold <= 0, 2000, dbConf.SlowThreshold)

		if dbConf.Dsn == "" || dbConf.Name == "" {
			return fmt.Errorf("你正在加载数据库 [%s] 模块，但配置缺少，请先添加配置", dbConf.Name)
		}

		var funcName string
		if dbConf.UseGorm {
			gdb, err := newGormDB(&dbConf)
			if err != nil {
				return err
			}
			core.SetDb(dbConf.Name, gdb, nil)
			if isDefault {
				core.SetDb("default", gdb, nil)
			}
			funcName = bvutils.Ternary(isDefault, "core.Gorm()", fmt.Sprintf(`core.Gorm("%s")`, dbConf.Name))
		} else {
			sdb, err := newSqlxDB(&dbConf)
			if err != nil {
				return err
			}
			core.SetDb(dbConf.Name, nil, sdb)
			if isDefault {
				core.SetDb("default", nil, sdb)
			}
			funcName = bvutils.Ternary(isDefault, "core.Sqlx()", fmt.Sprintf(`core.Sqlx("%s")`, dbConf.Name))
		}

		console.Echo.Infof("✅  提示: [%s] DB 模块加载成功, 你可以使用 `%s` 进行数据操作\n", dbConf.Name, funcName)
	}

	return nil
}
