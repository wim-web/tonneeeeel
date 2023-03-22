package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "like ecs execute-command",
	Run: func(cmd *cobra.Command, args []string) {
		command, err := cmd.Flags().GetString("command")
		if err != nil {
			log.Fatalln(err)
		}
		err = handler.ExecHandler(command)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().String("command", "bash", "exec command(default: bash)")
}
