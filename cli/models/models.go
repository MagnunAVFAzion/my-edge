package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

const OUT_DIR string = "./out"
const MANIFEST_PATH string = OUT_DIR + "/manifest.json"

type Project struct {
	Id             uuid.UUID `json:"id"`
	ClientId       string    `json:"clientId"`
	Name           string    `json:"name"`
	CurrentVersion string    `json:"version"`
}

type Manifest struct {
	Project Project           `json:"project"`
	Routes  map[string]string `json:"routes"`
	Cache   map[string]string `json:"cache"`
}

func (m *Manifest) Generate() error {
	fmt.Println("\n* Generating manifest ...")

	jsonData, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	err = os.WriteFile(MANIFEST_PATH, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	fmt.Println("Manifest created!")

	return nil
}
