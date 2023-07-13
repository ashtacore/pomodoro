package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

type model struct {
	title        string
	altscreen    bool
	duration     time.Duration
	passed       time.Duration
	start        time.Time
	timer        timer.Model
	progress     progress.Model
	quitting     bool
	interrupting bool
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		ms := m.duration.Milliseconds()
		if ms == 0 {
			fmt.Println("Duration is zero")
			return m, nil
		}
		m.passed += m.timer.Interval
		pct := m.passed.Milliseconds() * 100 / ms
		cmds = append(cmds, m.progress.SetPercent(float64(pct)/100))

		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		winHeight, winWidth = msg.Height, msg.Width
		if !m.altscreen && m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		focusTime = !focusTime
		if focusTime {
			m.duration = focusDuration
			m.title = focusTitle
		} else {
			m.duration = breakDuration
			m.title = breakTitle
		}
		m.start = time.Now()
		m.passed = 0

		interval := time.Second
		if m.duration < time.Minute {
			interval = 100 * time.Millisecond
		}
		m.timer = timer.NewWithInterval(m.duration, interval)
		//m.progress = progress.New(progress.WithDefaultGradient())
		return m, m.timer.Start()

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		if key.Matches(msg, intKeys) {
			m.interrupting = true
			return m, tea.Quit
		}
		if key.Matches(msg, pauseKeys) {
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting || m.interrupting {
		return "\n"
	}

	result := boldStyle.Render(m.start.Format(time.Kitchen))
	if m.title != "" {
		result += ": " + italicStyle.Render(m.title)
	}
	result += " - " + boldStyle.Render(m.timer.View()) + "\n" + m.progress.View()
	if m.altscreen {
		textWidth, textHeight := lipgloss.Size(result)
		return lipgloss.NewStyle().Margin((winHeight-textHeight)/2, (winWidth-textWidth)/2).Render(result)
	}
	return result
}

var (
	focusTitle          = "Focus"
	breakTitle          = "Chill"
	focusDuration       time.Duration
	breakDuration       time.Duration
	focusTime           = true
	altscreen           bool
	notifications       bool
	sounds              bool
	winHeight, winWidth int
	version             = "dev"
	intKeys             = key.NewBinding(key.WithKeys("esc", "q", "ctrl+c"))
	pauseKeys           = key.NewBinding(key.WithKeys(" "))
	boldStyle           = lipgloss.NewStyle().Bold(true)
	italicStyle         = lipgloss.NewStyle().Italic(true)
)

const (
	padding  = 2
	maxWidth = 80
)

var rootCmd = &cobra.Command{
	Use:          "pomo",
	Short:        "Pomodoro is a system were you focus for an amount of time, then relax for an amount of time",
	Version:      version,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var opts []tea.ProgramOption
		if altscreen {
			opts = append(opts, tea.WithAltScreen())
		}

		addSuffixIfArgIsNumber(&(args[0]), "s")
		addSuffixIfArgIsNumber(&(args[1]), "s")
		duration, err := time.ParseDuration(args[0])
		if err != nil {
			return err
		}
		focusDuration = duration
		breakDuration, err = time.ParseDuration(args[1])
		if err != nil {
			return err
		}

		if !focusTime {
			duration = breakDuration
		}

		interval := time.Second
		if duration < time.Minute {
			interval = 100 * time.Millisecond
		}
		m, err := tea.NewProgram(model{
			duration:  duration,
			timer:     timer.NewWithInterval(duration, interval),
			progress:  progress.New(progress.WithDefaultGradient()),
			title:     focusTitle,
			altscreen: altscreen,
			start:     time.Now(),
		}, opts...).Run()

		if err != nil {
			return err
		}
		if m.(model).interrupting {
			return fmt.Errorf("user exited")
		}
		return nil
	},
}

var manCmd = &cobra.Command{
	Use:                   "man",
	Short:                 "Generates man pages",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		manPage, err := mcobra.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

func init() {
	rootCmd.Flags().StringVarP(&focusTitle, "focus", "f", "Focus", "pomo focus")
	rootCmd.Flags().StringVarP(&breakTitle, "break", "b", "Chill", "pomo break")
	rootCmd.Flags().BoolVarP(&altscreen, "fullscreen", "a", false, "fullscreen")
	rootCmd.Flags().BoolVarP(&notifications, "notify", "n", true, "notify")
	rootCmd.Flags().BoolVarP(&sounds, "sound", "s", false, "sound")

	rootCmd.AddCommand(manCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addSuffixIfArgIsNumber(s *string, suffix string) {
	_, err := strconv.ParseFloat(*s, 64)
	if err == nil {
		*s = *s + suffix
	}
}
