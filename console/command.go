package console

import (
	"sort"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type runTask struct {
	Priority int
	Name     string
	Cmd      *cobra.Command
}

var (
	runTasks []runTask
	Echo     *zap.SugaredLogger
	CoreCmd  = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStartupTasks()
		},
	}
)

func RegisterTask(priority int, cmd *cobra.Command) {
	CoreCmd.AddCommand(cmd)
	runTasks = append(runTasks, runTask{
		Priority: priority,
		Name:     cmd.Name(),
		Cmd:      cmd,
	})
}

func runStartupTasks() error {
	// 1. 根据优先级排序，数字越大越先执行
	sort.SliceStable(runTasks, func(i, j int) bool {
		return runTasks[i].Priority > runTasks[j].Priority
	})

	// 2. 按顺序执行任务，任何一个任务失败，立即中止
	for _, task := range runTasks {
		if err := task.Cmd.RunE(task.Cmd, []string{}); err != nil {
			return err
		}
	}

	return nil
}
