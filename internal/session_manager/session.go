package session_manager

import (
	"context"
	"os/exec"
)

const SESSION_MANAGER_COMMAND = "session-manager-plugin"

func MakeStartSessionCmd(ctx context.Context, response string, region string) *exec.Cmd {
	const OperationName = "StartSession"

	// https://github.com/aws/session-manager-plugin/blob/1.2.463.0/src/sessionmanagerplugin/session/session.go#L163-L178
	return exec.CommandContext(
		ctx,
		SESSION_MANAGER_COMMAND,
		response,
		region,
		OperationName,
	)
}
