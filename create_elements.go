package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

// holds the inputs and outputs of the previous render
// and the extracted local states tree
type runResultWithStateTree struct {
	RunResult      runResultComponent
	LocalStateTree *localStateTree
}

type createElementsResult struct {
	Elements       []anyElement
	NewElements    []anyElement
	RemoveElements []anyElement
	TreeState      runResultWithStateTree
}

func (r createElementsResult) String() string {
	result := ""

	result += "CreateElementsResult"
	result += fmt.Sprintf("Elements: %v", elementsList(r.Elements))

	return result
}

// given
func createElements(
	comp Comp,
	gc globalContext[any],
	stateTree *runResultWithStateTree,
	logger *zap.Logger,
) *createElementsResult {
	// logger := GetLogger()
	// logger := zap.NewNop()

	logger.Debug("CreateElements",
		zap.String("compId", reflectCompId(comp)),
		zap.Any("props", reflectCompProps(comp)),
	)

	logger.Debug("StateTree", zap.Any("stateTree", stateTree))

	if stateTree == nil {
		logger.Debug("Running first time (stateTree == nil)")

		// this is the first render
		runResult := runComponentTree(&runContext{
			logger:         logger,
			globalContext:  gc,
			localStateTree: nil,
			componentIndex: []int{0},
			parents:        make([]elementComponent, 0),
		}, comp)

		elements := runResult.ExtractElements()

		logger.Debug("Extracting local state tree from the run")
		localStateTree := runResult.ExtractLocalStateTree()

		return &createElementsResult{
			Elements:       elements,
			NewElements:    elements,
			RemoveElements: make([]anyElement, 0),
			TreeState: runResultWithStateTree{
				RunResult:      runResult,
				LocalStateTree: localStateTree,
			},
		}
	}

	logger.Debug("This is not the first render (stateTree != nil)")
	// logger.Debug("used context", zap.Any("context", stateTree.RunResult.inputContext))

	rerunResult := rerunComponentTree(
		&rerunContext{
			logger:         logger,
			globalContext:  gc,
			prevRunResult:  stateTree.RunResult,
			localStateTree: *stateTree.LocalStateTree,
			componentIndex: []int{0},
			parents:        make([]elementComponent, 0),
		},
		comp,
	)

	logger.Debug("Extracting local state tree from rerun")

	var localStateTree *localStateTree

	switch r := rerunResult.(type) {
	case *rerunResultUnchanged:
		localStateTree = r.ExtractLocalStateTree()
	case *rerunResultUpdated:
		localStateTree = r.ExtractLocalStateTree()
	}

	logger.Debug("localStateTree", zap.String("localStateTree", localStateTree.String()))

	logger.Debug("Forming RunResult from rerun")

	rr := runResultFromRerun(rerunResult)

	aa := extractElementsFromRerun(rerunResult)

	if rrr, ok := rr.(*runResultComponent); ok {

		return &createElementsResult{
			Elements:       aa.elements,
			NewElements:    aa.newElements,
			RemoveElements: aa.removedElements,
			TreeState: runResultWithStateTree{
				RunResult:      *rrr,
				LocalStateTree: localStateTree,
			},
		}
	}

	// fmt.Print(rr)
	panic("not a run result")

}

func runResultFromRerun(rerunResult RerunResult) runResult {
	switch r := rerunResult.(type) {
	case *rerunResultUnchanged:
		output := make([]runResult, len(r.output))

		for i, o := range r.rerunOutput {
			output[i] = runResultFromRerun(o)
		}

		return &runResultComponent{
			comp:                   r.comp,
			compID:                 r.compID,
			inputProps:             r.inputProps,
			inputLocalStateClosure: r.inputLocalStateClosure,
			inputContext:           r.inputContext,
			output:                 output,
		}

	case *rerunResultUpdated:
		return &r.runResultComponent
	case *reRunResultElement:
		return &r.element
	}

	return nil
}
