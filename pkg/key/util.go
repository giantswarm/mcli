package key

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"go.uber.org/config"
	"gopkg.in/yaml.v3"
)

func GetSecretValue(key string, data string) (string, error) {
	value, err := GetValue(key, data)
	if err != nil {
		return "", fmt.Errorf("failed to get %s.\n%w", key, err)
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Debug().Msgf("failed to decode %s. Value does not seem to be base64 encoded. Returning raw value.", key)
		return value, nil
	}
	return string(decoded), nil
}

func GetValue(key string, data string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`%s: (\S+)`, key))
	matches := re.FindStringSubmatch(data)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to find %s in file", key)
	}
	if strings.HasPrefix(matches[1], "|") {
		return GetMultiLineValue(key, data)
	}
	return matches[1], nil
}

func GetMultiLineValue(key string, data string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`(?s)%s:(.+)`, key))
	matches := re.FindStringSubmatch(data)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find %s in file %s but found %v matches %s", key, data, len(matches), matches)
	}
	value := matches[1]
	// remove any other keys that might be in the value
	re = regexp.MustCompile(`\n\S+: \S+`)
	matches = re.FindStringSubmatch(value)
	if len(matches) > 0 {
		value = strings.Split(value, matches[0])[0]
	}

	// remove unneeded whitespaces
	value = strings.TrimPrefix(value, " ")
	value = strings.TrimPrefix(value, "|-")
	value = strings.TrimPrefix(value, "|")
	re = regexp.MustCompile(`\n\s+`)
	value = re.ReplaceAllString(value, "\n")
	value = strings.TrimPrefix(value, "\n")
	value = strings.TrimSuffix(value, "\n")

	return value, nil
}

func GetNamespacedName(data string) (name, namespace string) {
	metadata, err := GetMultiLineValue("metadata", data)
	if err != nil {
		return "", ""
	}
	name, err = GetValue("name", metadata)
	if err != nil {
		return "", ""
	}
	namespace, err = GetValue("namespace", metadata)
	if err != nil {
		return "", ""
	}
	return name, namespace
}

func GetData(data any) ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func MergeValues(valuesA string, valuesB string) (string, error) {
	base := strings.NewReader(valuesA)
	override := strings.NewReader(valuesB)

	provider, err := config.NewYAML(config.Source(base), config.Source(override))
	if err != nil {
		return "", fmt.Errorf("failed to create yaml provider.\n%w", err)
	}

	var c map[string]interface{}
	err = provider.Get(config.Root).Populate(&c)
	if err != nil {
		return "", fmt.Errorf("failed to populate config.\n%w", err)
	}

	data, err := GetData(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged values.\n%w", err)
	}

	return string(data), nil
}
