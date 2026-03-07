package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Artifact struct {
	ArtifactID string `json:"artifact_id"`
}

type LogEntry struct {
	ArtifactHash string `json:"artifact_hash"`
}

func baseDir() string {

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exe)
}

func artifactsDir() string {
	return filepath.Join(baseDir(), "artifacts")
}

func logPath() string {
	return filepath.Join(baseDir(), "log.json")
}

func appendArtifact(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading artifact:", err)
		return
	}

	var artifact Artifact

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		fmt.Println("Invalid artifact format:", err)
		return
	}

	if artifact.ArtifactID == "" {
		fmt.Println("artifact_id missing")
		return
	}

	err = os.MkdirAll(artifactsDir(), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating artifacts directory:", err)
		return
	}

	target := filepath.Join(artifactsDir(), artifact.ArtifactID+".json")

	err = os.WriteFile(target, data, 0644)
	if err != nil {
		fmt.Println("Error writing artifact:", err)
		return
	}

	var logEntries []LogEntry

	if _, err := os.Stat(logPath()); err == nil {

		logData, err := os.ReadFile(logPath())
		if err == nil {
			json.Unmarshal(logData, &logEntries)
		}
	}

	entry := LogEntry{
		ArtifactHash: artifact.ArtifactID,
	}

	logEntries = append(logEntries, entry)

	out, err := json.MarshalIndent(logEntries, "", "  ")
	if err != nil {
		fmt.Println("Error encoding log:", err)
		return
	}

	err = os.WriteFile(logPath(), out, 0644)
	if err != nil {
		fmt.Println("Error writing log:", err)
		return
	}

	fmt.Println("Record appended to registry")
	fmt.Println("Stored as:", target)
	fmt.Println("Artifact ID:", artifact.ArtifactID)
}

func verifyRegistry() {

	data, err := os.ReadFile(logPath())
	if err != nil {
		fmt.Println("Error reading registry log:", err)
		return
	}

	var entries []LogEntry

	err = json.Unmarshal(data, &entries)
	if err != nil {
		fmt.Println("Invalid registry log")
		return
	}

	fmt.Println("Registry entries:")

	for _, entry := range entries {
		fmt.Println(entry.ArtifactHash)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: registry [append|verify] <artifact>")
		return
	}

	cmd := os.Args[1]

	switch cmd {

	case "append":

		if len(os.Args) < 3 {
			fmt.Println("Usage: registry append <artifact.json>")
			return
		}

		appendArtifact(os.Args[2])

	case "verify":

		verifyRegistry()

	default:

		fmt.Println("Unknown command")
	}
}