/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
)

var remoteportforwardCmd = &cobra.Command{
	Use:   "remoteportforward",
	Short: "like start-session --document-name AWS-StartPortForwardingSessionToRemote",
	Run: func(cmd *cobra.Command, args []string) {
		local, err := cmd.Flags().GetString("local-port")
		if err != nil {
			log.Fatalln(err)
		}
		remote, err := cmd.Flags().GetString("remote-port")
		if err != nil {
			log.Fatalln(err)
		}
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatalln(err)
		}

		params := map[string][]string{
			"portNumber":      {remote},
			"localPortNumber": {local},
			"host":            {host},
		}
		err = handler.PortforwardHandler(handler.REMOTE_PORT_FORWARD_DOCUMENT_NAME, params)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(remoteportforwardCmd)

	remoteportforwardCmd.Flags().String("local-port", "", "local port")
	remoteportforwardCmd.Flags().String("remote-port", "", "remote port")
	remoteportforwardCmd.Flags().String("host", "", "host")
}
