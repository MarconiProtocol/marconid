package msys_cmd_utils

import (
	"../vars"
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Parsing the /etc/os-release file to find out the distro/flavor of linux being run
func ParseLinuxVersion() (int, error) {
	const NAME_KEY = "NAME"
	const VERSION_ID_KEY = "VERSION_ID"
	const OS_RELEASE_FILE = "/etc/os-release"

	const CENTOS_NAME = "CentOS"
	const UBUNTU_NAME = "Ubuntu"

	const VERSION_1604 = "16.04"
	const VERSION_1804 = "18.04"
	const VERSION_7 = "7"

	file, err := os.Open(OS_RELEASE_FILE)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var nameLine string
	var versionLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, NAME_KEY) {
			nameLine = line
		} else if strings.HasPrefix(line, VERSION_ID_KEY) {
			versionLine = line
		}
	}

	if strings.Contains(nameLine, UBUNTU_NAME) {
		if strings.Contains(versionLine, VERSION_1604) {
			return msys_cmd_vars.UBUNTU1604, nil
		} else if strings.Contains(versionLine, VERSION_1804) {
			return msys_cmd_vars.UBUNTU1804, nil
		} else {
			return msys_cmd_vars.UNKNOWN, nil
		}
	} else if strings.Contains(nameLine, CENTOS_NAME) {
		if strings.Contains(versionLine, VERSION_7) {
			return msys_cmd_vars.CENTOS7, nil
		} else {
			return msys_cmd_vars.UNKNOWN, nil
		}
	}

	return 0, nil
}

func RunCmd(cmd *exec.Cmd) (string, error) {
	outWriter := msys_cmd_vars.NewSafeWriter()
	errWriter := msys_cmd_vars.NewSafeWriter()
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	if err := cmd.Run(); err != nil {
		return "", err
	}
	res, err := string(outWriter.GetBytes()), string(errWriter.GetBytes())
	if err != "" {
		return "", errors.New(err)
	}
	return res, nil
}

func RunPipedCmds(cmds ...*exec.Cmd) (string, error) {
	outWriter := msys_cmd_vars.NewSafeWriter()
	errWriter := msys_cmd_vars.NewSafeWriter()

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		cmds[i+1].Stdin, _ = cmd.StdoutPipe()
		cmd.Stderr = errWriter
	}
	cmds[last].Stdout, cmds[last].Stderr = outWriter, errWriter

	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			return "", err
		}
	}
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			return "", err
		}
	}

	res, err := string(outWriter.GetBytes()), string(errWriter.GetBytes())
	if err != "" {
		return "", errors.New(err)
	}
	return res, nil
}