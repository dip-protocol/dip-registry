package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Record struct {
	ID       string                 `json:"decision_id"`
	Data     map[string]interface{} `json:"data"`
	PrevHash string                 `json:"prev_hash"`
	Hash     string                 `json:"hash"`
}

func hashRecord(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func getLastHash() string {

	file, err := os.Open("registry.log")
	if err != nil {
		return ""
	}

	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if lastLine == "" {
		return ""
	}

	var rec Record
	json.Unmarshal([]byte(lastLine), &rec)

	return rec.Hash
}

func appendRecord(record Record) error {

	record.PrevHash = getLastHash()

	dataBytes, err := json.Marshal(record.Data)
	if err != nil {
		return err
	}

	hashInput := append(dataBytes, []byte(record.PrevHash)...)
	record.Hash = hashRecord(hashInput)

	file, err := os.OpenFile("registry.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = file.Write(append(bytes, '\n'))
	return err
}

func verifyChain() error {

	file, err := os.Open("registry.log")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var prevHash string

	for scanner.Scan() {

		var rec Record
		err := json.Unmarshal(scanner.Bytes(), &rec)
		if err != nil {
			return err
		}

		if rec.PrevHash != prevHash {
			return fmt.Errorf("chain broken at record %s", rec.ID)
		}

		prevHash = rec.Hash
	}

	fmt.Println("Registry chain verified")
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run registry.go [append|verify]")
		return
	}

	switch os.Args[1] {

	case "append":

		rec := Record{
			ID: "example-001",
			Data: map[string]interface{}{
				"status": "approved",
			},
		}

		err := appendRecord(rec)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Record appended to registry")

	case "verify":

		err := verifyChain()
		if err != nil {
			fmt.Println("Registry verification failed:", err)
			return
		}
	}
}