package httpregistry

import "fmt"

// TestingT is the subset of [testing.T] (see also [testing.TB]) used by the httpregistry package.
type TestingT interface {
	Fail()
	Errorf(format string, args ...any)
}

// MockTestingT mocks the testing.T interface and it can be used to assert that test that should fail will fail
type MockTestingT struct {
	HasFailed bool
	Messages  []string
}

// Fail records that the Fail function was called
func (f *MockTestingT) Fail() {
	f.HasFailed = true
}

// Errorf records what error message was emitted
func (f *MockTestingT) Errorf(format string, args ...any) {
	f.Messages = append(f.Messages, fmt.Sprintf(format, args...))
	f.Fail()
}
