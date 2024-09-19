package utils

import (
	"testing"
)

func TestUUID(t *testing.T) {
	id, err := UUID()
	if err != nil {
		t.Errorf("Failed to generate UUID: %v", err)
	}
	if len(id) != 36 {
		t.Errorf("UUID length is not 36: %s", id)
	}
}

func TestJSONMarshal(t *testing.T) {
	type TestStruct struct {
		Foo string
		Bar int
	}
	s := TestStruct{"foo", 42}
	json, err := JSONMarshal(s)
	if err != nil {
		t.Errorf("Failed to marshal JSON: %v", err)
	}
	expected := `{"Foo":"foo","Bar":42}`
	if string(json) != expected {
		t.Errorf("JSON output does not match expected: %s != %s", json, expected)
	}
}

func TestHTTPGet(t *testing.T) {
	url := "https://example.com"
	resp, err := HTTPGet(url)
	if err != nil {
		t.Errorf("Failed to make GET request: %v", err)
	}
	if len(resp) == 0 {
		t.Errorf("Response body is empty")
	}
}

func TestURLParse(t *testing.T) {
	urlStr := "https://example.com/path?a=1&b=2"
	u, err := URLParse(urlStr)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
	}
	if u.Scheme != "https" {
		t.Errorf("URL scheme is not https: %s", u.Scheme)
	}
}

func TestFilePathJoin(t *testing.T) {
	paths := []string{"path", "to", "file"}
	expected := "path/to/file"
	actual := FilePathJoin(paths...)
	if actual != expected {
		t.Errorf("File path join does not match expected: %s != %s", actual, expected)
	}
}

func TestIsNil(t *testing.T) {
	var v interface{}
	if !IsNil(v) {
		t.Errorf("Expected nil value to be nil")
	}
	v = struct{}{}
	if IsNil(v) {
		t.Errorf("Expected non-nil value to not be nil")
	}
}

func TestIsEmpty(t *testing.T) {
	var s string
	if !IsEmpty(s) {
		t.Errorf("Expected empty string to be empty")
	}
	s = "not empty"
	if IsEmpty(s) {
		t.Errorf("Expected non-empty string to not be empty")
	}
}

func TestRandInt(t *testing.T) {
	min, max := 1, 10
	n := RandInt(min, max)
	if n < min || n > max {
		t.Errorf("Random int is out of range: %d", n)
	}
}

func TestRandString(t *testing.T) {
	length := 10
	s := RandString(length)
	if len(s) != length {
		t.Errorf("Random string length is not %d: %s", length, s)
	}
}

func TestTimeNow(t *testing.T) {
	now := TimeNow()
	if now.IsZero() {
		t.Errorf("Current time is zero")
	}
}

func TestTimeParse(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	value := "2022-07-25 14:30:00"
	t, err := TimeParse(layout, value)
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	if t.Format(layout) != value {
		t.Errorf("Parsed time does not match expected: %s != %s", t.Format(layout), value)
	}
}

func TestTimeFormat(t *testing.T) {
	t := TimeNow()
	layout := "2006-01-02 15:04:05"
	formatted := TimeFormat(t, layout)
	if formatted != t.Format(layout) {
		t.Errorf("Formatted time does not match expected: %s != %s", formatted, t.Format(layout))
	}
}

func TestRegexpMatch(t *testing.T) {
	pattern := "hello.*"
	str := "hello world"
	if !RegexpMatch(pattern, str) {
		t.Errorf("String does not match pattern: %s", str)
	}
}

func TestOS(t *testing.T) {
	os := OS()
	if os == "" {
		t.Errorf("Operating system is empty")
	}
}

func TestArch(t *testing.T) {
	arch := Arch()
	if arch == "" {
		t.Errorf("Architecture is empty")
	}
}
