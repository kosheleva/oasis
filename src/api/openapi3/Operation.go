package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Operation provides access to OAS3-specific API data.
type Operation struct {
	*contract.OperationPrototype

	RequestMethod string
	RequestPath   string
	SpecPath      *openapi3.PathItem
	SpecOp        *openapi3.Operation
}

// ID ...
func (op *Operation) ID() string {
	return op.SpecOp.OperationID
}

// Name ...
func (op *Operation) Name() string {
	return op.SpecOp.Summary
}

// Description ...
func (op *Operation) Description() string {
	return op.SpecOp.Description
}

// Method ...
func (op *Operation) Method() string {
	return op.RequestMethod
}

// Path ...
func (op *Operation) Path() string {
	return op.RequestPath
}
