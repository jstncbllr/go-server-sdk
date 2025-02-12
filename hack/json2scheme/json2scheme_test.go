package json2scheme

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/markkurossi/scheme"
	"github.com/stretchr/testify/require"
)

func TestJsonToScheme(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"segments":{},
		"flags":{
			"flag2":{
				"key":"flag2",
				"on":false,
				"prerequisites":[],
				"targets":[],
				"contextTargets":[],
				"rules":[],
				"fallthrough":{},
				"offVariation":null,
				"variations":[],
				"clientSide":false,
				"salt":"",
				"trackEvents":false,
				"trackEventsFallthrough":false,
				"debugEventsUntilDate":null,
				"version":0,
				"deleted":false
			}
		}
	}`

	// Convert to RawMessage
	rawMsg := json.RawMessage(jsonData)

	// Call jsonToScheme
	result, err := JsonToScheme(rawMsg)
	if err != nil {
		t.Fatalf("jsonToScheme failed: %v", err)
	}

	resultExpr := parse(t, result)

	// Create and parse the expected Scheme expression
	expected := `(
  (segments ())
  (flags (
    (flag2
      (key "flag2")
      (on #f)
      (prerequisites ())
      (targets ())
      (contextTargets ())
      (rules ())
      (fallthrough ())
      (offVariation nil)
      (variations ())
      (clientSide #f)
      (salt "")
      (trackEvents #f)
      (trackEventsFallthrough #f)
      (debugEventsUntilDate nil)
      (version 0)
      (deleted #f)))))`

	expectedExpr := parse(t, expected)
	if err != nil {
		t.Fatalf("Failed to parse expected as Scheme: %v", err)
	}

	// Compare the parsed Scheme expressions
	if !resultExpr.Equal(expectedExpr) {
		t.Errorf("Scheme expressions do not match\nexpected: %v\ngot: %v", expectedExpr, resultExpr)
	}
}

func parse(t *testing.T, s string) (result *scheme.Library) {
	// Parse the result into a Scheme expression
	scm, err := scheme.New()
	require.NoError(t, err)
	resultExpr, err := scheme.NewParser(scm).Parse("", strings.NewReader(s))
	if err != nil {
		t.Fatalf("Failed to parse result as Scheme: %v", err)
	}
	return resultExpr
}
