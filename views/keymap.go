package views

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/new-er/zmk-flasher/views/backend"
)

type keyMap struct {
	left    key.Binding
	right   key.Binding
	confirm key.Binding

	automationBindings []key.Binding

	help key.Binding

	quit key.Binding
}

func newKeyMap() keyMap {
	automationBindings := []key.Binding{}
	for i, a := range backend.AutomationStrategies {
		keyStr := strconv.FormatInt(int64(i), 10)
		automationBindings = append(automationBindings,
			key.NewBinding(
				key.WithKeys(keyStr),
				key.WithHelp(keyStr, a.String(-1)),
			))
	}
	return keyMap{
		left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h", "navigate left"),
		),
		right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l", "navigate right"),
		),
		confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		automationBindings: automationBindings,
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.left, k.right},
		k.automationBindings,
		{k.confirm, k.help, k.quit},
	}
}
