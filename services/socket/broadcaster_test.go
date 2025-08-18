package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"testing"
)

func TestParallelBroadcasters(t *testing.T) {
	broadcasterGroup := NewBroadcasterGroup()

	connectionsCount := 100000
	connections := make([]base.ContextWriter, connectionsCount)

	groupsCount := 100

	eg, _ := errgroup.WithContext(context.Background())
	for i, connection := range connections {
		eg.Go(func() error {
			broadcasterGroup.Join(i%groupsCount, i, connection)
			return nil
		})
		eg.Go(func() error {
			broadcasterGroup.Join((i+1)%groupsCount, i, connection)
			return nil
		})
		if i > 0 {
			eg.Go(func() error {
				broadcasterGroup.Join((i-1)%groupsCount, i, connection)
				return nil
			})
		}
	}
	err := eg.Wait()
	require.Nil(t, err, "wait error")

	for i, connection := range connections {
		eg.Go(func() error {
			if i%2 == 0 {
				broadcasterGroup.Join(i%groupsCount, i, connection)
			} else {
				broadcasterGroup.Leave(i%groupsCount, i)
			}
			return nil
		})
	}
	err = eg.Wait()
	require.Nil(t, err, "wait error")

	for i := range connections {
		eg.Go(func() error {
			broadcasterGroup.LeaveAll(i)
			broadcasterGroup.LeaveAll(i)
			broadcasterGroup.LeaveAll(i)
			return nil
		})
	}
	err = eg.Wait()
	require.Nil(t, err, "wait error")

	hasBroadcasters := false
	broadcasterGroup.broadcasters.Range(func(key, value any) bool {
		hasBroadcasters = true
		return false
	})
	require.False(t, hasBroadcasters, "broadcasters should be empty")

	hasPlayerBroadcasters := false
	broadcasterGroup.itemBroadcasters.Range(func(key, value any) bool {
		hasPlayerBroadcasters = true
		return false
	})
	require.False(t, hasPlayerBroadcasters, "item broadcasters should be empty")
}
