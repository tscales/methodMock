package methodMock

import (
	"bytes"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type Mock struct {
	ExpectedCalls []*Call
	mutex         sync.Mutex
}

func (m *Mock) On(methodName string, arguments ...interface{}) *Call {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	c := newCall(m, methodName, arguments...)
	m.ExpectedCalls = append(m.ExpectedCalls, c)
	return c
}

func (m *Mock) Called(arguments ...interface{}) []interface{} {
	// get the method called just before this one on the stack
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := runtime.FuncForPC(pc).Name()
	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]

	return m.MethodCalled(functionName, arguments...)
}

func (m *Mock) MethodCalled(methodName string, arguments ...interface{}) []interface{} {
	m.mutex.Lock()
	call := m.findExpectedCall(methodName, arguments...)
	m.mutex.Unlock()
	if call != nil {
		m.mutex.Lock()
		returnArgs := call.ReturnArguments
		m.mutex.Unlock()

		return returnArgs
	} else {
		panic("unexpected method call")
	}
}

func (m *Mock) findExpectedCall(method string, arguments ...interface{}) *Call {
	var expectedCall *Call

	for _, call := range m.ExpectedCalls {
		if call.Method == method {
			diffCount := argumentDiff(call.Arguments, arguments)
			if diffCount == 0 {
				expectedCall = call
			}
		}
	}
	return expectedCall
}

// argumentDiff checks that the arguments for a defined mock call
// match the arguments received from a mock instance
func argumentDiff(callArgs []interface{}, arguments []interface{}) int {
	var differences int
	maxArgCount := len(callArgs)
	if len(arguments) > maxArgCount {
		maxArgCount = len(arguments)
	}
	for i := 0; i < maxArgCount; i++ {
		var actual, expected interface{}
		if len(arguments) <= i {
			actual = "(Missing)"
		} else {
			actual = arguments[i]
		}

		if len(callArgs) <= i {
			expected = "(Missing)"
		} else {
			expected = callArgs[i]
		}
		if !objectsAreEqual(actual, expected) {
			differences++
		}
	}
	return differences
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

type Call struct {
	Parent          *Mock
	Method          string
	Arguments       []interface{}
	ReturnArguments []interface{}
}

func (c *Call) lock() {
	c.Parent.mutex.Lock()
}

func (c *Call) unlock() {
	c.Parent.mutex.Unlock()
}

func (c *Call) Return(returnArguments ...interface{}) *Call {
	c.lock()
	defer c.unlock()

	c.ReturnArguments = returnArguments
	return c
}

func newCall(parent *Mock, methodName string, methodArguments ...interface{}) *Call {
	return &Call{
		Parent:          parent,
		Method:          methodName,
		Arguments:       methodArguments,
		ReturnArguments: make([]interface{}, 0),
	}
}
