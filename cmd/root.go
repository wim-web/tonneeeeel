package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tonneeeeel",
	Short: "",
	Long:  ``,
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("command", "bash", "exec command(default: bash)")
}
