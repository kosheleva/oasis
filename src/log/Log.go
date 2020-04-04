package log

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/xeipuuv/gojsonschema"
)

// ILogger - interface for execution loggers
type ILogger interface {
	Usage()
	Error(err error)
	LoadingSpec(path string)

	PrintOperations(ops []*api.Operation)
	TestingProject(p *api.ProjectInfo)
	TestingOperation(res *api.Operation)

	UsingHost(p *api.Host)
	UsingDefaultHost()
	HostNotFound(h string)

	// UsingSecurity(sec *api.Security)
	// UsingRequest(req *api.Request)
	UsingResponse(req *api.Response)

	Overriding(what string)
	Requesting(url string)

	ParameterHasNoExample(paramName string, in string, container string)
	UsingParameterExample(paramName string, in string, container string)

	// PropertyHasNoValue(prop *api.Property, ctx *utility.Context)
	HeaderHasNoValue(header *api.Header)
	ResponseHasWrongStatus(schema *api.Response, actualStatus int)
	ResponseHasWrongContentType(schema *api.Response, actualCT string)
	ResponseNotFound(CT string, status int)

	OperationOK(res *api.Operation)
	OperationFail(res *api.Operation)
	OperationNotFound(op string)

	SchemaTesting(schema *api.Schema, data interface{})
	SchemaOK(schema *api.Schema)
	SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError)
	// UnknownSchemaDataType(schema *api.Schema)
	// SchemaExpectedBoolean(schema *api.Schema, v interface{})
	// SchemaExpectedNumber(schema *api.Schema, v interface{})
	// SchemaExpectedInteger(schema *api.Schema, v interface{})
	// SchemaExpectedString(schema *api.Schema, v interface{})
	// SchemaExpectedArray(schema *api.Schema, v interface{})
	// SchemaExpectedObject(schema *api.Schema, v interface{})
}

// Log is a base type for loggers.
type Log struct {
	Level int64
}

// Print prints.
func (log Log) Print(l int64, msg string, args ...interface{}) {
	if l <= log.Level {
		fmt.Printf(msg, args...)
	}
}

// Println prints and adds a newline.
func (log Log) Println(l int64, msg string, args ...interface{}) {
	log.Print(l, msg+"\n", args...)
}
