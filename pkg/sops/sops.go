package sops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// this implementation relies on the sops binary being available
// it's not ideal but has been implemented for now since calling sops directly is how we are currently using it in mc-bootstrap
// the way sops is implemented did not allow importing the package directly in a way that would have been similar to the original command
// let's change this later

const EnvAgeKey = "SOPS_AGE_KEY"

// Decrypt decrypts the given data.
func decrypt(data string, path string) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Decrypting file %s", path))
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

	cmd := exec.Command("sops", "--decrypt", "--input-type", "yaml", "--output-type", "yaml", f.Name()) // #nosec G204
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
func encrypt(data string, path string, age string) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Encrypting file %s", path))
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

	cmd := exec.Command("sops", "--encrypt", "--input-type", "yaml", "--output-type", "yaml", "--age", age, "--encrypted-regex", "^(data|stringData)$", f.Name()) // #nosec G204
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file %s: %s\n%w", f.Name(), out, err)
	}
	// delete the temp file
	if err := os.Remove(f.Name()); err != nil {
		return "", fmt.Errorf("failed to remove temp file: %w", err)
	}
	return string(out), nil
}

func DecryptDir(input map[string]string) (map[string]string, error) {
	log.Debug().Msg("Decrypting directory data")
	data := map[string]string{}

	// create temp directory
	dir, err := os.MkdirTemp("", "sops")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	for k, v := range input {
		if IsEncrypted(v) {
			decrypted, err := decrypt(v, fmt.Sprintf("%s/%s", dir, k))
			if err != nil {
				return nil, err
			}
			data[k] = decrypted
		} else {
			data[k] = v
		}
	}

	// delete the temp directory
	if err := os.RemoveAll(dir); err != nil {
		return nil, fmt.Errorf("failed to remove temp directory: %w", err)
	}
	return data, nil
}

func EncryptDir(data map[string]string, age string) (map[string]string, error) {
	log.Debug().Msg("Encrypting directory data")

	// create temp directory
	dir, err := os.MkdirTemp("", "sops")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	for k, v := range data {
		encrypted, err := encrypt(v, fmt.Sprintf("%s/%s", dir, k), age)
		if err != nil {
			return nil, err
		}
		data[k] = encrypted
	}

	// delete the temp directory
	if err := os.RemoveAll(dir); err != nil {
		return nil, fmt.Errorf("failed to remove temp directory: %w", err)
	}
	return data, nil
}

func IsEncrypted(data string) bool {
	return strings.Contains(data, "-----BEGIN AGE ENCRYPTED FILE----")
}

func GetAgeKey() string {
	return os.Getenv(EnvAgeKey)
}
