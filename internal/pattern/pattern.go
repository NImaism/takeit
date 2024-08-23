package pattern

import (
	"encoding/json"
	"fmt"
	"github.com/nimaism/takeit/internal/model"
	"net/http"
	"os"
	"path/filepath"
)

const (
	patternStorage = "https://raw.githubusercontent.com/EdOverflow/can-i-take-over-xyz/master/fingerprints.json"
	fileName       = "pattern.json"
)

func GetPatternPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "takeit")

	if err = os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating config directory '%s': %w", appConfigDir, err)
	}

	return filepath.Join(appConfigDir, fileName), nil
}

func UpdatePatterns(exclude []string) error {
	resp, err := http.Get(patternStorage)
	if err != nil {
		return fmt.Errorf("failed to download pattern: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var patterns []model.Pattern
	if err = json.NewDecoder(resp.Body).Decode(&patterns); err != nil {
		return fmt.Errorf("failed to decode json: %v", err)
	}

	var validPatterns []model.Pattern
	for _, v := range patterns {
		if v.Status == "Vulnerable" && len(v.Fingerprint) > 0 && !IsPatternExclude(v.Service, &exclude) {
			validPatterns = append(validPatterns, v)
		}
	}

	patternPath, err := GetPatternPath()
	if err != nil {
		return fmt.Errorf("failed to get pattern path: %v", err)
	}

	patternFile, err := os.OpenFile(patternPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open pattern file: %v", err)
	}
	defer patternFile.Close()

	if err = json.NewEncoder(patternFile).Encode(validPatterns); err != nil {
		return fmt.Errorf("failed to write pattern to file: %v", err)
	}

	return nil
}

func LoadPatterns(exclude []string) (*[]model.Pattern, error) {
	var patterns []model.Pattern

	patternPath, err := GetPatternPath()
	if err != nil {
		return &patterns, fmt.Errorf("failed to get pattern path: %v", err)
	}

	file, err := os.ReadFile(patternPath)
	if err != nil {
		if err = UpdatePatterns(exclude); err != nil {
			return &patterns, fmt.Errorf("failed to update pattern file: %v", err)
		}

		file, err = os.ReadFile(patternPath)
		if err != nil {
			return &patterns, fmt.Errorf("failed to read pattern file: %v", err)
		}
	}

	if err = json.Unmarshal(file, &patterns); err != nil {
		return &patterns, fmt.Errorf("failed to unmarshal pattern file: %v", err)
	}

	return &patterns, nil
}

func IsPatternExclude(service string, list *[]string) bool {
	for _, v := range *list {
		if service == v {
			return true
		}
	}
	return false
}
