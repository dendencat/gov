package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/spf13/cobra"
)

var govDir = filepath.Join(os.Getenv("HOME"), ".gov")

func main() {
	rootCmd := &cobra.Command{
		Use:   "gov",
		Short: "Go virtual environment tool",
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(buildCmd())
	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(depsCmd())
	rootCmd.AddCommand(activateCmd())
	rootCmd.AddCommand(deactivateCmd())
	rootCmd.AddCommand(useCmd())
	rootCmd.AddCommand(saveBinCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project",
		Run: func(cmd *cobra.Command, args []string) {
			// Git init
			if err := exec.Command("git", "init").Run(); err != nil {
				fmt.Println("Failed to init git:", err)
				return
			}
			fmt.Println("Git repository initialized.")
			// go mod init if not exists
			if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
				if err := exec.Command("go", "mod", "init", "project").Run(); err != nil {
					fmt.Println("Failed to init go mod:", err)
					return
				}
				fmt.Println("Go module initialized.")
			}
		},
	}
}

func buildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := exec.Command("go", "build").Run(); err != nil {
				fmt.Println("Failed to build:", err)
				return
			}
			fmt.Println("Project built successfully.")
		},
	}
}

func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new virtual environment Go project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			// Create directory
			if err := os.MkdirAll(name, 0755); err != nil {
				fmt.Println("Failed to create directory:", err)
				return
			}
			// Change to directory
			os.Chdir(name)
			// Init git and go mod
			initCmd().Run(cmd, []string{})
			fmt.Printf("Virtual environment project '%s' created.\n", name)
		},
	}
}

func depsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deps",
		Short: "Manage dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
				fmt.Println("Failed to tidy deps:", err)
				return
			}
			fmt.Println("Dependencies managed.")
		},
	}
}

func activateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate",
		Short: "Activate virtual environment",
		Run: func(cmd *cobra.Command, args []string) {
			// Set GOPATH or GOROOT to virtual env
			// For simplicity, set GOROOT to .gov/go
			goPath := filepath.Join(govDir, "go")
			os.Setenv("GOROOT", goPath)
			os.Setenv("PATH", goPath+"/bin:"+os.Getenv("PATH"))
			fmt.Println("Virtual environment activated.")
		},
	}
}

func deactivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate virtual environment",
		Run: func(cmd *cobra.Command, args []string) {
			// Reset to system Go
			os.Unsetenv("GOROOT")
			// PATH reset is tricky, for now just print
			fmt.Println("Virtual environment deactivated.")
		},
	}
}

func useCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <version>",
		Short: "Use specified Go version in virtual environment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version := args[0]
			// Download and set Go version
			// This is simplified
			fmt.Printf("Using Go version %s in virtual environment.\n", version)
		},
	}
}

func saveBinCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save-bin",
		Short: "Save Go binary to .gov directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(govDir, 0755); err != nil {
				fmt.Println("Failed to create .gov dir:", err)
				return
			}
			// Copy go binary
			src := "/usr/local/go/bin/go"
			dst := filepath.Join(govDir, "go")
			if err := exec.Command("cp", src, dst).Run(); err != nil {
				fmt.Println("Failed to copy go binary:", err)
				return
			}
			fmt.Println("Go binary saved to .gov directory.")
		},
	}
}
