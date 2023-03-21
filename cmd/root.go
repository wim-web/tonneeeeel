package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tonneeeeel",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tonneeeeel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeCommand() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalln(err)
	}

	ecsService := ecs.NewFromConfig(cfg)

	cluster := ""

	input := &ecs.ListTasksInput{
		Cluster: aws.String(cluster),
	}

	res, err := ecsService.ListTasks(context.Background(), input)

	if err != nil {
		log.Fatalln(err)
	}

	taskArn := res.TaskArns[0]

	execInput := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Task:        aws.String(taskArn),
		Interactive: *aws.Bool(true),
		Command:     aws.String("ash"),
	}

	res2, err := ecsService.ExecuteCommand(context.Background(), execInput)

	if err != nil {
		log.Fatalln(err)
	}

	b, err := NewStartSessionCommandBuilder(res2, "ap-northeast-1")

	if err != nil {
		log.Fatalln(err)
	}

	cmd := b.Cmd()

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

type StartSessionCommandBuilder struct {
	Command       string
	Response      string
	Region        string
	OperationName string
}

func NewStartSessionCommandBuilder(response *ecs.ExecuteCommandOutput, region string) (StartSessionCommandBuilder, error) {
	r, err := json.Marshal(response.Session)

	if err != nil {
		return StartSessionCommandBuilder{}, err
	}

	return StartSessionCommandBuilder{
		Command:       "session-manager-plugin",
		Response:      string(r),
		Region:        region,
		OperationName: "StartSession",
	}, nil

}

func (b StartSessionCommandBuilder) Cmd() *exec.Cmd {
	return exec.Command(
		b.Command,
		b.Response,
		b.Region,
		b.OperationName,
	)
}
