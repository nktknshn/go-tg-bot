package tgbot

const (
	KindRunResultComponent = "RunResultComponent"
	KindRunResultElement   = "RunResultElement"
)

type RunResult interface {
	RunResultKind() string
}

type RunResultComponent struct {
	comp            Comp[any]
	inputProps      any
	inputLocalState any
	output          []RunResult
}

func (c *RunResultComponent) RunResultKind() string {
	return KindRunResultComponent
}

type RunResultElement struct {
	element BasicElement
}

func (c *RunResultElement) RunResultKind() string {
	return KindRunResultElement
}

type LocalState[S any] struct {
	value   *S
	updated bool
}

type LocalStateTree[S any] struct {
	localState LocalState[S]
	children   []*LocalStateTree[any]
}

type TreeState struct {
}

func CreateElements[A any](comp Comp[A], treeState *TreeState) {
	if treeState == nil {

	}

}

func firstRun[A any](comp Comp[A]) {

}
