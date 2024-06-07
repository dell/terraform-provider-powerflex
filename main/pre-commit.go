package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// gitleaksEnabled determines if the pre-commit hook for gitleaks is enabled.
func gitleaksEnabled() bool {
	cmd := exec.Command("git", "config", "--bool", "hooks.gitleaks")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking gitleaks configuration:", err)
		return false
	}
	return strings.TrimSpace(string(out)) != "false"
}

func main() {
	if gitleaksEnabled() {
		cmd1 := exec.Command("pwd")
		out, err1 := cmd1.Output()
		fmt.Print(string(out))
		if err1 != nil {
			fmt.Println("Error getting current directory:", err1)
			os.Exit(1)
		}

		cmd := exec.Command("./gitleaks", "detect", "--no-git", "--config", "gitleaks.toml")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			exitCode := cmd.ProcessState.ExitCode()
			if exitCode == 1 {
				fmt.Println(`Warning: gitleaks has detected sensitive information in your changes.To disable the gitleaks precommit hook run the following command:git config hooks.gitleaks false`)
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("gitleaks precommit disabled (enable with `git config hooks.gitleaks true`)")
	}
}
