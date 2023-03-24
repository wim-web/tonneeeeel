package session_manager

import "os/exec"

const SESSION_MANAGER_COMMAND = "session-manager-plugin"

func MakeStartSessionCmd(response string, region string) *exec.Cmd {
	const OperationName = "StartSession"

	// https://github.com/aws/session-manager-plugin/blob/1.2.463.0/src/sessionmanagerplugin/session/session.go#L163-L178
	return exec.Command(
		SESSION_MANAGER_COMMAND,
		response,
		region,
		OperationName,
	)
}
