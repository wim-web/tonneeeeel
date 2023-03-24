package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
)

var portforwardCmd = &cobra.Command{
	Use:   "portforward",
	Short: "like start-session --document-name AWS-StartPortForwardingSession",
	Run: func(cmd *cobra.Command, args []string) {
		local, err := cmd.Flags().GetString("local-port")
		if err != nil {
			log.Fatalln(err)
		}
		target, err := cmd.Flags().GetString("target-port")
		if err != nil {
			log.Fatalln(err)
		}

		params := map[string][]string{
			"portNumber":      {target},
			"localPortNumber": {local},
		}
		err = handler.PortforwardHandler(handler.PORT_FORWARD_DOCUMENT_NAME, params)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(portforwardCmd)
	portforwardCmd.Flags().String("local-port", "", "local port")
	portforwardCmd.Flags().String("target-port", "", "target port")
}
