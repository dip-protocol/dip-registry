package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Record struct {
	ID   string                 `json:"decision_id"`
	Data map[string]interface{} `json:"data"`
}

func appendRecord(record Record) error {

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

func main() {

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
}
