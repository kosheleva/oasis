package script

import (
	"os"
	"sync"

	"github.com/x1n13y84issmd42/gog/graph/comp"
	gcontract "github.com/x1n13y84issmd42/gog/graph/contract"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
)

// Executor executes an ExecutionGraph that comes from a script.
type Executor struct {
	contract.EntityTrait
}

// NewExecutor creates a new Executor instance.
func NewExecutor(log contract.Logger) *Executor {
	return &Executor{
		EntityTrait: contract.Entity(log),
	}
}

// Execute executes.
func (ex Executor) Execute(graph gcontract.Graph) {
	n0 := comp.MotherNode(graph)

	if n0 == nil {
		ex.Log.NOMESSAGE("Could not determine the starting node.")
		ex.Log.NOMESSAGE("It is a bug which is about to be fixed.")
		ex.Log.NOMESSAGE("At the moment please use simpler execution graphs.")
		ex.Log.NOMESSAGE("Aborting.")
		os.Exit(255)
	}

	// TODO: it returns nil sometimes (on script/noosa_test.yaml)
	ex.Log.ScriptExecutionStart(string(n0.ID()))

	results := make(contract.OperationResults)

	wg := sync.WaitGroup{}
	wg.Add(1)
	ex.Walk(graph, n0.(*ExecutionNode), &wg, &results)
	wg.Wait()

	success := true

	for nID, nRes := range results {
		if !nRes.Success {
			ex.Log.NOMESSAGE("Operation %s has failed.", nID)
			success = false
		}
	}

	if !success {
		os.Exit(255)
	}
}

// Walk walks the execution graph and executes operations.
func (ex Executor) Walk(
	graph gcontract.Graph,
	n *ExecutionNode,
	nwg *sync.WaitGroup,
	nresults *contract.OperationResults,
) {
	// Executing child nodes first (post-order).
	anwg := sync.WaitGroup{}
	anwg.Add(graph.AdjacentNodes(n.ID()).Count())
	anresults := contract.OperationResults{}

	for _an := range graph.AdjacentNodes(n.ID()).Range() {
		an := _an.(*ExecutionNode)
		go ex.Walk(graph, an, &anwg, &anresults)
	}

	anwg.Wait()

	//TODO: check for successful outcome of the previous ops.

	n.Lock()

	if n.Result == nil {
		// Executing the current node after it's children.
		n.Operation.Data().Reload()
		n.Operation.Data().Load(&n.Data)
		n.Operation.Data().URL.Load(n.Operation.Resolve().Host(""))

		enrichment := []contract.RequestEnrichment{
			n.Operation.Data().Query,
			n.Operation.Data().Headers,
			n.Operation.Data().Body,

			n.Operation.Resolve().Security(""),
		}

		ex.Log.TestingOperation(n.Operation)

		v := n.Operation.Resolve().Response(n.Expect.Status, "")

		v.Expect(expect.JSONBody(n.ExpectBody, graph, ex.Log))

		n.Result = test.Operation(n.Operation, &enrichment, v, ex.Log)
		(*nresults)[string(n.ID())] = n.Result
	}

	n.Unlock()

	nwg.Done()
}
