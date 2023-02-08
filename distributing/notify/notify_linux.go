package notify

import "os/exec"

var command = exec.Command

func (n *Notify) Send() error {
	cmdName := "notify-send"

	cmd, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}

	notifyCommand := command(cmd, "-u", n.severity.String(),
		n.title, n.message)
	return notifyCommand.Run()
}
