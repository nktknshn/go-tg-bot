package tgbot

import (
	"reflect"

	"go.uber.org/zap"
)

const (
	// first run kinds
	KindRunResultComponent = "RunResultComponent"
	KindRunResultElement   = "RunResultElement"
	// reruns
	KindRerunResultElement   = "ReRunResultElement"
	KindRerunResultUpdated   = "RerunResultUpdated"
	KindRerunResultUnchanged = "RerunResultUnchanged"
)

// possible RunResults: RunResultComponent, RunResultElement
type RunResult interface {
	RunResultKind() string
}

// #region RunResultComponent
// holds the inputs and outputs of the very first render of the root component
// and all the subcomponents
type RunResultComponent[A any] struct {
	// the component that was rendered
	comp   Comp[A]
	compID string

	// context
	inputContext *UsedContextValue
	// props the component was rendered with
	inputProps any
	// localState the component was rendered with
	inputLocalStateClosure LocalStateClosure[any]
	// elements the component rendered
	output []RunResult
}

func (c *RunResultComponent[A]) RunResultKind() string {
	return KindRunResultComponent
}

// extract local state tree from the component
func (c *RunResultComponent[A]) ExtractLocalStateTree() *LocalStateTree {

	children := make([]*LocalStateTree, len(c.output))

	for idx, e := range c.output {
		switch e := e.(type) {
		case *RunResultComponent[A]:
			s := e.ExtractLocalStateTree()
			children[idx] = s

		case *RunResultElement:
			// nil state for elements
			children[idx] = nil
		}
	}

	// this is what getset is linked
	closureCopy := (any)(c.inputLocalStateClosure).(LocalStateClosure[any])

	// c.inputLocalStateClosure = LocalStateClosure[any]{
	// 	Initialized: closurePtr.Initialized,
	// 	Value:       closurePtr.Value,
	// }

	return &LocalStateTree{
		CompId:            c.compID,
		LocalStateClosure: &closureCopy,
		Children:          &children,
	}
}

// recursively extract elements from the component
func (c *RunResultComponent[A]) ExtractElements() []Element {
	result := make([]Element, 0)

	for _, e := range c.output {
		switch e := e.(type) {
		case *RunResultComponent[A]:
			result = append(result, e.ExtractElements()...)
		case *RunResultElement:
			result = append(result, e.element)
		}
	}

	return result
}

type RunResultElement struct {
	element Element
}

func (c *RunResultElement) RunResultKind() string {
	return KindRunResultElement
}

// #endregion

// possible ReRunResultElement, RerunResultUpdated, RerunResultUnchanged
type RerunResult interface {
	RerunResultKind() string
}

type RerunResultUpdated[A any] struct {
	RunResultComponent[A]
	oldChildren []RunResult
}

func (c *RerunResultUpdated[A]) ExtractLocalStateTree() *LocalStateTree {
	return c.RunResultComponent.ExtractLocalStateTree()
}

type RerunResultUnchanged[A any] struct {
	RunResultComponent[A]
	rerunOutput []RerunResult
}

func (c *RerunResultUnchanged[A]) ExtractLocalStateTree() *LocalStateTree {
	ac := RunResultFromRerun[A](c)
	bc := ac.(*RunResultComponent[A])
	return bc.ExtractLocalStateTree()

}

func (c *RerunResultUpdated[A]) RerunResultKind() string {
	return "RerunResultUpdated"
}

func (c *RerunResultUnchanged[A]) RerunResultKind() string {
	return "RerunResultUnchanged"
}

// just an element
type ReRunResultElement struct {
	element RunResultElement
}

func (c *ReRunResultElement) RerunResultKind() string {
	return "ReRunResultElement"
}

func (rr *RerunResultUnchanged[A]) ExtractElements() []Element {
	res := make([]Element, 0)

	for _, r := range rr.rerunOutput {
		res = append(res, ExtractElementsFromRerun(r).elements...)
	}

	return res
}

func (rr *RerunResultUnchanged[A]) ExtractNewElements() []Element {
	res := make([]Element, 0)

	for _, r := range rr.rerunOutput {
		res = append(res, ExtractElementsFromRerun(r).newElements...)
	}

	return res
}

type RerunExtractedElements struct {
	// elements that were rendered
	elements        []Element
	newElements     []Element
	removedElements []Element
}

