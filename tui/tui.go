package tui

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"os"
	"strings"
)

type TUI struct{}

func (t *TUI) NewUploadBar(maxLength int, name string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		int64(maxLength),
		progressbar.OptionSetDescription(fmt.Sprintf("uploading %s", name)),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("/n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	return bar
}

func (t *TUI) AskForUsername() (string, error) {
	return (&promptui.Prompt{
		Label:   "Username",
		Default: "admin",
	}).Run()
}

func (t *TUI) AskForPassword(user string) (string, error) {
	return (&promptui.Prompt{
		Label:       fmt.Sprintf("Password for %q", user),
		HideEntered: true,
		Mask:        '*',
	}).Run()
}

type TaskList struct {
	tasks   []string
	current int
	maxWidth int
}

func (t *TUI) ShowTaskList(tasks ...string) *TaskList {
	var maxWidth int
	for _,t := range tasks {
		if len(t) > maxWidth {
			maxWidth = len(t)
		}
	}
	tl := &TaskList{tasks, 0, maxWidth}
	tl.showNext()
	return tl
}

func (tl *TaskList) Done() {
	fmt.Print("✔\n")
	tl.current += 1
	if tl.current < len(tl.tasks) {
		tl.showNext()
	}
}

func (tl *TaskList) showNext() {
	task := tl.tasks[tl.current]
	padding := strings.Repeat(" ", tl.maxWidth - len(task))
	fmt.Printf("[%d/%d] %s ", tl.current + 1, len(tl.tasks), task + padding)
}
