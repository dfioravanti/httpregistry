package httpregistry

import "fmt"

// TestingT is the subset of [testing.T] (see also [testing.TB]) used by the httpregistry package.
// The reason why this exists is so that we can mock in testa and check if failures happen when we expect.
// By design testing.TB make it impossible for the end user to implement the interface so this is the only way to do so
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
