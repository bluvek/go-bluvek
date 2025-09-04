package core

import "sync"

var (
	dbList sync.Map
	Config *BaseConfig
	Log    *iLog
)
