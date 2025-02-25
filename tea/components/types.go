package components

import (
	"github.com/Khan/genqlient/graphql"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/satisfactorymodding/ficsit-cli/cli"
)

type RootModel interface {
	GetGlobal() *cli.GlobalContext

	GetCurrentProfile() *cli.Profile
	SetCurrentProfile(profile *cli.Profile) error

	GetCurrentInstallation() *cli.Installation
	SetCurrentInstallation(installation *cli.Installation) error

	GetAPIClient() graphql.Client

	Size() tea.WindowSizeMsg
	SetSize(size tea.WindowSizeMsg)

	View() string
	Height() int
}
