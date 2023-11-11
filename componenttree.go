package tgbot

import (
	"fmt"
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
	// props the component was rendered with
	inputProps reflect.Value
	// localState the component was rendered with
	inputLocalStateClosure LocalStateClosure[any]
	// elements the component rendered
	output []RunResult
}

func (c *RunResultComponent[A]) RunResultKind() string {
	return KindRunResultComponent
}

// extract local state tree from the component
func (c *RunResultComponent[A]) ExtractLocalStateTree() *LocalStateTree[any] {

	children := make([]*LocalStateTree[any], len(c.output))

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

	return &LocalStateTree[any]{
		localStateClosure: &c.inputLocalStateClosure,
		children:          &children,
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

type RerunResultUnchanged[A any] struct {
	RunResultComponent[A]
	rerunOutput []RerunResult
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
type ComponentsTreeState[A any] struct {
	runResult      RunResultComponent[A]
	localStateTree *LocalStateTree[any]
}

type RunContext[A any] struct {
	logger *zap.Logger

	localStateTree *LocalStateTree[any]

	// position of the component in the tree
	componentIndex []int
	parents        []ElementComponent[A]
}

type RerunContext[A any] struct {
	logger        *zap.Logger
	prevRunResult RunResultComponent[A]

	localStateTree LocalStateTree[any]

	// position of the component in the tree
	componentIndex []int
	parents        []ElementComponent[A]
}

func RunComponentTree[A any](ctx *RunContext[A], comp Comp[A]) RunResultComponent[A] {

	ctx.logger.Debug("RunComponentTree", zap.Any("comp", comp))

	if len(ctx.componentIndex) == 0 {
		// root component
		ctx.logger.Debug("Root component")
		ctx.componentIndex = []int{0}
	}

	if ctx.localStateTree == nil {
		ctx.logger.Debug("Creating local state tree")
		ctx.localStateTree = NewLocalStateTree[any]()
	}

	ctx.logger.Debug("Creating local state")
	localState := NewLocalState[any](ctx.componentIndex, ctx.localStateTree.localStateClosure)

	ctx.logger.Debug("Running the component")
	elements := RunComponent[A](comp, localState.Getset)

	ctx.logger.Debug("Elements was rendered",
		zap.Int("len", len(elements)),
		Elements(elements).ZapField("elements"),
	)

	childrenState := ctx.localStateTree.children

	if childrenState == nil || len(*ctx.localStateTree.children) != len(elements) {
		// initialize local states for the children
		cs := make([]*LocalStateTree[any], len(elements))
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
				componentIndex: append(ctx.componentIndex, idx),
				parents:        append(ctx.parents, *e),
			}, e.comp)
			output = append(output, &subcompres)
		default:
			output = append(output, &RunResultElement{e})
		}
	}

	runResult := RunResultComponent[A]{
		comp:                   comp,
		inputProps:             reflectCompProps[A](comp),
		inputLocalStateClosure: *localState.LocalState,
		output:                 output,
		compID:                 reflectCompId[A](comp),
	}

	ctx.logger.Debug("RunComponentTree done", zap.Any("runResult", runResult))

	return runResult

}

