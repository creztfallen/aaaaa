package utils

import (
	"fmt"
	"os/exec"
)

func ExecChain(containerName string) (string, error) {
	output, err := exec.Command("docker", fmt.Sprintf("cp ./cmd/api/compile/config-file-downloads/def-decry-cfg %s:/root/", containerName)).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
