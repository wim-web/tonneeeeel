package handler

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/wim-web/tonneeeeel/internal/listview"
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

func ExecHandler() error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	ecsService := ecs.NewFromConfig(cfg)

	input2 := &ecs.ListClustersInput{}
	ress, err := ecsService.ListClusters(context.Background(), input2)

	if err != nil {
		return err
	}

	var clusters []string

	for _, arn := range ress.ClusterArns {
		v := strings.Split(arn, "/")
		clusters = append(clusters, v[1])
	}

	cluster, quit, err := listview.RenderList("title", clusters)

	if quit {
		return nil
	}

	if err != nil {
		return err
	}

	input := &ecs.ListTasksInput{
		Cluster: aws.String(cluster),
	}

	res, err := ecsService.ListTasks(context.Background(), input)

	if err != nil {
		return err
	}

	taskArn := res.TaskArns[0]

	execInput := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Task:        aws.String(taskArn),
		Interactive: *aws.Bool(true),
		Command:     aws.String("ash"),
	}

	res2, err := ecsService.ExecuteCommand(context.Background(), execInput)

	if err != nil {
		return err
	}

	b, err := NewStartSessionCommandBuilder(res2, "ap-northeast-1")

	if err != nil {
		return err
	}

	cmd := b.Cmd()

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}
