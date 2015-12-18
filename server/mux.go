package server

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var patternMatch string

type Mux struct {
	sync.RWMutex
	ctx    context.Context
	cancel func()
	m      map[string]*Channel
	socket Socket
}

const PatternRegexp = `^(\w+:)+(\*|\w+){1}$`

var patternReg, _ = regexp.Compile(PatternRegexp)

var ErrWrongPattern = errors.New(
	`wrong pattern, pattern should be ` + PatternRegexp + " , like: rooms:alex",
)

func NewMux() *Mux {
	ctx, cancel := context.WithCancel(context.Background())
	return &Mux{
		ctx:    ctx,
		cancel: cancel,
		m:      make(map[string]*Channel),
	}
}

func (mux *Mux) Handle(ws *websocket.Conn) {
	req := ws.Request()
	if err := req.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	ctx, err := mux.socket.Connect(mux.ctx, req.Form)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := mux.socket.ID(ctx)

	_ = NewTransport(mux, ctx, ws, id)
	// start tr
	// read join msg from tr
	// join tr to channel
}

// Add pattern to Handler map
// pattern:
//  rooms:alex
//  rooms:*
func (mux *Mux) Add(pattern string, h *Channel) error {
	if !patternReg.MatchString(pattern) {
		return ErrWrongPattern
	}
	size := len(pattern)
	if pattern[size-2:] == ":*" {
		pattern = pattern[:size-2]
	}

	mux.Lock()
	mux.m[pattern] = h
	mux.Unlock()

	return nil
}

func (mux *Mux) get(pattern string) *Channel {
	mux.RLock()
	h := mux.m[pattern]
	mux.RUnlock()
	return h
}

func (mux *Mux) Find(pattern string) (*Channel, string) {
	if len(pattern) < 3 {
		return nil, ""
	}

	if h := mux.get(pattern); h != nil {
		return h, ""
	}

	for i := len(pattern) - 1; i >= 0; i-- {
		if pattern[i] == ':' {
			return mux.get(pattern[:i]), pattern[i+1:]
		}
	}

	return nil, ""
}
