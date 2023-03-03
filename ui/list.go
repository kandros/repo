package ui

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v50/github"
	"golang.design/x/clipboard"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

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
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list                list.Model
	choice              string
	quitting            bool
	showRepoFolderInput bool
	repos               []*github.Repository
	inputModel          inputModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) selectedRepo() *github.Repository {
	return m.repos[m.list.Index()]
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
			clipboard.Write(clipboard.FmtText, []byte(m.selectedRepo().GetHTMLURL()))
			m.quitting = true
			return m, tea.Quit
		case "o":
			selectedRepo := m.selectedRepo()
			err := exec.Command("open", selectedRepo.GetHTMLURL()).Start()
			if err != nil {
				panic(err)
			}
			m.quitting = true
			return m, tea.Quit
		case "c":
			selectedRepo := m.selectedRepo()

			cmd := exec.Command("git", "clone", selectedRepo.GetCloneURL())
			/* renderizzare input per chiedere nome della cartella dove clonare */
			cmd.Stdout = os.Stdout
			err := cmd.Start()
			if err != nil {
				panic(err)
			}
			m.quitting = true
			m.showRepoFolderInput = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showRepoFolderInput {
		return "\n" + m.inputModel.View()
	}

	if m.quitting {

		return quitTextStyle.Render(fmt.Sprintf("%s? copied to clipboard.", m.choice))
	}
	return "\n" + m.list.View()
}

var (
	keyO = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open"),
	)
	keyC = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("c", "clone"),
	)
)

func List(repos []*github.Repository) {
	var items []list.Item

	for _, repo := range repos {
		items = append(items, item(*repo.FullName))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select recent repo"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	customKeys := func() []key.Binding {
		return []key.Binding{keyO}
	}
	l.AdditionalShortHelpKeys = customKeys
	l.AdditionalFullHelpKeys = customKeys

	m := model{list: l, repos: repos}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
