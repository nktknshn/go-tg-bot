package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

// holds the inputs and outputs of the previous render
// and the extracted local states tree
type RunResultWithStateTree[A any] struct {
	RunResult      RunResultComponent[A]
	LocalStateTree *LocalStateTree
}

type CreateElementsResult[A any] struct {
	Elements       []Element
	NewElements    []Element
	RemoveElements []Element
	TreeState      RunResultWithStateTree[A]
}

func (r CreateElementsResult[any]) String() string {
	result := ""

	result += "CreateElementsResult"
	result += fmt.Sprintf("Elements: %v", Elements(r.Elements))

	return result
}

// given
func CreateElements[A any](comp Comp[A], stateTree *RunResultWithStateTree[A]) *CreateElementsResult[A] {
	logger := GetLogger()

	logger.Debug("CreateElements", zap.Any("comp", comp), zap.Any("stateTree", stateTree))

	if stateTree == nil {
		logger.Debug("Running first time (stateTree == nil)")

		// this is the first render
		runResult := RunComponentTree[A](&RunContext[A]{
			logger:         GetLogger(),
			localStateTree: nil,
			componentIndex: []int{0},
			parents:        make([]ElementComponent[A], 0),
		}, comp)

		elements := runResult.ExtractElements()

		logger.Debug("Extracting local state tree from the run")
		localStateTree := runResult.ExtractLocalStateTree()

		return &CreateElementsResult[A]{
			Elements:       elements,
			NewElements:    elements,
			RemoveElements: make([]Element, 0),
			TreeState: RunResultWithStateTree[A]{
				RunResult:      runResult,
				LocalStateTree: localStateTree,
			},
		}
	}

	rerunResult := RerunComponentTree[A](
		&RerunContext[A]{
			logger:         GetLogger().With(zap.String("rerun", "rerun")),
			prevRunResult:  stateTree.RunResult,
			localStateTree: *stateTree.LocalStateTree,
			componentIndex: []int{0},
			parents:        make([]ElementComponent[A], 0),
		},
		comp,
	)

	aa := ExtractElementsFromRerun(rerunResult)

	logger.Debug("Extracting local state tree from rerun")

	localStateTree := stateTree.RunResult.ExtractLocalStateTree()

	logger.Debug("Forming RunResult from rerun")

	rr := RunResultFromRerun[A](rerunResult)

	if rrr, ok := rr.(*RunResultComponent[A]); ok {

		return &CreateElementsResult[A]{
			Elements:       aa.elements,
			NewElements:    aa.newElements,
			RemoveElements: aa.removedElements,
			TreeState: RunResultWithStateTree[A]{
				RunResult:      *rrr,
				LocalStateTree: localStateTree,
			},
		}
	}

	// fmt.Print(rr)
	panic("not a run result")

}

func RunResultFromRerun[A any](rerunResult RerunResult) RunResult {
	switch r := rerunResult.(type) {
	case *RerunResultUnchanged[A]:
		output := make([]RunResult, len(r.output))

		for i, o := range r.rerunOutput {
			output[i] = RunResultFromRerun[A](o)
		}

		return &RunResultComponent[A]{
			comp:                   r.comp,
			compID:                 r.compID,
			inputProps:             r.inputProps,
			inputLocalStateClosure: r.inputLocalStateClosure,
			output:                 output,
		}

	case *RerunResultUpdated[A]:
		return &r.RunResultComponent
	case *ReRunResultElement:
		return &r.element
	}

	return nil
}
