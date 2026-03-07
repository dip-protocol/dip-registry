package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Artifact struct {
	ArtifactVersion string    `json:"artifact_version"`
	ArtifactID      string    `json:"artifact_id"`
	Decision        any       `json:"decision"`
	Signature       Signature `json:"signature"`
}

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey []byte `json:"public_key"`
	Value     []byte `json:"value"`
}

func canonicalizeJSON(data []byte) ([]byte, error) {

	var obj any

	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err = enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(buf.Bytes()), nil
}

func computeArtifactID(canonical []byte) string {

	hash := sha256.Sum256(canonical)

	return hex.EncodeToString(hash[:])
}

func appendToRegistry() {

	cmd := exec.Command("../dip-registry/registry.exe", "append", "artifact.json")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Registry append failed:", err)
		return
	}

	fmt.Println(string(output))
}

func signDecision(inputFile string) {

	raw, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading decision file:", err)
		return
	}

	canonical, err := canonicalizeJSON(raw)
	if err != nil {
		fmt.Println("Canonicalization error:", err)
		return
	}

	var decision any

	json.Unmarshal(canonical, &decision)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Key generation error:", err)
		return
	}

	artifactID := computeArtifactID(canonical)

	sig := ed25519.Sign(priv, canonical)

	artifact := Artifact{
		ArtifactVersion: "1.0",
		ArtifactID:      artifactID,
		Decision:        decision,
		Signature: Signature{
			Algorithm: "ed25519",
			PublicKey: pub,
			Value:     sig,
		},
	}

	out, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		fmt.Println("Artifact encoding error:", err)
		return
	}

	err = os.WriteFile("artifact.json", out, 0644)
	if err != nil {
		fmt.Println("Artifact write error:", err)
		return
	}

	fmt.Println("DIP artifact created: artifact.json")
	fmt.Println("Artifact ID:", artifactID)

	appendToRegistry()
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  dip sign <decision.json>")
		return
	}

	cmd := os.Args[1]

	if cmd == "sign" {

		signDecision(os.Args[2])

		return
	}

	fmt.Println("Unknown command")
}