func ExtractElementsFromRerun(rerun RerunResult) *RerunExtractedElements {

	switch r := rerun.(type) {
	case *RerunResultUpdated[any]:
		return &RerunExtractedElements{
			elements:        r.ExtractElements(),
			removedElements: ExtractFromRunResults(r.oldChildren),
			newElements:     r.ExtractElements(),
		}
	case *RerunResultUnchanged[any]:
		return &RerunExtractedElements{
			elements:    r.ExtractElements(),
			newElements: r.ExtractNewElements(),
		}
	case *ReRunResultElement:
		return &RerunExtractedElements{
			elements: []Element{r.element.element},
		}
	}

	return nil
}

func ExtractFromRunResults(results []RunResult) []Element {
	resultsElements := make([]Element, 0)

	for _, r := range results {
		switch r := r.(type) {
		case *RunResultComponent[any]:
			resultsElements = append(resultsElements, r.ExtractElements()...)
		case *RunResultElement:
			resultsElements = append(resultsElements, r.element)
		}
	}

	return resultsElements
}

// holds the inputs and outputs of the previous render
// and the extracted local states tree
// type ComponentsTreeState[A any] struct {
// 	runResult      RunResultComponent[A]
// 	localStateTree *LocalStateTree[any]
// }

type RunContext[A any] struct {
	logger *zap.Logger

	localStateTree *LocalStateTree

	globalContext GlobalContextTyped[any]

	// position of the component in the tree
	componentIndex []int
	parents        []ElementComponent[A]
}

type RerunContext[A any] struct {
	logger        *zap.Logger
	prevRunResult RunResultComponent[A]

	localStateTree LocalStateTree

	globalContext GlobalContextTyped[any]

	// position of the component in the tree
	componentIndex []int
	parents        []ElementComponent[A]
}

func RunComponentTree[A any](ctx *RunContext[A], comp Comp[A]) RunResultComponent[A] {

	ctx.logger.Debug("RunComponentTree",
		zap.String("compId", reflectCompId[A](comp)),
		zap.Any("props", reflectCompProps[A](comp)),
		zap.Any("index", ctx.componentIndex),
	)

	if len(ctx.componentIndex) == 0 {
		// root component
		ctx.logger.Debug("Root component")
		ctx.componentIndex = []int{0}
	}

	if ctx.localStateTree == nil {
		ctx.logger.Debug("This is the first run for the component. Creating local state tree.")
		ctx.localStateTree = NewLocalStateTree()
		ctx.localStateTree.CompId = reflectCompId[A](comp)
	} else {
		ctx.logger.Debug("Existing local state tree will be used",
			zap.Any("localStateTree", ctx.localStateTree.LocalStateClosure))
	}

	// creates new GetSet reusing ctx.localStateTree.LocalStateClosure
	localState := NewLocalState[any](
		ctx.componentIndex, ctx.localStateTree.LocalStateClosure,
	)

	ctx.logger.Debug("localState index", zap.Any("index", localState.Getset.Index))

	ctx.logger.Debug("Running the component")

	elements, closure, usedContextValue := RunComponent[A](
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
		Elements(elements).ZapField("elements"),
	)

	childrenState := ctx.localStateTree.Children

	if childrenState == nil || len(*ctx.localStateTree.Children) != len(elements) {
		// initialize local states for the children
		cs := make([]*LocalStateTree, len(elements))
		childrenState = &cs
	}

	output := make([]RunResult, 0)

	for idx, e := range elements {
		ctx.logger.Debug("Running element", zap.Int("idx", idx), zap.Any("element", e))

		switch e := e.(type) {
		case *ElementComponent[A]:

			subcompres := RunComponentTree(&RunContext[A]{
				logger:         ctx.logger,
				localStateTree: (*childrenState)[idx],
				globalContext:  ctx.globalContext,
				componentIndex: append(ctx.componentIndex, idx),
				parents:        append(ctx.parents, *e),
			}, e.comp)
			output = append(output, &subcompres)
		default:
			output = append(output, &RunResultElement{e})
		}
	}

	runResult := RunResultComponent[A]{
		comp:         comp,
		inputProps:   reflectCompProps[A](comp),
		inputContext: usedContextValue,
		// TODO FIX
		inputLocalStateClosure: closure,
		output:                 output,
		compID:                 reflectCompId[A](comp),
	}

	// ctx.logger.Debug("local state", zap.Any("localState", runResult.inputLocalStateClosure))

	ctx.logger.Debug("RunComponentTree done", zap.String("compId", reflectCompId[A](comp)))
	// , zap.Any("runResult", runResult)

	return runResult

}

