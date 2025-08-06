// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package test_utils

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockedCoreEngine struct {
	mock.Mock
}

func (c *MockedCoreEngine) Fetch(path string, payload any, reqHeaders map[string]string, res any, ctx context.Context) (status int, resHeaders map[string]string, err error) {
	ret := c.Called(path, payload, reqHeaders, res, ctx)
	return ret.Int(0), ret.Get(1).(map[string]string), ret.Error(2)
}
