//go:build integration
// +build integration

// With build tags you can tell the Go compiler to only include files when the tag is present
// It was originally intended for different platform implementations but can also be used to separate tests
// With this you can call your tests with `$ go test -tags integration ...` to include the integration tests

package solver

import (
	"context"
	"net/http"
	"testing"
)

func TestRemoteSolver_ResolveIntegration(t *testing.T) {
	// Since this is an integration test we actually use a local running server to test
	rs := RemoteSolver{
		MathServerURL: "http://localhost:8080",
		Client:        http.DefaultClient,
	}
	data := []struct {
		name       string
		expression string
		result     float64
		errMsg     string
	}{
		{"case1", "2 + 2 * 10", 22, ""},
		{"case2", "( 2 + 2 ) * 10", 40, ""},
		{"case3", "( 2 + 2 * 10", 0, "invalid expression: ( 2 + 2 * 10"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := rs.Resolve(context.Background(), d.expression)
			if result != d.result {
				t.Errorf("expected `%f`, got `%f`", d.result, result)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.errMsg {
				t.Errorf("expected error `%s`, got `%s`", d.errMsg, errMsg)
			}
		})
	}
}
