package client

import (
	"os"
	"testing"
)

func TestParseRawData(t *testing.T) {
	f, err := os.Open("test_data.json")
	if err != nil {
		t.Errorf("Unable to open test file: %+v", err)
	}

	_, err = newParser().parseRawData(f)
	if err != nil {
		t.Errorf("Unable to parse test file: %+v", err)
	}
}

func TestParseRawData2(t *testing.T) {
	f, err := os.Open("test_data2.json")
	if err != nil {
		t.Errorf("Unable to open test file: %+v", err)
	}

	_, err = newParser().parseRawData(f)
	if err != nil {
		t.Errorf("Unable to parse test file: %+v", err)
	}
}

func TestParseRawData3(t *testing.T) {
	f, err := os.Open("test_data3.json")
	if err != nil {
		t.Errorf("Unable to open test file: %+v", err)
	}

	_, err = newParser().parseRawData(f)
	if err != nil {
		t.Errorf("Unable to parse test file: %+v", err)
	}
}
