package details

import (
	json "github.com/json-iterator/go"
	"testing"
)

func TestSource(t *testing.T) {
	result, _ := Mangakakalot("Dungeon Reset")

	marshal, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return
	}

	println(string(marshal))
}
