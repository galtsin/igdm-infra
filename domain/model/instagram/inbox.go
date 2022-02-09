package instagram

import "fmt"

type Inbox struct {
	HasNewer              bool
	HasOlder              bool
	OldestCursor          string
	UnseenCount           int
	UnseenCountTs         int64
	SeqID                 int64
	HasPendingTopRequests bool
	PendingRequestsTotal  int
	SnapshotAt            int64
}

type InboxWithThreads struct {
	Inbox
	Threads []ThreadWithItems
}

func (i Inbox) Validate() error {
	if i.SeqID == 0 {
		return fmt.Errorf("SeqID should not be empty")
	}

	if i.SnapshotAt == 0 {
		return fmt.Errorf("SnapshotAt should not be empty")
	}

	return nil
}
