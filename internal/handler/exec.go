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
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

func listClusters(c *ecs.Client) ([]string, error) {
	input := &ecs.ListClustersInput{}
	res, err := c.ListClusters(context.Background(), input)

	if err != nil {
		return nil, err
	}

	var clusters []string

	for _, arn := range res.ClusterArns {
		v := strings.Split(arn, "/")
		clusters = append(clusters, v[1])
	}

	return clusters, nil
}

func listTasks(c *ecs.Client, cluster string) ([]types.Task, error) {
	ltRes, err := c.ListTasks(context.Background(), &ecs.ListTasksInput{
		Cluster:       aws.String(cluster),
		DesiredStatus: types.DesiredStatusRunning,
	})

	if err != nil {
		return nil, err
	}

	var tasks []types.Task

	dtRes, err := c.DescribeTasks(context.Background(), &ecs.DescribeTasksInput{
		Tasks:   ltRes.TaskArns,
		Cluster: &cluster,
	})

	if err != nil {
		return nil, err
	}

	for _, task := range dtRes.Tasks {
		if strings.HasPrefix(*task.Group, "service") {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
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

	clusters, err := listClusters(ecsService)

	if err != nil {
		return err
	}

	cluster, quit, err := listview.RenderList("Select a cluster", clusters)

	if quit {
		return nil
	}

	if err != nil {
		return err
	}

	tasks, err := listTasks(ecsService, cluster)

	if err != nil {
		return err
	}

	var taskList []string
	taskMap := map[string]types.Task{}

	for _, t := range tasks {
		taskList = append(taskList, *t.Group)
		taskMap[*t.Group] = t
	}

	task, quit, err := listview.RenderList("Select a Task", taskList)

	if quit {
		return nil
	}

	if err != nil {
		return err
	}

	var containers []string

	for _, container := range taskMap[task].Containers {
		containers = append(containers, *container.Name)
	}

	var container *string

	if len(containers) > 1 {
		c, quit, err := listview.RenderList("Select a container", containers)
		if quit {
			return nil
		}

		if err != nil {
			return err
		}

		container = aws.String(c)
	} else if len(containers) == 1 {
		container = aws.String(containers[0])
	} else {
		return fmt.Errorf("%s don't have container", task)
	}

	return execCommand(
		ecsService,
		cluster,
		*taskMap[task].TaskArn,
		command,
		container,
		cfg.Region,
	)
}
