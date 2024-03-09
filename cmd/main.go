package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
)

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func main() {
	ignoreDependencies := flag.Bool("ignore-dependencies", false, "Ignore updating dependencies")
	ignoreDevDependencies := flag.Bool("ignore-devDependencies", false, "Ignore updating devDependencies")
	projectPath := flag.String("path", ".", "Path of the project")

	flag.Usage = func() {
		fmt.Printf("%sUsage:%s %s[options]%s\n\n", colorCyan, colorReset, colorYellow, colorReset)
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Printf("  %s --path %s/path/to/project %s--ignore-devDependencies%s\n", os.Args[0], colorYellow, colorYellow, colorReset)
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		ignoreDependencies := false
		ignoreDevDependencies := false
		projectPath := "."

		if err := UpdateDependencies(projectPath, ignoreDependencies, ignoreDevDependencies); err != nil {
			fmt.Println(colorRed+"Error updating dependencies:", err.Error()+colorReset)
			os.Exit(1)
		}
		return
	}

	if *ignoreDependencies && *ignoreDevDependencies {
		fmt.Println(colorRed + "Error: Cannot use both --ignore-dependencies and --ignore-devDependencies flags at the same time." + colorReset)
		flag.Usage()
		os.Exit(1)
	}

	if err := UpdateDependencies(*projectPath, *ignoreDependencies, *ignoreDevDependencies); err != nil {
		fmt.Println(colorRed+"Error updating dependencies:", err.Error()+colorReset)
		os.Exit(1)
	}
}

func UpdateDependencies(projectPath string, ignoreDependencies, ignoreDevDependencies bool) error {

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	fmt.Println(colorCyan + "Please wait. Reading package.json..." + colorReset)

	packageJSONPath := filepath.Join(absPath, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in the specified directory")
	}

	file, err := os.Open(packageJSONPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var pkgJSON PackageJSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pkgJSON); err != nil {
		return err
	}

	ignoredDependenciesCount := 0

	if !ignoreDependencies {
		for dep, version := range pkgJSON.Dependencies {
			if err := updateDependency(absPath, dep, version); err != nil {
				return err
			}
		}
	} else {
		ignoredDependenciesCount += len(pkgJSON.Dependencies)
	}

	if !ignoreDevDependencies {
		for dep, version := range pkgJSON.DevDependencies {
			if err := updateDependency(absPath, dep, version); err != nil {
				return err
			}
		}
	} else {
		ignoredDependenciesCount += len(pkgJSON.DevDependencies)
	}

	fmt.Printf(colorYellow+"%d dependencies ignored."+colorReset+"\n", ignoredDependenciesCount)

	fmt.Println(colorGreen + "Everything updated successfully." + colorReset)
	return nil
}

func updateDependency(absPath, dep, version string) error {
	cmd := exec.Command("npm", "show", dep, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check version for %s: %v", dep, err)
	}

	currentVersion := strings.TrimSpace(string(output))
	if currentVersion == version {
		fmt.Printf("%s is already up to date (current version: %s%s%s)\n", colorYellow+dep+colorReset, colorGreen, version, colorReset)
		return nil
	}

	fmt.Printf("Updating %s to %s%s%s (old version: %s%s%s)...\n", colorYellow+dep+colorReset, colorGreen, version, colorReset, colorRed, currentVersion, colorReset)
	cmd = exec.Command("npm", "install", dep+"@"+version)
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
