package maps

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mapT = StdMap[string, int]
type mapTI = MapI[string, int]

func TestNewStdMap(t *testing.T) {
	m := NewStdMap(map[string]int{"a": 1})
	assert.Equal(t, 1, m.Get("a"))
}

func TestMap_Clear(t *testing.T) {
	var mNil mapT
	tests := []struct {
		name string
		m    mapTI
	}{
		{"zero", mNil},
		{"empty", NewStdMap[string, int]()},
		{"1 item", mapT{"a": 1}},
		{"2 items", mapT{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
			if tt.m.Len() != 0 {
				t.Errorf("StdMap not cleared: %q", tt.m)
			}
		})
	}
}

func TestMap_Merge(t *testing.T) {
	tests := []struct {
		name     string
		m1       mapT
		m2       mapT
		expected mapT
	}{
		{"1 to 1", mapT{"a": 1}, mapT{"b": 2}, mapT{"a": 1, "b": 2}},
		{"overwrite", mapT{"a": 1}, mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"to empty", mapT{}, mapT{"a": 1, "b": 2}, mapT{"a": 1, "b": 2}},
		{"from empty", mapT{"a": 1, "b": 2}, mapT{}, mapT{"a": 1, "b": 2}},
		{"from cast map", mapT{"a": 1}, Cast(map[string]int{"b": 2}), mapT{"a": 1, "b": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m1.Merge(tt.m2)
			if !tt.m1.Equal(tt.expected) {
				t.Errorf("Merge error. Expected: %q, got %q", tt.expected, tt.m1)
			}
		})
	}
}

func TestMap_MergePanic(t *testing.T) {
	assert.Panics(t, func() {
		var m mapT
		m.Merge(mapT{"a": 1})
	})
}

func TestMap_GetHasLoad(t *testing.T) {
	m := mapT{"a": 1, "b": 2}
	if m.Has("c") {
		t.Errorf("Expected false, got true")
	}
	if !m.Has("b") {
		t.Errorf("Expected true, got false")
	}
	if v := m.Get("b"); v != 2 {
		t.Errorf("Expected 2, got %q", v)
	}
	if v := m.Get("c"); v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}

	v, ok := m.Load("a")
	if v != 1 {
		t.Errorf("Expected 1, got %q", v)
	}
	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestMap_NilGetHasLoad(t *testing.T) {
	var m mapT
	if m.Has("c") {
		t.Errorf("Expected false, got true")
	}
	if v := m.Get("c"); v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}

	v, ok := m.Load("a")
	if v != 0 {
		t.Errorf("Expected 0, got %q", v)
	}
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestMap_Range(t *testing.T) {
	var mNil mapT

	tests := []struct {
		name     string
		m        mapTI
		expected int
	}{
		{"nil", mNil, 0},
		{"1", mapT{"a": 1}, 1},
		{"2", mapT{"a": 1, "b": 2}, 2},
		{"3", mapT{"a": 1, "b": 2, "c": 3}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := 0
			tt.m.Range(func(k string, v int) bool {
				count++
				if count > 1 {
					return false
				}
				return true
			})
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func ExampleCastMap() {
	m := map[string]int{"a": 1}
	b := Cast(m)
	fmt.Print(b.Len())
	//Output: 1
}

func TestStdMap_Set(t *testing.T) {
	var a mapT

	assert.Panics(t, func() {
		a.Set("a", 1)
	})

	a = mapT{}
	a.Set("a", 1)
	a.Set("b", 2)
	assert.Equal(t, 2, a.Get("b"))
}

func TestStdMap_Keys(t *testing.T) {
	m := NewStdMap(mapT{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "c")
}

func TestStdMap_Values(t *testing.T) {
	m := NewStdMap(mapT{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	assert.Len(t, values, 3)
	assert.Contains(t, values, 3)
}

func TestStdMap_Equal(t *testing.T) {
	tests := []struct {
		name string
		m    StdMap[string, int]
		m2   mapT
		want bool
	}{
		{"equal", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 1}, true},
		{"empty", StdMap[string, int]{}, StdMap[string, int]{}, true},
		{"dif len", StdMap[string, int]{"a": 1}, StdMap[string, int]{}, false},
		{"dif len 1", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 1, "b": 2}, false},
		{"dif value", StdMap[string, int]{"a": 1}, StdMap[string, int]{"a": 2}, false},
		{"dif key", StdMap[string, int]{"a": 1}, StdMap[string, int]{"b": 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.m.Equal(tt.m2), "Equal(%v)", tt.m2)
		})
	}
}
