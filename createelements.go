package tgbot

import "go.uber.org/zap"

// holds the inputs and outputs of the previous render
// and the extracted local states tree
type RunResultWithStateTree[A any] struct {
	runResult      RunResultComponent[A]
	localStateTree *LocalStateTree[any]
}

type CreateElementsResult[A any] struct {
	Elements       []Element
	NewElements    []Element
	RemoveElements []Element
	TreeState      RunResultWithStateTree[A]
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

		return &CreateElementsResult[A]{
			Elements:       elements,
			NewElements:    elements,
			RemoveElements: make([]Element, 0),
			TreeState: RunResultWithStateTree[A]{
				runResult:      runResult,
				localStateTree: runResult.ExtractLocalStateTree(),
				// localStateTree: runResult.,
			},
		}
	}

	rerunResult := RerunComponentTree[A](
		&RerunContext[A]{
			logger:         GetLogger().With(zap.String("rerun", "rerun")),
			prevRunResult:  stateTree.runResult,
			localStateTree: *stateTree.localStateTree,
			componentIndex: []int{0},
			parents:        make([]ElementComponent[A], 0),
		},
		comp,
		// GetLogger(),
		// stateTree.runResult,
		// comp,
		// *stateTree.localStateTree,
		// []int{0},
		// make([]ElementComponent[A], 0),
	)

	aa := ExtractElementsFromRerun(rerunResult)

	if rr, ok := RunResultFromRerun[A](rerunResult).(*RunResultComponent[A]); ok {

		return &CreateElementsResult[A]{
			Elements:       aa.elements,
			NewElements:    aa.newElements,
			RemoveElements: aa.removedElements,
			TreeState: RunResultWithStateTree[A]{
				runResult:      *rr,
				localStateTree: stateTree.runResult.ExtractLocalStateTree(),
			},
		}
	}

	panic("not a run result")

}

// func RunResultFromRerunComp[A any](rerunResult RerunResult) RunResultComponent[A] {
// 	switch r := rerunResult.(type) {
// 	case *RerunResultUnchanged[A]:
// 		return RunResultFromRerun[A](r).(*RunResultComponent[A])
// 	case *RerunResultUpdated[A]:
// 		return RunResultFromRerun[A](r).(*RunResultComponent[A])
// 	case *ReRunResultElement:
// 		panic("should be elemenet")
// 	}
// }

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
		return r
	case *ReRunResultElement:
		return &r.element
	}

	return nil
}
