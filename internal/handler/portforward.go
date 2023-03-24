package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/wim-web/tonneeeeel/internal/session_manager"
	"github.com/wim-web/tonneeeeel/internal/view"
)

type DocumentName string

const (
	PORT_FORWARD_DOCUMENT_NAME        DocumentName = "AWS-StartPortForwardingSession"
	REMOTE_PORT_FORWARD_DOCUMENT_NAME DocumentName = "AWS-StartPortForwardingSessionToRemoteHost"
)

func portForwardCommand(c *ssm.Client, cluster string, taskId string, containerId string, region string, doc DocumentName, params map[string][]string) error {
	input := &ssm.StartSessionInput{
		Target:       aws.String(fmt.Sprintf("ecs:%s_%s_%s", cluster, taskId, containerId)),
		DocumentName: aws.String(string(doc)),
		Parameters:   params,
	}

	res, err := c.StartSession(context.Background(), input)

	if err != nil {
		return err
	}

	r, err := json.Marshal(res)

	if err != nil {
		return err
	}

	cmd := session_manager.MakeStartSessionCmd(string(r), region)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PortforwardHandler(doc DocumentName, params map[string][]string) error {
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
		doc,
		params,
	)
}
