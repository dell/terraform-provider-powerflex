package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// go build -o pre-commit pre-commit.go
// cp pre-commit ../.git/hooks/.
// chmod +x ../.git/hooks/pre-commit
// git config hooks.gitleaks true

// gitleaksEnabled determines if the pre-commit hook for gitleaks is enabled.
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

// runCommand runs a shell command and returns an error if it fails.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running %s: %w", name, err)
	}
	return nil
}

func main() {
	// Run golangci-lint
	fmt.Println("Running formating...")
	if err := runCommand("gofmt", "-s", "-w", "."); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run golangci-lint
	fmt.Println("Running golangci-lint...")
	if err := runCommand("checks/golangci-lint", "run", "--fix", "--timeout", "5m"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run gosec
	fmt.Println("Running gosec...")
	if err := runCommand("checks/gosec", "./..."); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run go generate
	fmt.Println("Running go generate...")
	if err := runCommand("go", "generate", "./..."); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check for differences after go generate
	fmt.Println("Checking for differences after go generate...")
	cmd := exec.Command("git", "diff", "--compact-summary", "--exit-code")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			fmt.Println("\nUnexpected difference in directories after code generation. Run 'go generate ./...' command and commit.")
			os.Exit(1)
		}
	}

	// Run gitleaks if enabled
	if gitleaksEnabled() {
		fmt.Println("Running gitleaks...")
		if err := runCommand("checks/gitleaks", "detect", "--no-git", "--config", "checks/gitleaks.toml"); err != nil {
			fmt.Println(`Warning: gitleaks has detected sensitive information in your changes.To disable the gitleaks precommit hook run the following command:git config hooks.gitleaks false`)
			os.Exit(1)
		}
	} else {
		fmt.Println("gitleaks precommit disabled (enable with `git config hooks.gitleaks true`)")
	}
}
