# MethodMock

MethodMock is a demonstration of the core component's of the mock library in Testify. It includes minimum necessary To Mock a method.

### Mock

a Mock is simply a collection of Method Calls defined by the user.

```Go
type Mock struct {
	ExpectedCalls []*Call
	mutex         sync.Mutex
}
```

### Call

a Call stores the method name, arguments, and return values of a user defined method call. 
```Go
type Call struct {
	Parent          *Mock
	Method          string
	Arguments       []interface{}
	ReturnArguments []interface{}
}
```

### On() and Return()
`On()` Defines a New Call. `Return()` sets the `ReturnArguments` of a `Call`

### Called()

`Called()` does the work of checking that an invoked method has a defined call and that the arguments match.

### example

```Go
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
```

