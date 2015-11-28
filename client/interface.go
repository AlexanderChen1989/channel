package client

type GetPuter interface {
	Get() chan *Msg
	Put(chan *Msg)
}

type RefMaker interface {
	MakeRef() string
}

type PushPullRemover interface {
	Push(*Msg) error
	Pull(key interface{}, ch chan *Msg) error
	Remove(key interface{})
}

type Parent interface {
	GetPuter
	RefMaker
	PushPullRemover
}
