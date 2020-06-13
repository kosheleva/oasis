package apikey

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Security implements the security type that requires a unique API key
// to be present in each API request. The exact location of the key
// is implemented corresponding classes.
type Security struct {
	Name      string
	ParamName string
	Value     string
	Log       contract.Logger
}

// New creates a new API Key security scheme.
func New(name string, location string, paramName string, value string, logger contract.Logger) contract.Security {
	switch location {
	case "cookie":
		return Cookie{
			Security{name, paramName, value, logger},
		}

	case "header":
		return Header{
			Security{name, paramName, value, logger},
		}

	case "query":
		return Query{
			Security{name, paramName, value, logger},
		}
	}

	//TODO: return error
	panic(fmt.Sprintf("Unknown location \"%s\" for the API Key security parameter \"%s\".", location, paramName))
}

// GetName returns name.
func (sec Security) GetName() string {
	return sec.Name
}
