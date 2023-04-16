package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
	"github.com/wim-web/tonneeeeel/pkg/command"
	"github.com/wim-web/tonneeeeel/pkg/port"
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

		if local == "" {
			l, err := port.AvailablePort()
			if err != nil {
				log.Fatalln(err)
			}
			local = strconv.Itoa(l)
		}

		params := map[string][]string{
			"portNumber":      {target},
			"localPortNumber": {local},
		}
		err = handler.PortforwardHandler(command.PORT_FORWARD_DOCUMENT_NAME, params)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(portforwardCmd)

	portforwardCmd.Flags().StringP("local-port", "l", "", "local port. if not specify, auto assigned")

	portforwardCmd.Flags().StringP("target-port", "t", "", "target port")
	portforwardCmd.MarkFlagRequired("target-port")
}
