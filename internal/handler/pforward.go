package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/wim-web/tonneeeeel/internal/view"
)

type PortForwardCommandBuilder struct {
	Command       string
	Response      string
	Region        string
	OperationName string
}

func NewPortForwardCommandBuilder(response *ssm.StartSessionOutput, region string) (PortForwardCommandBuilder, error) {
	r, err := json.Marshal(response)

	if err != nil {
		return PortForwardCommandBuilder{}, err
	}

	return PortForwardCommandBuilder{
		Command:       "session-manager-plugin",
		Response:      string(r),
		Region:        region,
		OperationName: "StartSession",
	}, nil

}

func (b PortForwardCommandBuilder) Cmd() *exec.Cmd {
	return exec.Command(
		b.Command,
		b.Response,
		b.Region,
		b.OperationName,
	)
}

func portForwardCommand(c *ssm.Client, cluster string, taskId string, containerId string, region string) error {
	input := &ssm.StartSessionInput{
		Target:       aws.String(fmt.Sprintf("ecs:%s_%s_%s", cluster, taskId, containerId)),
		DocumentName: aws.String("AWS-StartPortForwardingSession"),
		Parameters: map[string][]string{
			"portNumber":      {"22"},
			"localPortNumber": {"10022"},
		},
	}

	res, err := c.StartSession(context.Background(), input)

	if err != nil {
		return err
	}

	b, err := NewPortForwardCommandBuilder(res, region)

	if err != nil {
		return err
	}

	cmd := b.Cmd()

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PForwardHandler() error {
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

	return portForwardCommand(
		ssmService,
		cluster,
		taskId,
		*container.RuntimeId,
		cfg.Region,
	)
}
