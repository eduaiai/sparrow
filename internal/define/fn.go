package define

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	UUID "github.com/google/uuid"
)

func GetBool(envKey string, def bool) bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env == "" {
		return def
	}
	return env == "on" || env == "true" || env == "1"
}

func GetPortString(envKey string, def int) string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env != "" {
		var REGEXP_PURE_NUMBER = regexp.MustCompile(`^\d{4,5}$`)
		appPortMatchFromEnv := REGEXP_PURE_NUMBER.FindAllString(env, -1)
		if len(appPortMatchFromEnv) == 1 {
			return fmt.Sprintf(":%s", appPortMatchFromEnv[0])
		}
	}
	return fmt.Sprintf(":%d", def)
}

func GetHostName(envKey string, def string) string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env == "" {
		return def
	}

	u, err := url.Parse(env)
	if !(err == nil && u.Scheme != "" && u.Host != "") {
		return def
	}

	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "ws" && u.Scheme != "wss" {
		return def
	}

	return strings.TrimSuffix(fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path), "/")
}

func GenerateRandomString(length int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var result string
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		result += string(letters[n.Int64()])
	}
	return result
}

func GenerateUUID() string {
	return UUID.New().String()
}

func GetGenerateSpeed(envKey string, def int) int {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env != "" {
		var REGEXP_PURE_NUMBER = regexp.MustCompile(`^\d{1,2}$`)
		appMultipleMatchFromEnv := REGEXP_PURE_NUMBER.FindAllString(env, -1)
		if len(appMultipleMatchFromEnv) == 1 {
			num, err := strconv.Atoi(appMultipleMatchFromEnv[0])
			if err != nil {
				return def
			}
			return num
		}
	}
	return def
}

func GetSecret(envKey string, def string) string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env != "" {
		return env
	}
	return def
}

func GetPath(envKey string, def string) string {
	env := strings.ToLower(strings.TrimSpace(os.Getenv(envKey)))
	if env != "" {
		return env
	}
	return def
}

func GetRandomNumber(min, max int) int {
	const fallback = 50
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return fallback
	}

	num := n.Int64()
	if num > int64(^uint(0)>>1) {
		return fallback
	}
	return int(num) + min
}

func MakeJSON(data any) (string, error) {
	ret, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(ret), nil
}
