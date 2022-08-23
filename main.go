package main

import (
	"context"
	"flag"
	"log"
	"os"

	"byndid/auth-commit-sig/action"
)

func main() {
	log.SetFlags(0)

	path := flag.String("path", ".", "Path to the git repository")
	ref := flag.String("ref", "HEAD", "Commit reference to check")
	flag.Parse()

	cfg := action.Config{
		RepoPath:                *path,
		CommitRef:               *ref,
		APIToken:                getRequiredEnv("API_TOKEN"),
		APIBaseURL:              getOptionalEnv("API_BASE_URL", "https://api.byndid.com/key-mgmt"),
		AllowlistConfigFilePath: getOptionalEnv("ALLOWLIST_CONFIG_FILE_PATH", ""),
	}

	err := action.Run(context.Background(), cfg)
	if err != nil {
		log.Println("Error from action:", err)
		os.Exit(1)
	}
}

func getRequiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Printf("Missing required environment variable: %q", name)
		os.Exit(2)
	}
	return value
}

func getOptionalEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
