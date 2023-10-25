// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package sunodo

import (
	"fmt"
	"os/exec"
)

func procGetChild(pid int) (int, error) {
	output, err := exec.Command("pgrep", "-P", fmt.Sprint(pid)).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to exec pgrep: %v: %v", err, string(output))
	}
	var childPid int
	_, err = fmt.Sscanf(string(output), "%d\n", &childPid)
	if err != nil {
		return 0, fmt.Errorf("failed to parse pid: %v", err)
	}
	return childPid, nil
}

func procInterrupt(pid int) error {
	output, err := exec.Command("kill", "-2", fmt.Sprint(pid)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to exec kill: %v: %v", err, string(output))
	}
	return nil
}

func procWait(pid int, timeout int) error {
	output, err := exec.Command(
		"timeout", fmt.Sprint(timeout),
		"tail", "--pid", fmt.Sprint(pid), "-f", "/dev/null").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to exec timeout/tail: %v: %v", err, string(output))
	}
	return nil
}
