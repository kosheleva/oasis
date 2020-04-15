package openapi3

import (
	"net/url"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestCreatePath(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"
	expectedPath := "/foo/p1_from_op/bar/p2_from_path/p3_from_override"

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "param2",
					Example: "p2_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param3",
					Required: true,
					Example:  "p3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_path",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param2",
					Required: true,
					Example:  "p2_from_path",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{
		Path: api.PathParameters{
			"param1": "",
			"param3": "p3_from_override",
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualPath, err := spec.CreatePath(path, oasPathItem, &oasOperation, &params)

	assert.Equal(T, actualPath, expectedPath)
	assert.Nil(T, err)
}

func TestCreatePath_ErrNoParameters(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}/qeq/{param4}"
	expectedPath := "/foo/p1_from_op/bar/{param2}/{param3}/qeq/{param4}"
	expectedErr := errors.NoParameters([]string{"param2", "param3", "param4"}, nil)

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "param2",
					Example: "p2_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualPath, err := spec.CreatePath(path, oasPathItem, &oasOperation, &params)

	assert.Equal(T, actualPath, expectedPath)
	assert.Equal(T, err.Error(), expectedErr.Error())
}

func TestCreateQuery(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"
	expectedQuery := &url.Values{
		"q1": []string{
			"q1_from_override",
			"q1_from_op",
			"q1_from_path",
		},

		"q2": []string{
			"q2_from_op",
		},

		"q3": []string{
			"q3_from_override",
			"q3_from_path",
		},
	}

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q2",
					Example:  "q2_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "q3",
					Example: "q3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_path",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q3",
					Example:  "q3_from_path",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{
		Query: url.Values{
			"q1": []string{
				"q1_from_override",
			},
			"q2": nil,
			"q3": []string{
				"q3_from_override",
			},
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualQuery, err := spec.CreateQuery(oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedQuery, actualQuery)
	assert.Nil(T, err)
}

func TestCreateQuery_ErrNoParameters(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"

	expectedQuery := &url.Values{
		"q1": []string{
			"q1_from_op",
		},
	}

	expectedErr := errors.NoParameters([]string{"q2", "q3", "q4"}, nil)

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
					Example:  "q1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q2",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "q3",
					Example: "q3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q1",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q3",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "query",
					Required: true,
					Name:     "q4",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualQuery, actualErr := spec.CreateQuery(oasPathItem, &oasOperation, &params)

	assert.Equal(T, expectedQuery, actualQuery)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestGetHost(T *testing.T) {
	expectedHost := &api.Host{
		URL:  "stage.localhost",
		Name: "Staging localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetHost("Staging localhost")

	assert.Equal(T, expectedHost, actualHost)
	assert.Nil(T, actualErr)
}

func TestGetHost_HostNotFound(T *testing.T) {
	expectedErr := errors.ErrHostNotFound{
		HostName: "Production localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetHost("Production localhost")

	assert.Nil(T, actualHost)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}

func TestGetDefaultHost(T *testing.T) {
	expectedHost := &api.Host{
		URL:  "localhost",
		Name: "Localhost",
	}

	OAS := &openapi3.Swagger{}
	OAS.Servers = append(OAS.Servers, &openapi3.Server{
		URL:         "localhost",
		Description: "Localhost",
	}, &openapi3.Server{
		URL:         "stage.localhost",
		Description: "Staging localhost",
	}, &openapi3.Server{
		URL:         "dev.localhost",
		Description: "Development localhost",
	})

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetDefaultHost()

	assert.Equal(T, expectedHost, actualHost)
	assert.Nil(T, actualErr)
}

func TestGetDefaultHost_HostNotFound(T *testing.T) {
	expectedErr := errors.ErrHostNotFound{
		HostName: "Default",
	}

	OAS := &openapi3.Swagger{}
	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualHost, actualErr := spec.GetDefaultHost()

	assert.Nil(T, actualHost)
	assert.Equal(T, expectedErr.Error(), actualErr.Error())
}