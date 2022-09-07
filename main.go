package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
		Repository:              getRequiredEnv("REPOSITORY"),
		AllowlistConfigFilePath: getOptionalEnv("ALLOWLIST_CONFIG_FILE_PATH", ""),
	}

	outcome := action.Run(context.Background(), cfg)
	outcomeJSON, err := jsonMarshal(outcome)
	if err != nil {
		log.Printf("Failed to marshal outcome JSON: %v", err)
		os.Exit(1)
	}
	log.Printf("Outcome JSON: \n%s", string(outcomeJSON))
	stdOut := []byte(fmt.Sprintf(`echo "::set-output name=outcome::%s"`, string(outcomeJSON)))
	_, err = os.Stdout.Write(stdOut)
	if err != nil {
		log.Printf("Failed to set output for this step: %v", err)
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

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