func RerunComponentTree[A any](
	ctx *RerunContext[A],
	comp Comp[A],
) RerunResult {

	localStateClosure := ctx.localStateTree.localStateClosure
	childrenState := ctx.localStateTree.children

	// localState := NewGetSet(ctx.componentIndex, localStateclosure)

	// signals if the component has to be rerun
	var rerun bool = false

	if ctx.prevRunResult.compID != reflectCompId[A](comp) {
		rerun = true
	} else if !reflect.DeepEqual(ctx.prevRunResult.inputProps, reflectCompProps[A](comp)) {
		rerun = true
	} else if !reflect.DeepEqual(ctx.prevRunResult.inputLocalStateClosure, localStateClosure) {
		rerun = true
	}

	if rerun {
		runResult := RunComponentTree(&RunContext[A]{
			logger: ctx.logger,
			localStateTree: &LocalStateTree[any]{
				localStateClosure: localStateClosure,
				children:          nil,
			},
			componentIndex: ctx.componentIndex,
			parents:        ctx.parents,
		}, comp)

		return &RerunResultUpdated[A]{
			oldChildren:        ctx.prevRunResult.output,
			RunResultComponent: runResult,
		}

	} else {
		returnOutput := make([]RerunResult, 0)

		if childrenState == nil || len(*childrenState) != len(ctx.prevRunResult.output) {
			cs := make([]*LocalStateTree[any], len(ctx.prevRunResult.output))
			childrenState = &cs
		}

		for idx, e := range ctx.prevRunResult.output {
			switch e := e.(type) {
			case *RunResultComponent[A]:
				rerunResult := RerunComponentTree[A](
					&RerunContext[A]{
						logger:         ctx.logger,
						prevRunResult:  *e,
						localStateTree: *(*childrenState)[idx],
						componentIndex: append(ctx.componentIndex, idx),
						parents:        append(ctx.parents, ElementComponent[A]{(*e).comp}),
					},
					(*e).comp,
				)

				returnOutput = append(returnOutput, rerunResult)
			case *RunResultElement:
				returnOutput = append(returnOutput, &ReRunResultElement{*e})
			}
		}

		return &RerunResultUnchanged[A]{
			RunResultComponent: RunResultComponent[A]{
				comp:                   comp,
				inputProps:             ctx.prevRunResult.inputProps,
				inputLocalStateClosure: ctx.prevRunResult.inputLocalStateClosure,
				output:                 ctx.prevRunResult.output,
				compID:                 ctx.prevRunResult.compID,
			},
			rerunOutput: returnOutput,
		}

	}
}

// sets the local state of the component if it has one defined
func ReflectCompLocalState[A any](comp Comp[A], ls GetSetLocalStateImpl[any]) Comp[A] {

	t := reflect.TypeOf(comp)

	fmt.Println("t: ", t)

	if ls.LocalState.Value == nil {
		// initialize local state
		return comp
	}

	stateField, ok := t.FieldByName("State")

	if !ok {
		return comp
	}

	v := reflect.ValueOf(comp)
	vp := reflect.ValueOf(&comp)

	v.Interface()

	fmt.Println("v: ", v)
	fmt.Println("v.Type(): ", v.Type())
	fmt.Println("vp: ", vp)
	fmt.Println("vp.Type(): ", vp.Type())

	fmt.Println("vp.Elem(): ", vp.Elem())
	fmt.Println("vp.Elem().Type(): ", vp.Elem().Type())
	fmt.Println("vp.Elem().Elem().Type(): ", vp.Elem().Elem().Type())

	sf := vp.Elem().Elem().FieldByName(stateField.Name)

	vls := reflect.ValueOf(ls)

	fmt.Printf("sf: %v\n", sf.Type())
	fmt.Printf("vls: %v\n", vls.Type())

	for i := 0; i < vls.NumField(); i++ {
		fmt.Printf("vls.Field(%v): %v\n", i, vls.Field(i).Type().Name())
	}

	vlsValue := vls.FieldByName("LocalState").FieldByName("Value")

	// fmt.Println("vlsValue: ", reflect.TypeOf(vlsValue))
	fmt.Println("vlsValue: ", reflect.TypeOf(vlsValue.Interface()))

	fmt.Println("sf.CanSet(): ", sf.CanSet())

	nt := reflect.New(t).Elem()

	fmt.Println("nt.Type(): ", nt.Type())

	// ntf := nt.Interface()
	// nts := reflect.ValueOf(&ntf).Elem()

	fmt.Println("nt.CanSet(): ", nt.CanSet())
	// fmt.Println("nt.CanSet(): ", nts.CanAddr())

	for i := 0; i < nt.NumField(); i++ {
		nt.Field(i).Set(v.Field(i))
		fmt.Printf("nt.Field(%v): %v\n", i, nt.Field(i))
	}

	fmt.Println("vls.Type()", vls.Type())
	fmt.Println("nt.state", nt.FieldByName("State").Type())

	nt.FieldByName("State").FieldByName("LocalState").FieldByName("Value").Set(
		reflect.ValueOf(vlsValue.Interface()),
	)

	nt.FieldByName("State").FieldByName("Index").Set(
		vls.FieldByName("Index"),
	)

	return nt.Interface().(Comp[A])
}

func RunComponent[A any](comp Comp[A], getset GetSetLocalStateImpl[any]) []Element {
	comp = ReflectCompLocalState[A](comp, getset)

	o := newOutput[A]()
	comp.Render(o)

	return o.result

}
