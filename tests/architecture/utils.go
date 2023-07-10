package architecture

import (
	"testing"
)

func assertNoError(t *testing.T, mockT *testingT) {
	t.Helper()
	if mockT.errored() {
		t.Fatalf("archtest should not have failed but, %s", mockT.message())
	}
}

type testingT struct {
	errors [][]interface{}
}

func (t *testingT) Error(args ...interface{}) {
	t.errors = append(t.errors, args)
}

func (t testingT) errored() bool {
	return len(t.errors) != 0
}

func (t *testingT) message() interface{} {
	return t.errors[0][0]
}
