package tgbot

import (
	"reflect"

	"go.uber.org/zap"
)

const (
	// first run kinds
	kindRunResultComponent = "RunResultComponent"
	kindRunResultElement   = "RunResultElement"
	// reruns
	kindRerunResultElement   = "ReRunResultElement"
	kindRerunResultUpdated   = "RerunResultUpdated"
	kindRerunResultUnchanged = "RerunResultUnchanged"
)

// possible RunResults: RunResultComponent, RunResultElement
type runResult interface {
	RunResultKind() string
}

type usedContextValue []reflect.Value

// #region runResultComponent
// holds the inputs and outputs of the very first render of the root component
// and all the subcomponents
type runResultComponent struct {
	// the component that was rendered
	comp   Comp
	compID string

	// context
	inputContext *usedContextValue
	// props the component was rendered with
	inputProps any
	// localState the component was rendered with
	inputLocalStateClosure localStateClosure[any]
	// elements the component rendered
	output []runResult
}

func (c *runResultComponent) RunResultKind() string {

	return kindRunResultComponent
}

// extract local state tree from the component
func (c *runResultComponent) ExtractLocalStateTree() *localStateTree {

	children := make([]*localStateTree, len(c.output))

	for idx, e := range c.output {
		switch e := e.(type) {
		case *runResultComponent:
			s := e.ExtractLocalStateTree()
			children[idx] = s

		case *runResultElement:
			// nil state for elements
			children[idx] = nil
		}
	}

	// this is what getset is linked
	closureCopy := (any)(c.inputLocalStateClosure).(localStateClosure[any])

	// c.inputLocalStateClosure = LocalStateClosure[any]{
	// 	Initialized: closurePtr.Initialized,
	// 	Value:       closurePtr.Value,
	// }

	return &localStateTree{
		CompId:            c.compID,
		LocalStateClosure: &closureCopy,
		Children:          &children,
	}
}

// recursively extract elements from the component
func (c *runResultComponent) ExtractElements() []anyElement {
	result := make([]anyElement, 0)

	for _, e := range c.output {
		switch e := e.(type) {
		case *runResultComponent:
			result = append(result, e.ExtractElements()...)
		case *runResultElement:
			result = append(result, e.element)
		}
	}

	return result
}

type runResultElement struct {
	element anyElement
}

func (c *runResultElement) RunResultKind() string {
	return kindRunResultElement
}

// #endregion

// possible ReRunResultElement, RerunResultUpdated, RerunResultUnchanged
type RerunResult interface {
	RerunResultKind() string
}

type rerunResultUpdated struct {
	runResultComponent
	oldChildren []runResult
}

func (c *rerunResultUpdated) ExtractLocalStateTree() *localStateTree {
	return c.runResultComponent.ExtractLocalStateTree()
}

type rerunResultUnchanged struct {
	runResultComponent
	rerunOutput []RerunResult
}

func (c *rerunResultUnchanged) ExtractLocalStateTree() *localStateTree {
	ac := runResultFromRerun(c)
	bc := ac.(*runResultComponent)
	return bc.ExtractLocalStateTree()

}

func (c *rerunResultUpdated) RerunResultKind() string {
	return "RerunResultUpdated"
}

func (c *rerunResultUnchanged) RerunResultKind() string {
	return "RerunResultUnchanged"
}

// just an element
type reRunResultElement struct {
	element runResultElement
}

func (c *reRunResultElement) RerunResultKind() string {
	return "ReRunResultElement"
}

func (rr *rerunResultUnchanged) ExtractElements() []anyElement {
	res := make([]anyElement, 0)

	for _, r := range rr.rerunOutput {
		res = append(res, extractElementsFromRerun(r).elements...)
	}

	return res
}

func (rr *rerunResultUnchanged) ExtractNewElements() []anyElement {
	res := make([]anyElement, 0)

	for _, r := range rr.rerunOutput {
		res = append(res, extractElementsFromRerun(r).newElements...)
	}

	return res
}

type rerunExtractedElements struct {
	// elements that were rendered
	elements        []anyElement
	newElements     []anyElement
	removedElements []anyElement
}

func extractElementsFromRerun(rerun RerunResult) rerunExtractedElements {

	switch r := rerun.(type) {
	case *rerunResultUpdated:
		return rerunExtractedElements{
			elements:        r.ExtractElements(),
			removedElements: extractFromRunResults(r.oldChildren),
			newElements:     r.ExtractElements(),
		}
	case *rerunResultUnchanged:
		return rerunExtractedElements{
			elements:    r.ExtractElements(),
			newElements: r.ExtractNewElements(),
		}
	case *reRunResultElement:
		return rerunExtractedElements{
			elements: []anyElement{r.element.element},
		}
	}

	return rerunExtractedElements{}
}

