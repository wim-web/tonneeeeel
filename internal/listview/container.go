package listview

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type containers []types.Container

func (cs containers) names() (containerList []string) {
	for _, c := range cs {
		containerList = append(containerList, *c.Name)
	}
	return
}

func (cs containers) findByName(name string) (types.Container, bool) {
	for _, c := range cs {
		if *c.Name == name {
			return c, true
		}
	}
	return types.Container{}, false
}

func SelectContainerView(task types.Task, auto bool) (types.Container, bool, error) {
	containers := containers(task.Containers)

	var container types.Container

	if auto && len(containers) == 1 {
		container = containers[0]
	} else if len(containers) > 0 {
		c, quit, err := RenderList("Select a container", containers.names())

		if quit {
			return container, true, nil
		}
		if err != nil {
			return container, false, err
		}

		container, _ = containers.findByName(c)
	} else {
		return container, false, fmt.Errorf("%s don't have container", *task.Group)
	}

	return container, false, nil
}
