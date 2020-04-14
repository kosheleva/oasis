package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/errors"
)

// Spec is an interface to access specification data.
type Spec interface {
	GetProjectInfo() *ProjectInfo
	GetHost(name string) *Host
	GetDefaultHost() *Host
	GetOperations(params *OperationParameters) []*Operation
	GetOperation(name string, params *OperationParameters) (*Operation, errors.IError)
}

// ProjectInfo is a generic project information.
type ProjectInfo struct {
	Title       string
	Description string
	Version     string
}

// ExampleList is a map of maps to keep request example data in.
type ExampleList map[string]ExampleObject

// ExampleObject --
type ExampleObject map[interface{}]interface{}

// MarshalJSON encodes an example map from the OAS spec as a JSON string.
func (ex ExampleObject) MarshalJSON() ([]byte, error) {
	props := []string{}

	for propKey, propVal := range ex {
		jp, err := json.Marshal(propVal)
		if err != nil {
			return nil, err
		}

		props = append(props, fmt.Sprintf("\"%s\":%s", propKey, jp))
	}

	return []byte(fmt.Sprintf("{%s}", strings.Join(props, ","))), nil
}

// Host is an API host description.
type Host struct {
	Name        string
	URL         string
	Description string
}
