/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wim-web/tonneeeeel/internal/handler"
	"github.com/wim-web/tonneeeeel/internal/port"
)

var remoteportforwardCmd = &cobra.Command{
	Use:   "remoteportforward",
	Short: "like start-session --document-name AWS-StartPortForwardingSessionToRemote",
	Run: func(cmd *cobra.Command, args []string) {
		remote, err := cmd.Flags().GetString("remote-port")
		if err != nil {
			log.Fatalln(err)
		}

		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatalln(err)
		}

		local, err := cmd.Flags().GetString("local-port")
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

	remoteportforwardCmd.Flags().StringP("local-port", "l", "", "local port. if not specify, auto assigned")

	remoteportforwardCmd.Flags().StringP("remote-port", "r", "", "remote port")
	remoteportforwardCmd.MarkFlagRequired("remote-port")

	remoteportforwardCmd.Flags().String("host", "", "host")
	remoteportforwardCmd.MarkFlagRequired("host")
}
