package syncmanager

import (
	"math/rand"
	"testing"
	"time"
)

func TestSyncManager(t *testing.T) {
	SyncManager.Init(4)

	for i := 0; i < 10; i++ {
		SyncManager.AddQueue(func(...any) {
			random := rand.Intn(4)
			duration := time.Second * time.Duration(random)
			time.Sleep(duration)
			//println("Hello ", random)
		})
	}

	SyncManager.WaitFinish()
}
