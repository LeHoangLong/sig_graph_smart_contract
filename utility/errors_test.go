package utility_test

import (
	"errors"
	"sig_graph/utility"
	"testing"

	"github.com/stretchr/testify/require"
)

var err = errors.New("abc")

func TestErrorLogging(t *testing.T) {
	testErr := utility.NewError(err)
	require.Contains(t, testErr.String(), "abc")
	require.Contains(t, testErr.String(), "sig_graph/utility/errors_test.go:14") // should contain line where error happens
}

func TestErrorIs(t *testing.T) {
	testErr := utility.NewError(err)
	require.True(t, testErr.Is(err))
}
