package core

import (
	"sync"

	"github.com/bluvek/go-bluvek/pkg/bvcache"
	casbinV2 "github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

var (
	Dbs    sync.Map
	Config *BaseConfig
	Log    *iLog
	Cache  *bvcache.CacheNode
	Casbin *casbinV2.SyncedEnforcer
	Rdb    redis.Cmdable
	Mdb    *mongo.Client
)

func Run() {
	if err := CoreCmd.Execute(); err != nil {
		Echo.Fatalf("❌  服务启动失败: [%s] \n", err)
	}
}

type instance struct {
	Name string
	GORM *gorm.DB
	SQLX *sqlx.DB
}

func SetDb(name string, gdb *gorm.DB, sdb *sqlx.DB) {
	if gdb != nil {
		Dbs.Store(name, &instance{Name: name, GORM: gdb})
	} else if sdb != nil {
		Dbs.Store(name, &instance{Name: name, SQLX: sdb})
	}
}

func Gorm(name ...string) *gorm.DB {
	if len(name) == 0 {
		name = []string{"default"}
	}
	if v, ok := Dbs.Load(name[0]); ok {
		return v.(*instance).GORM
	}

	return nil
}

func Sqlx(name ...string) *sqlx.DB {
	if len(name) == 0 {
		name = []string{"default"}
	}
	if v, ok := Dbs.Load(name[0]); ok {
		return v.(*instance).SQLX
	}

	return nil
}
