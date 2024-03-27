package sops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Config represents the configuration used to create a new Sops.
type Sops struct {
	AgeKey   string
	SopsFile string
}

// New creates a new configured Sops.
func New(key string, sops string) (*Sops, error) {
	return &Sops{
		AgeKey:   key,
		SopsFile: sops,
	}, nil
}

// Decrypt decrypts the given data.
func (s *Sops) decrypt(data string, path string) (string, error) {
	log.Debug().Msg("Decrypting file")
	// create any needed parent directories
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directories: %w", err)
	}
	// create a temp file
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	// write the data to the temp file
	if _, err := f.WriteString(data); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	cmd := exec.Command("sops", "--decrypt", "--input-type", "yaml", "--output-type", "yaml", f.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file %s: %s\n%w", f.Name(), out, err)
	}
	// delete the temp file
	if err := os.Remove(f.Name()); err != nil {
		return "", fmt.Errorf("failed to remove temp file: %w", err)
	}
	return string(out), nil
}

// Encrypt encrypts the given data.
func (s *Sops) Encrypt(file string) (string, error) {
	return "", nil
}

func (s *Sops) DecryptDir(data map[string]string) (map[string]string, error) {
	log.Debug().Msg("Decrypting directory data")

	// create temp directory
	dir, err := os.MkdirTemp("", "sops")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	// write sops file to .sops.yaml in the temp directory
	sopsFile := fmt.Sprintf("%s/.sops.yaml", dir)
	f, err := os.Create(sopsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create sops file: %w", err)
	}
	if _, err := f.WriteString(s.SopsFile); err != nil {
		return nil, fmt.Errorf("failed to write to sops file: %w", err)
	}

	for k, v := range data {
		if isEncrypted(v) {
			decrypted, err := s.decrypt(v, fmt.Sprintf("%s/%s", dir, k))
			if err != nil {
				return nil, err
			}
			data[k] = decrypted
		}
	}

	// delete the temp directory
	if err := os.RemoveAll(dir); err != nil {
		return nil, fmt.Errorf("failed to remove temp directory: %w", err)
	}
	return data, nil
}

func isEncrypted(data string) bool {
	return strings.Contains(data, "-----BEGIN AGE ENCRYPTED FILE----")
}
