//go:build windows

package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed payload/tripol.exe
var compiler []byte

//go:embed payload/hello.tri
var hello []byte

const setPath = `$ErrorActionPreference='Stop'; $dir=$env:TRIPOL_INSTALL_DIR; $p=[Environment]::GetEnvironmentVariable('Path','User'); $parts=@($p -split ';' | Where-Object { $_ -and $_.TrimEnd('\') -ine $dir.TrimEnd('\') }); $new=(@($dir)+$parts)-join ';'; [Environment]::SetEnvironmentVariable('Path',$new,'User')`
const removePath = `$ErrorActionPreference='Stop'; $dir=$env:TRIPOL_INSTALL_DIR; $p=[Environment]::GetEnvironmentVariable('Path','User'); $parts=@($p -split ';' | Where-Object { $_ -and $_.TrimEnd('\') -ine $dir.TrimEnd('\') }); [Environment]::SetEnvironmentVariable('Path',($parts -join ';'),'User')`

func pathChange(dir, script string) error {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Env = append(os.Environ(), "TRIPOL_INSTALL_DIR="+dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func setup(uninstall bool) error {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return fmt.Errorf("LOCALAPPDATA is not defined")
	}
	dir := filepath.Join(base, "Programs", "TRIPOL")
	exe := filepath.Join(dir, "tripol.exe")
	if uninstall {
		if err := pathChange(dir, removePath); err != nil {
			return err
		}
		for _, name := range []string{"tripol.exe"} {
			if err := os.Remove(filepath.Join(dir, name)); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		fmt.Println("TRIPOL removed from PATH. Your example and other files were preserved.")
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(exe, compiler, 0755); err != nil {
		return err
	}
	example := filepath.Join(dir, "hello.tri")
	if _, err := os.Stat(example); os.IsNotExist(err) {
		if err := os.WriteFile(example, hello, 0644); err != nil {
			return err
		}
	}
	if err := pathChange(dir, setPath); err != nil {
		return fmt.Errorf("compiler installed in %s but PATH failed: %w", dir, err)
	}
	fmt.Println("Installed TRIPOL 0.9.0 in:", dir)
	fmt.Println("Close this terminal and open a NEW PowerShell or Command Prompt.")
	fmt.Println("Run: tripol --version")
	fmt.Println("To try immediately in a new terminal: tripol run \"%LOCALAPPDATA%\\Programs\\TRIPOL\\hello.tri\" (Command Prompt)")
	return nil
}
func main() {
	quiet, uninstall := false, false
	for _, arg := range os.Args[1:] {
		switch strings.ToLower(arg) {
		case "--quiet":
			quiet = true
		case "--uninstall":
			uninstall = true
		default:
			fmt.Println("Usage: TRIPOL-Setup.exe [--quiet] [--uninstall]")
			os.Exit(2)
		}
	}
	err := setup(uninstall)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Installation error:", err)
	}
	if !quiet {
		fmt.Print("Press Enter to close...")
		fmt.Scanln()
	}
	if err != nil {
		os.Exit(1)
	}
}
