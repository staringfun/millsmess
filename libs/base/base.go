// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import "context"

type Base struct {
	Clock         Clock
	MemoryStorage MemoryStorage
	Ctx           context.Context
}
