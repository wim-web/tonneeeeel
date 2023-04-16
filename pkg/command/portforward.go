package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/wim-web/tonneeeeel/internal/session_manager"
)

type DocumentName string

const (
	PORT_FORWARD_DOCUMENT_NAME        DocumentName = "AWS-StartPortForwardingSession"
	REMOTE_PORT_FORWARD_DOCUMENT_NAME DocumentName = "AWS-StartPortForwardingSessionToRemoteHost"
)

func PortForwardCommand(c *ssm.Client, cluster string, taskId string, containerId string, region string, doc DocumentName, params map[string][]string) error {
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