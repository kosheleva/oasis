package openapi3

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Spec is an OAS3-backed API test spec.
type Spec struct {
	Log log.ILogger
	OAS *openapi3.Swagger
}

// GetOperations returns a list of all available test operations from the spec.
func (spec *Spec) GetOperations(params *api.OperationParameters) []*api.Operation {
	ops := []*api.Operation{}

	addOp := func(oasOp *openapi3.Operation, method string, oasPath string, oasPathItem *openapi3.PathItem) {
		if oasOp != nil {
			ops = append(ops, spec.makeOperation(method, oasOp, oasPath, oasPathItem, params))
		}
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
		addOp(oasPathItem.Get, "GET", oasPath, oasPathItem)
		addOp(oasPathItem.Post, "POST", oasPath, oasPathItem)
		addOp(oasPathItem.Put, "PUT", oasPath, oasPathItem)
		addOp(oasPathItem.Delete, "DELETE", oasPath, oasPathItem)
		addOp(oasPathItem.Patch, "PATCH", oasPath, oasPathItem)
		addOp(oasPathItem.Head, "HEAD", oasPath, oasPathItem)
		addOp(oasPathItem.Options, "OPTIONS", oasPath, oasPathItem)
		addOp(oasPathItem.Trace, "TRACE", oasPath, oasPathItem)
		addOp(oasPathItem.Connect, "CONNECT", oasPath, oasPathItem)
	}

	return ops
}

// GetOperation returns a list of all available test operations from the spec.
func (spec *Spec) GetOperation(name string, params *api.OperationParameters) *api.Operation {
	filterOp := func(oasOp *openapi3.Operation) bool {
		return (oasOp.Summary == name || oasOp.OperationID == name)
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
		if filterOp(oasPathItem.Get) {
			return spec.makeOperation("GET", oasPathItem.Get, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Post) {
			return spec.makeOperation("POST", oasPathItem.Post, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Put) {
			return spec.makeOperation("PUT", oasPathItem.Put, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Delete) {
			return spec.makeOperation("DELETE", oasPathItem.Delete, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Patch) {
			return spec.makeOperation("PATCH", oasPathItem.Patch, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Head) {
			return spec.makeOperation("HEAD", oasPathItem.Head, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Options) {
			return spec.makeOperation("OPTIONS", oasPathItem.Options, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Connect) {
			return spec.makeOperation("CONNECT", oasPathItem.Connect, oasPath, oasPathItem, params)
		}
		if filterOp(oasPathItem.Trace) {
			return spec.makeOperation("TRACE", oasPathItem.Trace, oasPath, oasPathItem, params)
		}
	}

	return nil
}

func (spec *Spec) makeOperation(
	method string,
	oasOp *openapi3.Operation,
	oasPath string,
	oasPathItem *openapi3.PathItem,
	params *api.OperationParameters,
) *api.Operation {
	specPath := spec.CreatePath(oasPath, oasPathItem, oasOp, params)
	specRequests := []*api.Request{}

	if oasOp.RequestBody != nil && oasOp.RequestBody.Value != nil {
		for oasCT, oasMT := range oasOp.RequestBody.Value.Content {
			specRequests = append(specRequests, spec.MakeRequest(method, specPath, oasOp, oasPathItem, oasCT, oasMT, params))
		}
	} else {
		specRequests = append(specRequests, spec.MakeRequest(method, specPath, oasOp, oasPathItem, "", nil, params))
	}

	return &api.Operation{
		ID:          oasOp.OperationID,
		Name:        oasOp.Summary,
		Description: oasOp.Description,
		Method:      method,
		Path:        specPath,
		Requests:    specRequests,
	}
}

// MakeRequest creates an api.Rrequest instance from available operation spec data.
func (spec *Spec) MakeRequest(
	method string,
	specPath string,
	oasOp *openapi3.Operation,
	oasPathItem *openapi3.PathItem,
	oasCT string,
	oasMT *openapi3.MediaType,
	params *api.OperationParameters,
) *api.Request {

	specReq := &api.Request{}
	specReq.Method = method
	specReq.Path = specPath

	return specReq
}

// CreatePath creates an operation path with parameters replaced by actual values from `example`.
// Examples from operation-level parameters override examples from the path-level ones.
func (spec *Spec) CreatePath(
	oasPath string,
	oasPathItem *openapi3.PathItem,
	oasOp *openapi3.Operation,
	params *api.OperationParameters,
) string {
	path := oasPath

	fixPath := func(ppn string, ppv string, container string) {
		RX, _ := regexp.Compile("\\{" + ppn + "\\}")

		if RX.Match([]byte(path)) {
			if ppv != "" {
				path = string(RX.ReplaceAll([]byte(path), []byte(ppv)))
				spec.Log.UsingParameterExample(ppn, "path", container)
			} else {
				spec.Log.ParameterHasNoExample(ppn, "path", container)
			}
		}
	}

	useParameters := func(specParams openapi3.Parameters, container string) {
		for _, specP := range specParams {
			if specP == nil || specP.Value == nil || specP.Value.In != "path" {
				continue
			}

			fixPath(specP.Value.Name, specP.Value.Example.(string), container)
		}
	}

	for ppn, ppv := range params.Path {
		fixPath(ppn, ppv, "override")
	}

	useParameters(oasOp.Parameters, "operation")
	useParameters(oasPathItem.Parameters, "path")

	return path
}

// GetProjectInfo returns project info from the spec.info object.
func (spec *Spec) GetProjectInfo() *api.ProjectInfo {
	return &api.ProjectInfo{
		Title:       spec.OAS.Info.Title,
		Description: spec.OAS.Info.Description,
		Version:     spec.OAS.Info.Version,
	}
}

// GetHost returns an API host by requested description
// from the spec.servers list.
func (spec *Spec) GetHost(name string) *api.Host {
	for _, oasServer := range spec.OAS.Servers {
		if oasServer.Description == name {
			return &api.Host{
				Name:        "Default",
				Description: oasServer.Description,
				URL:         oasServer.URL,
			}
		}
	}
	return nil
}

// GetDefaultHost returns the fisr host from the spec.servers list as default.
func (spec *Spec) GetDefaultHost() *api.Host {
	if len(spec.OAS.Servers) > 0 {
		return &api.Host{
			Name:        "Default",
			Description: spec.OAS.Servers[0].Description,
			URL:         spec.OAS.Servers[0].URL,
		}
	}
	return nil
}
