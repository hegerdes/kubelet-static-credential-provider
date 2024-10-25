// internal/utils/utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kcpv1 "k8s.io/kubelet/pkg/apis/credentialprovider/v1"
)

type Config struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	CacheType string `yaml:"cache_type"`
}

// type CredentialProviderResponse struct {
// 	kcpv1.CredentialProviderResponse
// 	APIVersion string `json:"apiVersion"`
// 	Kind       string `json:"kind"`
// }

func GetConfig(path string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing YAML file:", err)
		return config, err

	}
	return config, nil
}

// Parse parses a JSON string into a CredentialRequest struct.
func GetRequestImage(credrequest string) (string, error) {
	var cred kcpv1.CredentialProviderRequest
	err := json.Unmarshal([]byte(credrequest), &cred)
	if err != nil {
		return "", err
	}

	// Compile the regex that matches ':' only if it is preceded by '/'
	// Replace the matched part to remove the characters after ':'
	re := regexp.MustCompile(`(\/[^:]*)\:.*`)
	baseImage := re.ReplaceAllString(cred.Image, `$1`)

	return baseImage, nil
}

var cacheTypeMap = map[string]kcpv1.PluginCacheKeyType{
	"registry": kcpv1.RegistryPluginCacheKeyType,
	"image":    kcpv1.ImagePluginCacheKeyType,
	"global":   kcpv1.GlobalPluginCacheKeyType,
}

func CreateImageRequestResponse(image, username, password, cachetype string) string {

	// Set cache type
	var usedCacheType = kcpv1.ImagePluginCacheKeyType
	if value, ok := cacheTypeMap[strings.ToLower(cachetype)]; ok {
		usedCacheType = value
	}

	// Set cache duration
	duration, err := time.ParseDuration("8h")
	if err != nil {
		return ""
	}
	k8sDuration := metav1.Duration{Duration: duration}

	cred := kcpv1.CredentialProviderResponse{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "credentialprovider.kubelet.k8s.io/v1",
			Kind:       "CredentialProviderResponse"},
		CacheKeyType:  usedCacheType,
		CacheDuration: &k8sDuration,
		Auth: map[string]kcpv1.AuthConfig{
			image: {
				Username: username,
				Password: password,
			},
		}}
	resp, err := json.Marshal(cred)
	if err != nil {
		return ""
	}

	return string(resp)
}
