//go:build darwin

package sleepless

import "os/exec"

func PreventSleep(appName string, reason string) (func(), error) {
	caffeinate := exec.Command("caffeinate", "-d", "-i", "-s", "-m")
	err := caffeinate.Start()
	if err != nil {
		return func() {}, err
	}

	return func() {
		caffeinate.Process.Kill()
		caffeinate.Process.Wait()
	}, nil
}