func extractFromRunResults(results []runResult) []anyElement {
	resultsElements := make([]anyElement, 0)

	for _, r := range results {
		switch r := r.(type) {
		case *runResultComponent:
			resultsElements = append(resultsElements, r.ExtractElements()...)
		case *runResultElement:
			resultsElements = append(resultsElements, r.element)
		}
	}

	return resultsElements
}

// holds the inputs and outputs of the previous render
// and the extracted local states tree
type runContext struct {
	logger *zap.Logger

	localStateTree *localStateTree

	globalContext globalContext[any]

	// position of the component in the tree
	componentIndex []int
	parents        []elementComponent
}

type rerunContext struct {
	logger        *zap.Logger
	prevRunResult runResultComponent

	localStateTree localStateTree

	globalContext globalContext[any]

	// position of the component in the tree
	componentIndex []int
	parents        []elementComponent
}

func runComponentTree(ctx *runContext, comp Comp) runResultComponent {

	ctx.logger.Debug("RunComponentTree",
		zap.String("compId", reflectCompId(comp)),
		zap.Any("props", reflectCompId(comp)),
		zap.Any("index", ctx.componentIndex),
	)

	if len(ctx.componentIndex) == 0 {
		// root component
		ctx.logger.Debug("Root component")
		ctx.componentIndex = []int{0}
	}

	if ctx.localStateTree == nil {
		ctx.logger.Debug("This is the first run for the component. Creating local state tree.")
		ctx.localStateTree = newLocalStateTree()
		ctx.localStateTree.CompId = reflectCompId(comp)
	} else {
		ctx.logger.Debug("Existing local state tree will be used",
			zap.Any("localStateTree", ctx.localStateTree.LocalStateClosure))
	}

	// creates new GetSet reusing ctx.localStateTree.LocalStateClosure
	localState := newLocalState[any](
		ctx.componentIndex, ctx.localStateTree.LocalStateClosure,
	)

	ctx.logger.Debug("localState index", zap.Any("index", localState.Getset.Index))

	ctx.logger.Debug("Running the component")

	elements, closure, usedContextValue := runComponent(
		ctx.logger, comp, ctx.globalContext, localState.Getset,
	)

	ctx.logger.Debug("Local state status after the run",
		zap.Bool("initialzied", closure.Initialized),
		zap.Any("value", closure.Value),
	)

	if usedContextValue == nil {
		ctx.logger.Debug("The component doesn't use global context")
	} else {
		ctx.logger.Debug("Used context value", zap.Any("usedContextValue", usedContextValue))
	}

	localState.LocalState = &closure

	ctx.logger.Debug("Elements was rendered",
		zap.Int("len", len(elements)),
		elementsList(elements).ZapField("elements"),
	)

	childrenState := ctx.localStateTree.Children

	if childrenState == nil || len(*ctx.localStateTree.Children) != len(elements) {
		// initialize local states for the children
		cs := make([]*localStateTree, len(elements))
		childrenState = &cs
	}

	output := make([]runResult, 0)

	for idx, e := range elements {
		ctx.logger.Debug("Running element", zap.Int("idx", idx), zap.Any("element", e))

		switch e := e.(type) {
		case *elementComponent:

			subcompres := runComponentTree(&runContext{
				logger:         ctx.logger,
				localStateTree: (*childrenState)[idx],
				globalContext:  ctx.globalContext,
				componentIndex: append(ctx.componentIndex, idx),
				parents:        append(ctx.parents, *e),
			}, e.comp)
			output = append(output, &subcompres)
		default:
			output = append(output, &runResultElement{e})
		}
	}

	runResult := runResultComponent{
		comp:         comp,
		inputProps:   reflectCompProps(comp),
		inputContext: usedContextValue,
		// TODO FIX
		inputLocalStateClosure: closure,
		output:                 output,
		compID:                 reflectCompId(comp),
	}

	// ctx.logger.Debug("local state", zap.Any("localState", runResult.inputLocalStateClosure))

	ctx.logger.Debug("RunComponentTree done", zap.String("compId", reflectCompId(comp)))
	// , zap.Any("runResult", runResult)

	return runResult

}

