package models

import (
	"encoding/json"
	"fmt"
	"os"
)

const OUT_DIR string = "./out"
const MANIFEST_PATH string = OUT_DIR + "/manifest.json"

type Project struct {
	Id             string `json:"id"`
	ClientId       string `json:"clientId"`
	Name           string `json:"name"`
	CurrentVersion string `json:"version"`
}

type Manifest struct {
	Project Project           `json:"project"`
	Routes  map[string]string `json:"routes"`
	Cache   map[string]string `json:"cache"`
}

func (m *Manifest) Save() error {
	fmt.Println("\n* Saving manifest ...")

	jsonData, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	err = os.WriteFile(MANIFEST_PATH, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	fmt.Println("Manifest saved!")

	return nil
}

func (m *Manifest) GetProjectStateIdentifier() string {
	return fmt.Sprintf("%s-%s-%s", m.Project.ClientId, m.Project.Id, m.Project.CurrentVersion)
}

func (m *Manifest) ReadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file content: %w", err)
	}

	err = json.Unmarshal(content, m)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %w", err)
	}

	return nil
}

// must ignore other attributes
type VercelConfig struct {
	Runtime string `json:"runtime"`
	Name    string `json:"name"`
}

func (v *VercelConfig) ReadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file content: %w", err)
	}

	err = json.Unmarshal(content, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %w", err)
	}

	return nil
}
