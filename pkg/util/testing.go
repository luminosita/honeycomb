package util

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"runtime/debug"
	"testing"
)

func AssertPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("PANIC %+v\n%s", r, string(debug.Stack()))
	}
}

func ErrEq(x error) gomock.Matcher { return errMatcher{x} }

type errMatcher struct {
	x error
}

func (e errMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(e.x.Error(), x.(error).Error())
}

func (e errMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}
