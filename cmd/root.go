package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

func executeCommand() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalln(err)
	}

	ecsService := ecs.NewFromConfig(cfg)

	input2 := &ecs.ListClustersInput{}
	ress, err := ecsService.ListClusters(context.Background(), input2)

	if err != nil {
		log.Fatalln(err)
	}

	var items []list.Item

	for _, arn := range ress.ClusterArns {
		v := strings.Split(arn, "/")
		items = append(items, item(v[1]))
	}

	l := list.New(items, itemDelegate{}, 20, listHeight)
	m := model{list: l}
	m.list.Title = "select cluster"

	p := tea.NewProgram(m, tea.WithAltScreen())

	mi, err := p.Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	m, ok := mi.(model)

	if !ok {
		log.Fatalln("mi.model failed")
	}

	cluster := m.list.SelectedItem().FilterValue()

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
