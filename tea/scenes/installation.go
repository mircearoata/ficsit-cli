package scenes

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

var _ tea.Model = (*installation)(nil)

type installation struct {
	list         list.Model
	root         components.RootModel
	parent       tea.Model
	installation *cli.Installation
	error        *components.ErrorComponent
	hadRenamed   bool
}

func NewInstallation(root components.RootModel, parent tea.Model, installationData *cli.Installation) tea.Model {
	model := installation{
		root:         root,
		parent:       parent,
		installation: installationData,
	}

	items := []list.Item{
		utils.SimpleItem[installation]{
			ItemTitle: "Select",
			Activate: func(msg tea.Msg, currentModel installation) (tea.Model, tea.Cmd) {
				if err := root.SetCurrentInstallation(installationData); err != nil {
					errorComponent, cmd := components.NewErrorComponent(err.Error(), time.Second*5)
					currentModel.error = errorComponent
					return currentModel, cmd
				}

				return currentModel.parent, nil
			},
		},
		utils.SimpleItem[installation]{
			ItemTitle: "Delete",
			Activate: func(msg tea.Msg, currentModel installation) (tea.Model, tea.Cmd) {
				if err := root.GetGlobal().Installations.DeleteInstallation(installationData.Path); err != nil {
					errorComponent, cmd := components.NewErrorComponent(err.Error(), time.Second*5)
					currentModel.error = errorComponent
					return currentModel, cmd
				}

				return currentModel.parent, updateInstallationListCmd
			},
		},
	}

	model.list = list.New(items, utils.NewItemDelegate(), root.Size().Width, root.Size().Height-root.Height())
	model.list.SetShowStatusBar(false)
	model.list.SetFilteringEnabled(false)
	model.list.Title = fmt.Sprintf("Installation: %s", installationData.Path)
	model.list.Styles = utils.ListStyles
	model.list.SetSize(model.list.Width(), model.list.Height())
	model.list.StatusMessageLifetime = time.Second * 3
	model.list.DisableQuitKeybindings()

	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("q", "back")),
		}
	}

	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("q", "back")),
		}
	}

	return model
}

func (m installation) Init() tea.Cmd {
	return nil
}

func (m installation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case KeyControlC:
			return m, tea.Quit
		case "q":
			if m.parent != nil {
				m.parent.Update(m.root.Size())

				if m.hadRenamed {
					return m.parent, updateInstallationNamesCmd
				}

				return m.parent, nil
			}
			return m, nil
		case KeyEnter:
			i, ok := m.list.SelectedItem().(utils.SimpleItem[installation])
			if ok {
				if i.Activate != nil {
					newModel, cmd := i.Activate(msg, m)
					if newModel != nil || cmd != nil {
						if newModel == nil {
							newModel = m
						}
						return newModel, cmd
					}
					return m, nil
				}
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(2, 2).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	case updateInstallationNames:
		m.hadRenamed = true
		m.list.Title = fmt.Sprintf("Installation: %s", m.installation.Path)
	case components.ErrorComponentTimeoutMsg:
		m.error = nil
	}

	return m, nil
}

func (m installation) View() string {
	if m.error != nil {
		err := m.error.View()
		m.list.SetSize(m.list.Width(), m.root.Size().Height-m.root.Height()-lipgloss.Height(err))
		return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), err, m.list.View())
	}

	m.list.SetSize(m.list.Width(), m.root.Size().Height-m.root.Height())
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}
