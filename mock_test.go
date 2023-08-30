package methodMock

import (
	"testing"
)

type counter interface {
	AddOne(n int) int
}

type mockCounter struct {
	Mock
}

func (m *mockCounter) AddOne(n int) int {
	args := m.Called(n)
	if len(args) < 1 {
		panic("error: expected 1 return argument for mocked function AddOne")
	}
	return args[0].(int)
}

func TestMock(t *testing.T) {
	m := &mockCounter{}
	m.On("AddOne", 1).Return(2)

	res := m.AddOne(1)
	if res != 2 {
		t.Fatalf("expected 2, got %d", res)
	}
}
