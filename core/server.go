package core

import (
	"fmt"

	"github.com/bluvek/go-bluvek/pkg/bvutils"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Web 项目服务启动",
	Long:  `通过注册指定路由启动 HTTP 服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer closeService()
		if len(serviceList) <= 0 {
			return fmt.Errorf("请务必通过实现接口 `core.IService` 注册你要启动的服务")
		}

		var eg errgroup.Group
		for _, service := range serviceList {
			service := service
			eg.Go(func() error {
				if err := service.OnStart(); err != nil {
					return fmt.Errorf("服务 %s: %v", bvutils.GetCallerName(service), err)
				}
				return nil
			})
		}
		return nil
	},
}

var serviceList []IService

func RegisterService(service ...IService) {
	serviceList = append(serviceList, service...)
}

func closeService() {
	// 停止服务
	_ = Echo.Sync()
	_ = Log.Sync()
	if rotationSchedulerProcess != nil {
		rotationSchedulerProcess.Stop()
	}
}

type IService interface {
	OnStart() error
}

type IServer struct {
	IService
}

func (s *IServer) OnStart() error {
	// 启动服务
	return nil
}
