// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import "time"

type Clock interface {
	UTCTime() time.Time
}

type DefaultClock struct{}

func (*DefaultClock) UTCTime() time.Time {
	return time.Now().UTC()
}
