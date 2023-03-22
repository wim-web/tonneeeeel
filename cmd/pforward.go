package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
)

var pforwardCmd = &cobra.Command{
	Use:   "pforward",
	Short: "like start-session --document-name AWS-StartPortForwardingSession",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.PForwardHandler()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pforwardCmd)
}
