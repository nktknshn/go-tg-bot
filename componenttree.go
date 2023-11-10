package tgbot

import (
	"reflect"

	"go.uber.org/zap"
)

const (
	KindRunResultComponent = "RunResultComponent"
	KindRunResultElement   = "RunResultElement"
)

type RunResult interface {
	RunResultKind() string
}

type RunResultComponent[A any] struct {
	element         ElementComponent[A]
	comp            Comp[A]
	compID          string
	inputProps      reflect.Value
	inputLocalState LocalStateClosure[any]
	output          []RunResult
}

func (c *RunResultComponent[A]) RunResultKind() string {
	return KindRunResultComponent
}

type RunResultElement struct {
	element Element
}

func (c *RunResultElement) RunResultKind() string {
	return KindRunResultElement
}

type TreeState struct {
}

func CreateElements[A any](comp Comp[A], treeState *TreeState) {
	if treeState == nil {

	}

}

type RunContext[A any] struct {
	logger *zap.Logger

	localStateTree *LocalStateTree[any]

	// position of the component in the tree
	componentIndex []int
	parents        []ElementComponent[A]
}

func reflectCompProps[A any](comp Comp[A]) reflect.Value {
	t := reflect.TypeOf(comp)

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "state" {
			continue
		}

		fs = append(fs, f)
	}

	props := reflect.New(reflect.StructOf(fs)).Elem()
	props.Set(reflect.ValueOf(comp))

	return props
}

func reflectCompId[A any](comp Comp[A]) string {
	t := reflect.TypeOf(comp)
	return t.Name()
}

func reflectCompLocalState[A any](comp Comp[A], ls GetSetLocalStateImpl[any]) {
	t := reflect.TypeOf(comp)

	f, ok := t.FieldByName("state")

	if ok == false {
		return
	}

	reflect.ValueOf(comp).FieldByName(f.Name).Set(reflect.ValueOf(ls))

}

func RunComponentTree[A any](ctx *RunContext[A], comp Comp[A]) RunResultComponent[A] {

	// assert if ctx.localStateTree == nil then len(ctx.componentIndex) == 0
	// and vice versa

	// if ctx.localStateTree == nil && len(ctx.componentIndex) != 0 {
	// 	panic("invalid state")
	// }

	// if len(ctx.componentIndex) == 0 && ctx.localStateTree != nil {
	// 	panic("invalid state")
	// }

	if len(ctx.componentIndex) == 0 {
		// root component
		ctx.componentIndex = []int{0}
		ctx.logger.Debug("Root component")
	}

	if ctx.localStateTree == nil {
		ctx.localStateTree = NewLocalStateTree[any]()
	}

	ctx.logger.Debug("firstRun", zap.Any("comp", comp))

	localState := NewLocalState[any](ctx.componentIndex, ctx.localStateTree.localState)

	elements := RunComponent[A](comp, localState.getset)

	if ctx.localStateTree.children == nil || len(*ctx.localStateTree.children) != len(elements) {
		childrenState := make([]LocalStateTree[any], len(elements))
		ctx.localStateTree.children = &childrenState
	}

	output := make([]RunResult, 0)

	for idx, e := range elements {
		switch e := e.(type) {
		case *ElementComponent[A]:
			subcompres := RunComponentTree(&RunContext[A]{
				logger:         ctx.logger,
				localStateTree: &(*ctx.localStateTree.children)[idx],
				componentIndex: append(ctx.componentIndex, idx),
				parents:        append(ctx.parents, *e),
			}, e.comp)
			output = append(output, &subcompres)
		default:
			output = append(output, &RunResultElement{e})
		}
	}

	runResult := RunResultComponent[A]{
		comp:            comp,
		inputProps:      reflectCompProps[A](comp),
		inputLocalState: *ctx.localStateTree.localState,
		output:          output,
		compID:          reflectCompId[A](comp),
	}

	return runResult

}

type RerunResult interface {
	RerunResultKind() string
}

type ReRunResultElement struct {
	element RunResultElement
}

func (c *ReRunResultElement) RerunResultKind() string {
	return "ReRunResultElement"
}

type RerunResultUpdated[A any] struct {
	oldChildren []RunResult
	RunResultComponent[A]
}

func (c *RerunResultUpdated[A]) RerunResultKind() string {
	return "RerunResultUpdated"
}

type RerunResultUnchanged[A any] struct {
	RunResultComponent[A]
	RerunOutput []RerunResult
}

func (c *RerunResultUnchanged[A]) RerunResultKind() string {
	return "RerunResultUnchanged"
}

func RerunComponentTree[A any](
	logger *zap.Logger,
	prevRunResult RunResultComponent[A],
	comp Comp[A],
	localStateTree LocalStateTree[any],
	index []int,
	parents []ElementComponent[A],
) RerunResult {

	ls := localStateTree.localState
	childrenState := localStateTree.children

	localState := createGetSet(index, ls)

	var rerender bool = false

	if prevRunResult.compID != reflectCompId[A](comp) {
		rerender = true
	}

	if !reflect.DeepEqual(prevRunResult.inputProps, reflectCompProps[A](comp)) {
		rerender = true
	}

	if !reflect.DeepEqual(prevRunResult.inputLocalState, ls) {
		rerender = true
	}

	if rerender {
		runResult := RunComponentTree(&RunContext[A]{
			logger: logger,
			localStateTree: &LocalStateTree[any]{
				localState: localState.localState,
				children:   nil,
			},
			componentIndex: index,
			parents:        parents,
		}, comp)

		return &RerunResultUpdated[A]{
			oldChildren:        prevRunResult.output,
			RunResultComponent: runResult,
		}

	} else {
		returnOutput := make([]RerunResult, 0)

		if childrenState == nil || len(*childrenState) != len(prevRunResult.output) {
			cs := make([]LocalStateTree[any], len(prevRunResult.output))
			childrenState = &cs
		}

		for idx, e := range prevRunResult.output {
			switch e := e.(type) {
			case *RunResultComponent[A]:
				rerunResult := RerunComponentTree[A](
					logger,
					*e,
					(*e).comp,
					(*childrenState)[idx],
					append(index, idx),
					append(parents, ElementComponent[A]{(*e).comp}),
				)

				returnOutput = append(returnOutput, rerunResult)
			case *RunResultElement:
				returnOutput = append(returnOutput, &ReRunResultElement{*e})
			}
		}

		return &RerunResultUnchanged[A]{
			RunResultComponent: RunResultComponent[A]{
				comp:            comp,
				inputProps:      prevRunResult.inputProps,
				inputLocalState: prevRunResult.inputLocalState,
				output:          prevRunResult.output,
				compID:          prevRunResult.compID,
			},
			RerunOutput: returnOutput,
		}

	}
}

func RunComponent[A any](comp Comp[A], getset GetSetLocalStateImpl[any]) []Element {
	reflectCompLocalState[A](comp, getset)

	o := newOutput[A]()
	comp.Render(o)

	return o.result

}
