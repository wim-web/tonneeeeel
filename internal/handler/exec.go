package handler

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/wim-web/tonneeeeel/internal/view"
)

type StartSessionCommandBuilder struct {
	Command       string
	Response      string
	Region        string
	OperationName string
}

func NewStartSessionCommandBuilder(response *ecs.ExecuteCommandOutput, region string) (StartSessionCommandBuilder, error) {
	r, err := json.Marshal(response.Session)

	if err != nil {
		return StartSessionCommandBuilder{}, err
	}

	return StartSessionCommandBuilder{
		Command:       "session-manager-plugin",
		Response:      string(r),
		Region:        region,
		OperationName: "StartSession",
	}, nil

}

func (b StartSessionCommandBuilder) Cmd() *exec.Cmd {
	return exec.Command(
		b.Command,
		b.Response,
		b.Region,
		b.OperationName,
	)
}

func execCommand(c *ecs.Client, cluster string, task string, command string, container *string, region string) error {
	input := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Task:        aws.String(task),
		Interactive: *aws.Bool(true),
		Command:     aws.String(command),
		Container:   container,
	}

	res, err := c.ExecuteCommand(context.Background(), input)

	if err != nil {
		return err
	}

	b, err := NewStartSessionCommandBuilder(res, region)

	if err != nil {
		return err
	}

	cmd := b.Cmd()

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ExecHandler(command string) error {
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

	return execCommand(
		ecsService,
		cluster,
		*task.TaskArn,
		command,
		container.Name,
		cfg.Region,
	)
}
