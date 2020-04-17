package visualtest

import (
	"io/ioutil"
	"os"
	"testing"
)

func GoldenValue(t *testing.T, goldenFile string, actual string, update bool) string {
	t.Helper()
	goldenPath := "testdata/" + goldenFile + ".golden"
	f, err := os.Open(goldenPath)
	defer f.Close()
	if update {
		err := ioutil.WriteFile(goldenPath, []byte(actual), 0644)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", goldenPath, err)
		}
		return actual
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", goldenPath, err)
	}
	return string(content)
}
