package script

import (
	"github.com/x1n13y84issmd42/gog/graph/comp"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// OperationRef ...
type OperationRef struct {
	OperationID string           `yaml:"operationId"`
	Use         OperationDataUse `yaml:"use"`
	Expect      OperationDataUse `yaml:"expect"`
}

// OperationDataMap ...
type OperationDataMap map[string]string

// Iterate creates an iterable channel.
func (m OperationDataMap) Iterate() contract.ParameterIterator {
	ch := make(contract.ParameterIterator)

	go func() {
		for pN, pV := range m {
			ch <- contract.ParameterTuple{
				N: pN,
				Parameter: contract.Parameter{
					V:      params.Value(pV),
					Source: "script OperationRef",
				},
			}
		}

		close(ch)
	}()

	return ch
}

// OperationDataUse ...
type OperationDataUse struct {
	Path    OperationDataMap
	Query   OperationDataMap
	Headers OperationDataMap
}

// Script ...
type Script struct {
	contract.EntityTrait
	Spec       string                   `yaml:"spec"`
	Operations map[string]*OperationRef `yaml:"operations"`
}

// GetExecutionGraph builds and returns an operation execution graph.
func (script *Script) GetExecutionGraph() gcontract.Graph {
	spec := utility.Load(script.Spec, script.Log)
	graph := NewExecutionGraph(script.Log)

	for opRefID, opRef := range script.Operations {
		// script.Log.NOMESSAGE("Loop %s", opRefID)
		//TODO: opRef.OperationID may be absent, use opID then.
		op := spec.GetOperation(opRef.OperationID)
		opNode := script.GetNode(graph, opRefID, op)

		var err error

		err = script.SetupDependencies(spec, graph, &opRef.Use.Path, opNode.Data.URL, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Query, opNode.Data.Query, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		err = script.SetupDependencies(spec, graph, &opRef.Use.Headers, opNode.Data.Headers, opNode, opRefID)
		if err != nil {
			return NoGraph(err, script.Log)
		}

		// script.Log.NOMESSAGE("Execution node %s is ready.", opRefID)
		// opNode.Data.URL.Print()
		// script.Log.NOMESSAGE("")
		// script.Log.NOMESSAGE("")
	}

	// Checking for cycles.
	cycle := comp.Cycle(graph)
	if len(*cycle) > 0 {
		return NoGraph(errors.GraphHasCycles(cycle, nil), script.Log)
	}

	return graph
}

// SetupDependencies iterates over the provided map, looks for reference values,
// collects a list of references operations along with ParameterAccess functions.
func (script *Script) SetupDependencies(
	spec contract.Spec,
	graph *ExecutionGraph,
	srcParams *OperationDataMap,
	dstParams contract.Set,
	opNode *ExecutionNode,
	opRefID string,
) error {
	refParams := params.NewReferenceSource()
	memParams := params.NewMemorySource("script data")

	// script.Log.NOMESSAGE("Setting dependencies for the %s operation.", opRefID)

	for pn, pv := range *srcParams {
		isref, op2RefID, selector := Dereference(pv)
		if isref {
			// Retrieving the referenced operation.
			opRef2 := script.Operations[op2RefID]
			if opRef2 == nil {
				return errors.NotFound("Operation reference", op2RefID, nil)
			}

			op2 := spec.GetOperation(opRef2.OperationID)
			if opRef2 == nil {
				return errors.NotFound("Spec operation", opRef2.OperationID, nil)
			}

			// Adding the value so it's available for op later.
			refParams.AddReference(pn, op2.ID(), op2.Result(), selector)

			// Adding an edge to the execution graph.
			graph.AddEdge(opNode.ID(), script.GetNode(graph, op2RefID, op2).ID())
		} else {
			memParams.Add(pn, pv)
		}
	}

	// refParams.Print()

	dstParams.Load(refParams)
	dstParams.Load(memParams)

	return nil
}

// GetNode returns an ExecutionNode instance corresponding to the opRefID.
// If such a node exists in the graph, it is be returned, otherwise a new
// one is created.
func (script *Script) GetNode(graph gcontract.Graph, opRefID string, op contract.Operation) *ExecutionNode {
	var opNode *ExecutionNode
	_opNode := graph.Node(gcontract.NodeID(opRefID))
	if _opNode != nil {
		opNode = _opNode.(*ExecutionNode)
	} else {
		opNode = NewExecutionNode(op, opRefID, script.Log)
		graph.AddNode(opNode)
	}

	return opNode
}
