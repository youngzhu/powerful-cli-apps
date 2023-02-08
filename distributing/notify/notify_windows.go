package notify

import (
	"fmt"
	"os/exec"
)

var command = exec.Command

func (n *Notify) Send() error {
	cmdName := "powershell.exe"

	cmd, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}

	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
$notify = New-Object System.Windows.Forms.NotifyIcon
$notify.Icon = [System.Drawing.SystemIcons]::Information
$notify.BalloonTipIcon = %q
$notify.BalloonTipTitle = %q
$notify.BalloonTipText = %q
$notify.Visible = $True
$notify.ShowBalloonTip(10000)
`, n.severity, n.title, n.message)

	args := []string{
		"-NoProfile",
		"-NonInteractive",
	}
	args = append(args, psScript)

	notifyCommand := command(cmd, args...)
	return notifyCommand.Run()
}