// rerun the component tree
func rerunComponentTree(
	ctx *rerunContext,
	comp Comp,
) RerunResult {

	ctx.logger.Debug("Detect if the component needs a rerun",
		zap.Any("comp", reflectCompId(comp)),
		zap.Any("props", reflectCompProps(comp)),
		zap.Any("used_context", ctx.prevRunResult.inputContext),
	)

	localStateClosure := ctx.localStateTree.LocalStateClosure

	ctx.logger.Debug("Local state",
		zap.Any("before", ctx.prevRunResult.inputLocalStateClosure.Value),
		zap.Any("now", localStateClosure.Value),
	)

	childrenState := ctx.localStateTree.Children

	currentUsedContext := reflectTypedContextSelect[any](comp, ctx.globalContext.Get())

	// localState := NewGetSet(ctx.componentIndex, localStateclosure)

	// signals if the component has to be rerun
	var rerun bool = false

	bothContextsNil := ctx.prevRunResult.inputContext == nil && currentUsedContext == nil
	onlyOneContextNil1 := ctx.prevRunResult.inputContext == nil && currentUsedContext != nil
	onlyOneContextNil2 := ctx.prevRunResult.inputContext != nil && currentUsedContext == nil

	if bothContextsNil {
		ctx.logger.Debug("Both contexts are nil")
	} else if onlyOneContextNil1 || onlyOneContextNil2 {
		ctx.logger.Debug("One of the contexts is nil",
			zap.Any("before", ctx.prevRunResult.inputContext),
			zap.Any("now", currentUsedContext),
		)

		rerun = true
	} else if !ctx.prevRunResult.inputContext.Equal(*currentUsedContext) {
		ctx.logger.Debug("The global context has changed",
			zap.Any("before", ctx.prevRunResult.inputContext),
			zap.Any("now", currentUsedContext),
		)

		rerun = true
	} else {
		ctx.logger.Debug("The global context is the same")
	}

	if ctx.prevRunResult.compID != reflectCompId(comp) {
		ctx.logger.Debug("Component is different now",
			zap.String("compID", ctx.prevRunResult.compID),
			zap.String("newCompID", reflectCompId(comp)),
		)

		rerun = true
	} else if !reflect.DeepEqual(ctx.prevRunResult.inputProps, reflectCompProps(comp)) {
		ctx.logger.Debug("The Props has changed",
			zap.Any("before", ctx.prevRunResult.inputProps),
			zap.Any("now", reflectCompProps(comp)),
		)

		rerun = true

	} else if !reflect.DeepEqual(ctx.prevRunResult.inputLocalStateClosure.Value, localStateClosure.Value) {
		ctx.logger.Debug("The LocalState has changed",
			zap.Any("before", ctx.prevRunResult.inputLocalStateClosure),
			zap.Any("now", localStateClosure),
		)
		rerun = true
	}

	if rerun {
		ctx.logger.Debug("Rerunning component", zap.Any("comp", reflectCompId(comp)))

		runResult := runComponentTree(&runContext{
			logger: ctx.logger,
			localStateTree: &localStateTree{
				CompId:            reflectCompId(comp),
				LocalStateClosure: localStateClosure,
				Children:          nil,
			},
			globalContext:  ctx.globalContext,
			componentIndex: ctx.componentIndex,
			parents:        ctx.parents,
		}, comp)

		return &rerunResultUpdated{
			oldChildren:        ctx.prevRunResult.output,
			runResultComponent: runResult,
		}

	} else {
		ctx.logger.Debug("Component hasn't changed. Rerunning children")

		returnOutput := make([]RerunResult, 0)

		if childrenState == nil || len(*childrenState) != len(ctx.prevRunResult.output) {
			cs := make([]*localStateTree, len(ctx.prevRunResult.output))
			childrenState = &cs
		}

		for idx, e := range ctx.prevRunResult.output {

			switch e := e.(type) {
			case *runResultComponent:
				ctx.logger.Debug("Reruning comp", zap.String("compId", e.compID))

				rerunResult := rerunComponentTree(
					&rerunContext{
						logger:         ctx.logger,
						prevRunResult:  *e,
						localStateTree: *(*childrenState)[idx],
						globalContext:  ctx.globalContext,
						componentIndex: append(ctx.componentIndex, idx),
						parents:        append(ctx.parents, elementComponent{e.comp}),
					},
					e.comp,
				)

				ctx.logger.Debug("rerun done", zap.String("compId", e.compID))

				returnOutput = append(returnOutput, rerunResult)
			case *runResultElement:
				returnOutput = append(returnOutput, &reRunResultElement{*e})
			}
		}

		ctx.logger.Debug("RerunComponentTree done", zap.String("compId", reflectCompId(comp)))

		return &rerunResultUnchanged{
			runResultComponent: runResultComponent{
				comp:                   comp,
				inputContext:           ctx.prevRunResult.inputContext,
				inputProps:             ctx.prevRunResult.inputProps,
				inputLocalStateClosure: ctx.prevRunResult.inputLocalStateClosure,
				output:                 ctx.prevRunResult.output,
				compID:                 ctx.prevRunResult.compID,
			},
			rerunOutput: returnOutput,
		}

	}
}
