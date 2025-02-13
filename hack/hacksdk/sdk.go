package hacksdk

import (
	"fmt"
	"strings"

	"github.com/markkurossi/scheme"
)

// SDK represents the main interface for the Hack SDK
type SDK struct {
	scm *scheme.Scheme
}

// NewSDK creates a new instance of the SDK
func NewSDK() (*SDK, error) {
	scm, err := scheme.NewWithParams(scheme.Params{
		Verbose:   false,
		NoRuntime: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize scheme: %w", err)
	}

	return &SDK{scm: scm}, nil
}

// Eval evaluates a scheme expression and returns the result
func (s *SDK) Eval(schemeCode string) (scheme.Value, error) {
	payload := strings.NewReader(schemeCode)
	res, err := s.scm.Eval("", payload)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %w", err)
	}
	return res, nil
}
