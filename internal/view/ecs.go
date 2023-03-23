package view

import (
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/wim-web/tonneeeeel/internal/listview"
)

func Cluster2Task2Container(ecsService *ecs.Client) (string, types.Task, types.Container, bool, error) {
	var clusterName string
	var task types.Task
	var container types.Container

	clusterName, quit, err := listview.SelectClusterView(ecsService)

	if quit {
		return clusterName, task, container, true, nil
	}
	if err != nil {
		return clusterName, task, container, false, err
	}

	task, quit, err = listview.SelectTaskView(ecsService, clusterName)

	if quit {
		return clusterName, task, container, true, nil
	}
	if err != nil {
		return clusterName, task, container, false, err
	}

	container, quit, err = listview.SelectContainerView(task, true)

	if quit {
		return clusterName, task, container, true, nil
	}
	if err != nil {
		return clusterName, task, container, false, err
	}

	return clusterName, task, container, false, nil
}
