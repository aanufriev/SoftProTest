package workerspool

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	sport := "soccer"
	json := `{"lines":{"SOCCER":"1.667"}}`

	expected := "1.667"
	res, _ := parseLine([]byte(json), sport)

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Test failed. Results not match\nGot:\n%v\nExpected:\n%v", res, expected)
	}
}
