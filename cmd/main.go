package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/torjan0/secure-obfuscation-framework-go/obfuscator" // Replace with your module path
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "go-obfuscate",
		Usage: "A tool to obfuscate Go code",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "Build the project with obfuscation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "level",
						Aliases: []string{"l"},
						Value:   "light",
						Usage:   "Obfuscation level: none, light, medium, heavy",
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Enable verbose logging",
					},
				},
				Action: func(c *cli.Context) error {
					level := c.String("level")
					verbose := c.Bool("verbose")

					// Validate obfuscation level
					switch level {
					case "none", "light", "medium", "heavy":
						// Valid levels
					default:
						return fmt.Errorf("invalid obfuscation level: %s", level)
					}

					// Create a temporary directory for obfuscation
					tempDir, err := os.MkdirTemp("", "obfuscated")
					if err != nil {
						return fmt.Errorf("failed to create temp dir: %v", err)
					}
					defer os.RemoveAll(tempDir)
					if verbose {
						log.Printf("Created temp directory: %s", tempDir)
					}

					// Copy source code to temp directory
					err = copyDir(".", tempDir)
					if err != nil {
						return fmt.Errorf("failed to copy directory: %v", err)
					}
					if verbose {
						log.Printf("Copied source to %s", tempDir)
					}

					// Apply obfuscation transformations
					err = transformSource(tempDir, level, verbose)
					if err != nil {
						return fmt.Errorf("failed to transform source: %v", err)
					}

					// Build with garble (assumes garble is installed)
					cmd := exec.Command("garble", "build", "./...")
					cmd.Dir = tempDir
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						return fmt.Errorf("garble build failed: %v", err)
					}
					if verbose {
						log.Println("Build completed successfully")
					}

					// Move the binary back to the current directory
					return os.Rename(filepath.Join(tempDir, "main"), "obfuscated_main")
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// transformSource applies transformations with logging
func transformSource(dir string, level string, verbose bool) error {
	if verbose {
		log.Printf("Applying transformations with level: %s", level)
	}
	return obfuscator.TransformSource(dir, level)
}
