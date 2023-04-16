package handler

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/wim-web/tonneeeeel/internal/view"
	"github.com/wim-web/tonneeeeel/pkg/command"
)

func ExecHandler(cmd string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	ecsService := ecs.NewFromConfig(cfg)

	cluster, task, container, quit, err := view.Cluster2Task2Container(ecsService)

	if quit {
		return nil
	}
	if err != nil {
		return err
	}

	exeCmd, err := command.ExecCommand(
		context.Background(),
		ecsService,
		cluster,
		*task.TaskArn,
		cmd,
		container.Name,
		cfg.Region,
	)

	if err != nil {
		return err
	}

	return exeCmd.Run()
}
