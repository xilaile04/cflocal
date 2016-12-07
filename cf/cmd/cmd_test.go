package cmd_test

import (
	"bytes"

	"github.com/sclevine/cflocal/cf/cmd/mocks"

	"github.com/golang/mock/gomock"
)

type mockBufferCloser struct {
	*mocks.MockCloser
	*bytes.Buffer
}

func newMockBufferCloser(ctrl *gomock.Controller, contents ...string) *mockBufferCloser {
	bc := &mockBufferCloser{mocks.NewMockCloser(ctrl), &bytes.Buffer{}}
	for _, v := range contents {
		bc.Buffer.Write([]byte(v))
	}
	return bc
}
