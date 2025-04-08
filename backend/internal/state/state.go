package state

type State struct {
	Env map[string]string
}

func NewState() *State {

	return &State{}

}
