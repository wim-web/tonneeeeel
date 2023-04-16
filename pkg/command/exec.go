package command

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/wim-web/tonneeeeel/internal/session_manager"
)

func ExecCommand(ctx context.Context, c *ecs.Client, cluster string, task string, command string, container *string, region string) (*exec.Cmd, error) {
	input := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Task:        aws.String(task),
		Interactive: *aws.Bool(true),
		Command:     aws.String(command),
		Container:   container,
	}

	res, err := c.ExecuteCommand(context.Background(), input)

	if err != nil {
		return nil, err
	}

	r, err := json.Marshal(res.Session)

	if err != nil {
		return nil, err
	}

	cmd := session_manager.MakeStartSessionCmd(ctx, string(r), region)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd, nil
}
