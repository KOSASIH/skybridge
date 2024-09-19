package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// UUID generates a random UUID
func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// JSONMarshal marshals a struct to JSON
func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONUnmarshal unmarshals JSON to a struct
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// HTTPGet makes a GET request to a URL
func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// URLParse parses a URL
func URLParse(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// FilePathJoin joins file paths
func FilePathJoin(paths ...string) string {
	return filepath.Join(paths...)
}

// IsNil checks if a value is nil
func IsNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

// IsEmpty checks if a string, slice, or map is empty
func IsEmpty(v interface{}) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return v.(string) == ""
	case reflect.Slice, reflect.Array:
		return reflect.ValueOf(v).Len() == 0
	case reflect.Map:
		return reflect.ValueOf(v).Len() == 0
	default:
		return false
	}
}

// RandInt generates a random integer
func RandInt(min, max int) int {
	if min >= max {
		return min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return min + int(n.Int64())
}

// RandString generates a random string of a given length
func RandString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// TimeNow returns the current time
func TimeNow() time.Time {
	return time.Now()
}

// TimeParse parses a time string
func TimeParse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// TimeFormat formats a time
func TimeFormat(t time.Time, layout string) string {
	return t.Format(layout)
}

// RegexpMatch checks if a string matches a regular expression
func RegexpMatch(pattern, str string) bool {
	re, _ := regexp.Compile(pattern)
	return re.MatchString(str)
}

// OS returns the current operating system
func OS() string {
	return runtime.GOOS
}

// Arch returns the current architecture
func Arch() string {
	return runtime.GOARCH
}
