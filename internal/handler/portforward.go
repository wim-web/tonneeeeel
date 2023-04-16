package handler

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/wim-web/tonneeeeel/internal/view"
	"github.com/wim-web/tonneeeeel/pkg/command"
)

func PortforwardHandler(doc command.DocumentName, params map[string][]string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	ssmService := ssm.NewFromConfig(cfg)
	ecsService := ecs.NewFromConfig(cfg)

	cluster, task, container, quit, err := view.Cluster2Task2Container(ecsService)

	if quit {
		return nil
	}
	if err != nil {
		return err
	}

	taskId := strings.Split(*task.TaskArn, "/")[2]

	cmd, err := command.PortForwardCommand(
		context.Background(),
		ssmService,
		cluster,
		taskId,
		*container.RuntimeId,
		cfg.Region,
		doc,
		params,
	)

	if err != nil {
		return err
	}

	return cmd.Run()
}
