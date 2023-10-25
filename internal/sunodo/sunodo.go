// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package sunodo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Get the machine hash from the .sunodo/image/hash file.
func GetMachineHash() (nilhash common.Hash, err error) {
	bytes, err := os.ReadFile(".sunodo/image/hash")
	if err != nil {
		return nilhash, fmt.Errorf("failed to read sunodo image: %v", err)
	}
	if len(bytes) != common.HashLength {
		return nilhash, fmt.Errorf("invalid hash size at .sunodo/image/hash: %v", err)
	}
	return (common.Hash)(bytes), nil
}

// Get the sunodo docker-compose project based on the image hash.
func GetSunodoComposeProject() (string, error) {
	hash, err := GetMachineHash()
	if err != nil {
		return "", err
	}
	project := strings.ToLower(common.Bytes2Hex(hash[:4]))
	return project, nil
}

// Check if sunodo is running in no-backend moode.
func IsNoBackendRunning() (bool, error) {
	output, err := exec.Command("docker", "compose", "ls").CombinedOutput()
	if err != nil {
		msg := "failed to run docker compose ls: %v: %v"
		return false, fmt.Errorf(msg, err, string(output))
	}
	return strings.Contains(string(output), "sunodo-node"), nil
}

// Get the DApp address from the running Anvil container.
func GetDAppAddress() (niladdr common.Address, err error) {
	isRunning, err := IsRunning()
	if err != nil {
		return niladdr, err
	}
	if !isRunning {
		return niladdr, fmt.Errorf("sunodo is not running")
	}

	noBackend, err := IsNoBackendRunning()
	if err != nil {
		return niladdr, err
	}

	var project string
	if noBackend {
		project = "sunodo-node"
	} else {
		project, err = GetSunodoComposeProject()
		if err != nil {
			return niladdr, err
		}
	}

	// Read deployment file from anvil docker container
	const deploymentPath = "/usr/share/sunodo/dapp.json"
	cmd := exec.Command(
		"docker", "compose",
		"--project-name", project,
		"exec", "anvil",
		"cat", deploymentPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		msg := "failed to get deployment file: %v: %v"
		return niladdr, fmt.Errorf(msg, err, stderr.String())
	}

	var deployment struct {
		Address string
	}
	err = json.Unmarshal(output, &deployment)
	if err != nil {
		return niladdr, fmt.Errorf("failed to parse deployment file: %v", err)
	}

	return common.HexToAddress(deployment.Address), nil
}

// Check if sunodo is running.
func IsRunning() (bool, error) {
	output, err := exec.Command("docker", "compose", "ls").CombinedOutput()
	if err != nil {
		msg := "failed to run docker compose ls: %v: %v"
		return false, fmt.Errorf(msg, err, string(output))
	}
	return strings.Contains(string(output), "@sunodo/cli"), nil
}

// Execute the sunodo build command.
func Build(target string, verbose bool) error {
	args := []string{"build"}
	if target != "" {
		args = append(args, "--target", target)
	}
	cmd := exec.Command("sunodo", args...)
	cmd.Stdout = os.Stdout
	if verbose {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return nil
}

// Session of a sunodo run.
type Session struct {
	cmd *exec.Cmd
}

func (s *Session) Close() error {
	// Sending sigint directly to sunodo doesn't work.
	// So, we get the PID of the docker-compose child process and kill it.
	dockerPid, err := procGetChild(s.cmd.Process.Pid)
	if err != nil {
		return fmt.Errorf("failed to get docker pid: %v", err)
	}
	dockerComposePid, err := procGetChild(dockerPid)
	if err != nil {
		return fmt.Errorf("failed to get docker compose pid: %v", err)
	}
	err = procInterrupt(dockerComposePid)
	if err != nil {
		return fmt.Errorf("failed to send interrupt to docker compose: %v", err)
	}
	err = procWait(dockerComposePid, 30)
	if err != nil {
		return fmt.Errorf("failed to wait for docker compose: %v", err)
	}
	return nil
}

// Execute the sunodo run command.
func Run(verbose bool) (*Session, error) {
	cmd := exec.Command("sunodo", "run")
	if verbose {
		cmd.Stderr = os.Stderr
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	ready := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(outPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "Press Ctrl+C to stop the node") {
				ready <- struct{}{}
			}
		}
	}()
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("exec failed: %v", err)
	}
	<-ready

	return &Session{cmd}, nil
}