// rerun the component tree
func RerunComponentTree[A any](
	ctx *RerunContext[A],
	comp Comp[A],
) RerunResult {

	ctx.logger.Debug("Detect if the component needs a rerun",
		zap.Any("comp", reflectCompId[A](comp)),
		zap.Any("props", reflectCompProps[A](comp)),
		zap.Any("used_context", ctx.prevRunResult.inputContext),
	)

	localStateClosure := ctx.localStateTree.LocalStateClosure

	ctx.logger.Debug("Local state",
		zap.Any("before", ctx.prevRunResult.inputLocalStateClosure.Value),
		zap.Any("now", localStateClosure.Value),
	)

	childrenState := ctx.localStateTree.Children

	currentUsedContext := ReflectTypedContextSelect[A](comp, ctx.globalContext.Get())

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

	if ctx.prevRunResult.compID != reflectCompId[A](comp) {
		ctx.logger.Debug("Component is different now",
			zap.String("compID", ctx.prevRunResult.compID),
			zap.String("newCompID", reflectCompId[A](comp)),
		)

		rerun = true
	} else if !reflect.DeepEqual(ctx.prevRunResult.inputProps, reflectCompProps[A](comp)) {
		ctx.logger.Debug("The Props has changed",
			zap.Any("before", ctx.prevRunResult.inputProps),
			zap.Any("now", reflectCompProps[A](comp)),
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
		ctx.logger.Debug("Rerunning component", zap.Any("comp", reflectCompId[A](comp)))

		runResult := RunComponentTree(&RunContext[A]{
			logger: ctx.logger,
			localStateTree: &LocalStateTree{
				CompId:            reflectCompId[A](comp),
				LocalStateClosure: localStateClosure,
				Children:          nil,
			},
			globalContext:  ctx.globalContext,
			componentIndex: ctx.componentIndex,
			parents:        ctx.parents,
		}, comp)

		return &RerunResultUpdated[A]{
			oldChildren:        ctx.prevRunResult.output,
			RunResultComponent: runResult,
		}

	} else {
		ctx.logger.Debug("Component hasn't changed. Rerunning children")

		returnOutput := make([]RerunResult, 0)

		if childrenState == nil || len(*childrenState) != len(ctx.prevRunResult.output) {
			cs := make([]*LocalStateTree, len(ctx.prevRunResult.output))
			childrenState = &cs
		}

		for idx, e := range ctx.prevRunResult.output {

			switch e := e.(type) {
			case *RunResultComponent[A]:
				ctx.logger.Debug("Reruning comp", zap.String("compId", e.compID))

				rerunResult := RerunComponentTree[A](
					&RerunContext[A]{
						logger:         ctx.logger,
						prevRunResult:  *e,
						localStateTree: *(*childrenState)[idx],
						globalContext:  ctx.globalContext,
						componentIndex: append(ctx.componentIndex, idx),
						parents:        append(ctx.parents, ElementComponent[A]{e.comp}),
					},
					e.comp,
				)

				ctx.logger.Debug("rerun done", zap.String("compId", e.compID))

				returnOutput = append(returnOutput, rerunResult)
			case *RunResultElement:
				returnOutput = append(returnOutput, &ReRunResultElement{*e})
			}
		}

		ctx.logger.Debug("RerunComponentTree done", zap.String("compId", reflectCompId[A](comp)))

		return &RerunResultUnchanged[A]{
			RunResultComponent: RunResultComponent[A]{
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

func RunComponent[A any](logger *zap.Logger, comp Comp[A], globalContext GlobalContextTyped[any], getset State[any]) ([]Element, LocalStateClosure[any], *UsedContextValue) {

	// contextQueryResult := ReflectContextQueryResultGet(comp, globalContext)
	logger.Debug("RunComponent",
		zap.String("compId", reflectCompId[A](comp)),
		zap.Any("globalContext", globalContext),
	)

	comp = ReflectCompLocalState[A](comp, getset)

	// refState := reflect.ValueOf(comp).Elem().FieldByName("State")
	// fmt.Println(refState)

	comp, usedContextValue := ReflectTypedContext[A](comp, globalContext.Get())

	// if contextQueryResult != nil {
	// 	comp = ReflectContextQueryResultSet[A](comp, contextQueryResult)
	// }

	o := NewOutput[A]()
	comp.Render(o)

	if !ReflectHasState(comp) {
		return o.Result, getset.LocalState, usedContextValue
	}

	// s, ok := reflect.ValueOf(comp).Elem().FieldByName("State")

	ls := ReflectDeref(reflect.ValueOf(comp)).FieldByName("State").FieldByName("LocalState")

	vi := ls.FieldByName("Initialized")
	vv := ls.FieldByName("Value")

	// fmt.Println("vi", vi)
	// fmt.Println("vv", vv)

	return o.Result, LocalStateClosure[any]{
		Initialized: vi.Bool(),
		Value:       vv.Interface(),
	}, usedContextValue

}
