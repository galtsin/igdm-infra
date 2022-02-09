package instagram

import "fmt"

type Thread struct {
	ID                string
	V2ID              string
	LastActivityAt    int64
	Pending           bool
	Archived          bool
	ThreadType        string
	HasOlder          bool
	HasNewer          bool
	NewestCursor      string // TODO: Если это предыдущий threadID то мы можем иметь четкую последовательность сообщений
	OldestCursor      string
	ViewerUserID      string
	InviterUserID     string
	LastPermanentItem struct {
		ItemID    string
		UserID    string
		Timestamp int64
		ItemType  string
	}
}

type ThreadWithItems struct {
	Thread
	Items []ThreadItem
	Users []User
}

func (t Thread) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("ID should not be empty")
	}

	if t.LastActivityAt == 0 {
		return fmt.Errorf("LastActivityAt should not be empty")
	}

	return nil
}
