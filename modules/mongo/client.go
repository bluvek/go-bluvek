package mongo

import (
	"context"
	"fmt"

	"github.com/bluvek/go-bluvek/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {
	core.RegisterTask(20, mongoCmd)
}

var mongoCmd = &cobra.Command{
	Use:   "mongoDB",
	Short: "Init MongoDB",
	Long:  `加载MongoDB模块之后，可以通过 gooze.Mdb 进行数据操作`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := viper.GetString("Mongo.Url")
		if url == "" {
			return fmt.Errorf("你正在加载MongoDB模块，但是你未配置Mongo.Url，请先添加配置")
		}

		return initClient(url)
	},
}

func initClient(url string) error {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB连接失败: %w", err)
	}

	// 检查连接
	if err = client.Ping(context.TODO(), nil); err != nil {
		return fmt.Errorf("MongoDB连接失败: %w", err)
	}

	core.Mdb = client
	core.Echo.Info("✅ 提示: [Mongo] 模块加载成功, 你可以使用 `core.Mdb` 进行数据操作\n")
	return nil
}
