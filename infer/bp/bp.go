package bp

import (
	"github.com/jesand/stats"
	"github.com/jesand/stats/factor"
	"github.com/jesand/stats/variable"
	"sync"
)

// Perform exact belief propogation on a factor graph with a tree structure.
// This works by iteratively passing messages from nodes with at most one
// pending in-message to all adjacent nodes. If all possible messages are
// passed, the exact marginals over the factor graph will have been calculated.
// If any messages cannot be passed due to the graph structure, the method
// will fail with an error.
func InferForTree(graph factor.FactorGraph) error {

	// Create the BP tree structure
	tree := buildBPTree(graph, 1)

	// Prepare the message processing pipeline
	const mux = 10
	var (
		msgC  = make(chan *bpNode, mux)
		doneC = make(chan *bpNode, mux)
		wait  sync.WaitGroup
	)
	for i := 0; i < mux; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for node := range msgC {

				doneC <- node
			}
		}()
	}
	wait.Add(1)
	go func() {
		defer wait.Done()
		// for node := range doneC {

		// }
	}()

	// Build the queue of leaf nodes
	var queue []*bpNode
	for _, node := range tree {
		if len(node.In) <= 1 {
			queue = append(queue, node)
		}
	}
	if len(queue) == 0 {
		return stats.ErrGraphNotTree
	}

	// Move down the queue until all messages have been sent
	// for len(queue) > 0 {

	// }
	return nil
}

func buildBPTree(graph factor.FactorGraph, initialMessage float64) (tree bpGraph) {
	var (
		vNodes = make(map[*variable.DiscreteRV]*bpNode)
	)
	for _, f := range graph.Factors {
		fNode := &bpNode{Factor: f}
		tree = append(tree, fNode)
		for _, v := range graph.AdjToFactor(f) {
			dv, ok := v.(*variable.DiscreteRV)
			if !ok {
				panic(stats.ErrDiscreteOnly)
			}
			vNode, ok := vNodes[dv]
			if !ok {
				vNode = &bpNode{Variable: dv}
				tree = append(tree, vNode)
				vNodes[dv] = vNode
			}
			var (
				vals   = dv.Space().Size()
				msgOut = &bpMessage{From: fNode, To: vNode, Value: make([]float64, vals)}
				msgIn  = &bpMessage{From: vNode, To: fNode, Value: make([]float64, vals)}
			)
			for i := 0; i < vals; i++ {
				msgOut.Value[i] = initialMessage
				msgIn.Value[i] = initialMessage
			}
			fNode.Out = append(fNode.Out, msgOut)
			fNode.In = append(fNode.In, msgIn)
			vNode.Out = append(vNode.Out, msgIn)
			vNode.In = append(vNode.In, msgOut)
		}
	}
	return
}

type bpGraph []*bpNode

type bpNode struct {
	Variable *variable.DiscreteRV
	Factor   factor.Factor
	Out      []*bpMessage
	In       []*bpMessage
}

type bpMessage struct {

	// The current message value for each variable assignment in the To node.
	// We assume only discrete random variables are used.
	Value []float64

	// The number of times this message has been updated
	Iter int

	// The nodes involved in the message
	From, To *bpNode
}
