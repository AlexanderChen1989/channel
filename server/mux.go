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
	m      map[string]WSHandler
	socket Socket
	conns  map[string]*websocket.Conn
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
		m:      make(map[string]WSHandler),
		conns:  make(map[string]*websocket.Conn),
	}
}

func (r *Mux) Handle(ws *websocket.Conn) {
	req := ws.Request()
	if err := req.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	ctx, err := r.socket.Connect(r.ctx, req.Form)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := r.socket.ID(ctx)
	fmt.Println(id)
}

// Add pattern to Handler map
// pattern:
//  rooms:alex
//  rooms:*
func (r *Mux) Add(pattern string, h WSHandler) error {
	if !patternReg.MatchString(pattern) {
		return ErrWrongPattern
	}
	size := len(pattern)
	if pattern[size-2:] == ":*" {
		pattern = pattern[:size-2]
	}

	r.Lock()
	r.m[pattern] = h
	r.Unlock()

	return nil
}

func (r *Mux) Del(patten string) {
	r.Lock()
	delete(r.m, patten)
	r.Unlock()
}

func (r *Mux) get(pattern string) WSHandler {
	r.RLock()
	h := r.m[pattern]
	r.RUnlock()
	return h
}

func (r *Mux) Find(pattern string) (WSHandler, string) {
	if len(pattern) < 3 {
		return nil, ""
	}

	if h := r.get(pattern); h != nil {
		return h, ""
	}

	for i := len(pattern) - 1; i >= 0; i-- {
		if pattern[i] == ':' {
			return r.get(pattern[:i]), pattern[i+1:]
		}
	}

	return nil, ""
}
