package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type match struct {
	match   string
	matched bool
	tail    string
}

var tcs = []struct {
	pattern      string
	wrongPattern bool
	matchs       []match
}{
	{"", true, nil},
	{"hello", true, nil},
	{"hello:", true, nil},
	{":xx", true, nil},
	{":*", true, nil},
	{"h:e", false, nil},
	{"hello:alex", false, []match{
		{"hello:alex", true, ""},
		{"hello:alex1", false, "alex1"},
	}},
	{"hello:*", false, []match{
		{"hell:no", false, "no"},
		{"hel", false, ""},
		{"hel:", false, ""},
		{":hel", false, "hel"},
		{":::", false, ""},
		{"hello:", true, ""},
		{"hello:alex", true, "alex"},
		{"hello:alex1", true, "alex1"},
	}},
}

func TestRouter(t *testing.T) {
	for _, tc := range tcs {
		r := NewMux()
		handler := &JoinHandler{}
		err := r.Add(tc.pattern, handler)

		if tc.wrongPattern {
			assert.Equal(t, err, ErrWrongPattern)
			continue
		}

		assert.Nil(t, err)
		for _, m := range tc.matchs {
			h, tail := r.Find(m.match)
			assert.Equal(t, m.tail, tail)
			if !m.matched {
				assert.Nil(t, h)
				continue
			}
			assert.NotNil(t, h)
		}
	}
}

func BenchmarkRouter(b *testing.B) {
	r := Mux{m: make(map[string]WSHandler)}
	r.Add("rooms:alex", &JoinHandler{})
	r.Add("rooms:*", &JoinHandler{})
	for i := 0; i < b.N; i++ {
		r.Find("rooms:alex")
	}
}
