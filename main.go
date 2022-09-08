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
	ghOut := flag.Bool("github-output", false, "Flag to output outcome to a github action output.")
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

	// Pretty print JSON in logs.
	prettyOutcomeJSON, err := prettyString(string(outcomeJSON))
	if err != nil {
		log.Printf("Failed to pretty outcome JSON: %v", err)
		os.Exit(1)
	} else {
		log.Printf("Outcome JSON: \n%s", prettyOutcomeJSON)
	}

	// If github-outfit is true, produce output.
	if *ghOut {
		fmt.Printf(`::set-output name=outcome::%s`, outcomeJSON)
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
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func prettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
