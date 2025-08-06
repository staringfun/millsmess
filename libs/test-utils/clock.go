// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package test_utils

import "time"

type MockedClock struct {
	Value time.Time
}

func (c *MockedClock) UTCTime() time.Time {
	return c.Value
}
