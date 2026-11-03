package truto

import (
	"encoding/json"
	"os"
	"testing"
)

// testdata/qs_cases.json was generated with qs.stringify(input, {sort}) from
// the qs package the TypeScript SDK and the Truto API use.
func TestEncodeQueryMatchesQS(t *testing.T) {
	raw, err := os.ReadFile("testdata/qs_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name     string `json:"name"`
		Input    Params `json:"input"`
		Expected string `json:"expected"`
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no cases")
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := encodeQuery(c.Input)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.Expected {
				t.Errorf("\n got: %s\nwant: %s", got, c.Expected)
			}
		})
	}
}

func TestEncodeQueryGoValues(t *testing.T) {
	type filter struct {
		Status string `json:"status"`
	}
	got, err := encodeQuery(Params{
		"ids":    []string{"a", "b"},
		"limit":  25,
		"filter": filter{Status: "open"},
		"flag":   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "filter%5Bstatus%5D=open&flag=true&ids%5B0%5D=a&ids%5B1%5D=b&limit=25"
	if got != want {
		t.Errorf("\n got: %s\nwant: %s", got, want)
	}
}
