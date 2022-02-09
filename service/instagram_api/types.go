package instagram_api

import (
	"strconv"

	"channels-instagram-dm/domain/model/instagram"
	"channels-instagram-dm/service/instagram_api/types"
)

type InboxResponse struct {
	Inbox                `json:"inbox"`
	Viewer               types.User `json:"viewer"`
	SeqID                string     `json:"seq_id"`
	SnapshotAtMs         string     `json:"snapshot_at_ms"`
	PendingRequestsTotal int        `json:"pending_requests_total"` // TODO: Вот этот параметр можно смотреть, чтобы пендить только когда надо
	Status               string     `json:"status"`
}

type Inbox struct {
	Threads               []Thread `json:"threads"`
	HasNewer              bool     `json:"has_newer"`
	HasOlder              bool     `json:"has_older"`
	OldestCursor          string   `json:"oldest_cursor"`
	UnseenCount           int      `json:"unseen_count"`
	UnseenCountTs         string   `json:"unseen_count_ts"`
	HasPendingTopRequests bool     `json:"has_pending_top_requests"`
	PendingRequestsTotal  int      `json:"pending_requests_total"`
	// BlendedInboxEnabled bool   `json:"blended_inbox_enabled"`
}

type Thread struct {
	ID             string       `json:"thread_id"`
	V2ID           string       `json:"thread_v2_id"`
	Items          []ThreadItem `json:"items"`
	LastActivityAt string       `json:"last_activity_at"`
	Muted          bool         `json:"muted"`
	IsPin          bool         `json:"is_pin"`
	Named          bool         `json:"named"`
	Pending        bool         `json:"pending"`
	Archived       bool         `json:"archived"`
	ThreadType     string       `json:"thread_type"`
	ViewerID       interface{}  `json:"viewer_id"`
	// Title                string       `json:"thread_title"`
	// Folder               uint         `json:"folder"`
	// HasDropIn            bool         `json:"thread_has_drop_in"`
	// IsGroup              bool         `json:"is_group"`
	// BusinessThreadFolder uint         `json:"business_thread_folder"`
	// ReadState            uint         `json:"read_state"`
	IsVerified   bool         `json:"is_verified_thread"`
	HasOlder     bool         `json:"has_older"`
	HasNewer     bool         `json:"has_newer"`
	NewestCursor string       `json:"newest_cursor"`
	OldestCursor string       `json:"oldest_cursor"`
	NextCursor   string       `json:"next_cursor"`
	PrevCursor   string       `json:"prev_cursor"`
	Inviter      types.User   `json:"inviter"`
	Users        []types.User `json:"users"`
	// LastSeenAt           map[string]struct {
	// 	Timestamp string `json:"timestamp"`
	// 	ItemID    string `json:"item_id"`
	// } `json:"last_seen_at"`
	LastPermanentItem struct {
		ItemID        interface{} `json:"item_id"`
		UserID        interface{} `json:"user_id"`
		Timestamp     string      `json:"timestamp"`
		ItemType      string      `json:"item_type"`
		ClientContext string      `json:"client_context"`
	} `json:"last_permanent_item"`
	TqSeqId int `json:"tq_seq_id"`
	UqSeqId int `json:"uq_seq_id"`
}

type RealtimeUpdate struct {
	ThreadID     string     `json:"thread_id"`
	ThreadItemID string     `json:"thread_item_id"`
	ThreadItem   ThreadItem `json:"thread_item"`
}

type ThreadItem struct {
	ID            string      `json:"item_id"`
	UserID        string      `json:"user_id"`
	Timestamp     int64       `json:"timestamp"`
	ClientContext string      `json:"client_context"`
	Type          string      `json:"item_type"`
	Text          types.Text  `json:"text"`
	Like          types.Like  `json:"like"`
	Media         types.Media `json:"media"`
	Link          types.Link
	ActionLog     types.ActionLog     `json:"action_log"`
	VisualMedia   types.VisualMedia   `json:"visual_media"`
	AnimatedMedia types.AnimatedMedia `json:"animated_media"`
	VoiceMedia    types.VoiceMedia    `json:"voice_media"`
	MediaShare    types.MediaShare    `json:"media_share"`
	StoryShare    types.StoryShare    `json:"story_share"`
	ReelShare     types.ReelShare     `json:"reel_share"`
	Clip          types.Clip          `json:"clip"`
	Profile       types.Profile       `json:"profile"`
}

type LoginRequired struct {
	Required string            `json:"required"`
	Data     LoginRequiredData `json:"data,omitempty"`
}

type LoginRequiredData struct {
	Identifier    string `json:"identifier"`
	Step          string `json:"step"`
	CheckpointUrl string `json:"checkpoint_url"`
	Method        string `json:"method"`
	Code          string `json:"code"`
}

func (i InboxResponse) toModel() (instagram.InboxWithThreads, error) {
	inboxModel := instagram.InboxWithThreads{
		Inbox: instagram.Inbox{
			HasNewer:              i.Inbox.HasNewer,
			HasOlder:              i.Inbox.HasOlder,
			OldestCursor:          i.Inbox.OldestCursor,
			UnseenCount:           i.Inbox.UnseenCount,
			HasPendingTopRequests: i.HasPendingTopRequests,
			PendingRequestsTotal:  i.PendingRequestsTotal,
		},
		Threads: make([]instagram.ThreadWithItems, 0, len(i.Inbox.Threads)),
	}

	unseenCountTs, err := strconv.ParseInt(i.Inbox.UnseenCountTs, 10, 64)
	if err != nil {
		return instagram.InboxWithThreads{}, err
	}
	inboxModel.Inbox.UnseenCountTs = unseenCountTs

	seqID, err := strconv.ParseInt(i.SeqID, 10, 64)
	if err != nil {
		return instagram.InboxWithThreads{}, err
	}
	inboxModel.Inbox.SeqID = seqID

	snapshotAtMs, err := strconv.ParseInt(i.SnapshotAtMs, 10, 64)
	if err != nil {
		return instagram.InboxWithThreads{}, err
	}
	inboxModel.Inbox.SnapshotAt = snapshotAtMs * 1000

	for _, thread := range i.Inbox.Threads {
		threadModel, err := thread.toModel()
		if err != nil {
			return instagram.InboxWithThreads{}, err
		}

		inboxModel.Threads = append(inboxModel.Threads, threadModel)
	}

	if err := inboxModel.Inbox.Validate(); err != nil {
		return instagram.InboxWithThreads{}, err
	}

	return inboxModel, nil
}

func (t Thread) toModel() (instagram.ThreadWithItems, error) {
	threadModel := instagram.ThreadWithItems{
		Thread: instagram.Thread{
			ID:            t.ID,
			V2ID:          t.V2ID,
			Pending:       t.Pending,
			Archived:      t.Archived,
			ThreadType:    t.ThreadType,
			HasOlder:      t.HasOlder,
			HasNewer:      t.HasNewer,
			NewestCursor:  t.NewestCursor,
			OldestCursor:  t.OldestCursor,
			ViewerUserID:  types.ValueToString(t.ViewerID),
			InviterUserID: types.ValueToString(t.Inviter.ID),
			LastPermanentItem: struct {
				ItemID    string
				UserID    string
				Timestamp int64
				ItemType  string
			}{
				ItemID:   types.ValueToString(t.LastPermanentItem.ItemID),
				UserID:   types.ValueToString(t.LastPermanentItem.UserID),
				ItemType: t.LastPermanentItem.ItemType,
			},
		},
		Items: make([]instagram.ThreadItem, 0, len(t.Items)),
		Users: make([]instagram.User, 0, len(t.Users)),
	}

	lastActivityAt, err := strconv.ParseInt(t.LastActivityAt, 10, 64)
	if err != nil {
		return instagram.ThreadWithItems{}, err
	}
	threadModel.Thread.LastActivityAt = lastActivityAt

	if t.LastPermanentItem.Timestamp != "" {
		timestamp, err := strconv.ParseInt(t.LastPermanentItem.Timestamp, 10, 64)
		if err != nil {
			return instagram.ThreadWithItems{}, err
		}
		threadModel.LastPermanentItem.Timestamp = timestamp
	}

	for _, item := range t.Items {
		itemModel, err := item.toModel()
		if err != nil {
			return instagram.ThreadWithItems{}, err
		}

		threadModel.Items = append(threadModel.Items, itemModel)
	}

	for _, user := range t.Users {
		userModel, err := user.ToModel()
		if err != nil {
			return instagram.ThreadWithItems{}, err
		}

		threadModel.Users = append(threadModel.Users, userModel)
	}

	if err := threadModel.Thread.Validate(); err != nil {
		return instagram.ThreadWithItems{}, err
	}

	return threadModel, nil
}

func (t ThreadItem) toModel() (instagram.ThreadItem, error) {
	threadModel := instagram.ThreadItem{
		ID:            t.ID,
		UserID:        t.UserID,
		Timestamp:     t.Timestamp,
		ClientContext: t.ClientContext,
		Type:          instagram.MessageTypeUndefined,
	}

	switch t.Type {
	case "text":
		model, err := t.Text.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeText
		threadModel.Text = model

	case "like":
		model, err := t.Like.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeLike
		threadModel.Text = model

	case "action_log":
		model, err := t.ActionLog.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeActionLog
		threadModel.Text = model

	case "link":
		model, err := t.Link.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeLink
		threadModel.Link = model

	case "media":
		model, err := t.Media.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "raven_media":
		model, err := t.VisualMedia.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "animated_media":
		model, err := t.AnimatedMedia.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "voice_media":
		model, err := t.VoiceMedia.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "media_share":
		model, err := t.MediaShare.Media.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "story_share":
		model, err := t.StoryShare.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "clip":
		model, err := t.Clip.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeMedia
		threadModel.Media = model

	case "profile":
		model, err := t.Profile.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeLink
		threadModel.Link = model

	case "reel_share":
		model, err := t.ReelShare.ToModel()
		if err != nil {
			return instagram.ThreadItem{}, err
		}

		threadModel.Type = instagram.MessageTypeLink
		threadModel.Link = model

	default:
		threadModel.Type = instagram.MessageTypeUndefined
	}

	if err := threadModel.Validate(); err != nil {
		return instagram.ThreadItem{}, err
	}

	return threadModel, nil
}

func (l LoginRequired) toModel() instagram.Required {
	loginModel := instagram.Required{}
	switch l.Required {
	case string(instagram.RequiredStep2F):
		loginModel.Case = instagram.RequiredStep2F
		loginModel.Options = instagram.RequiredOptions{
			Identifier: l.Data.Identifier,
			Method:     l.Data.Method,
			Code:       l.Data.Code,
		}
	case string(instagram.RequiredStepChallenge):
		loginModel.Case = instagram.RequiredStepChallenge
		loginModel.Options = instagram.RequiredOptions{
			Step:          l.Data.Step,
			CheckpointUrl: l.Data.CheckpointUrl,
			Code:          l.Data.Code,
		}
	default:
		loginModel.Case = instagram.RequiredStepNone
	}

	return loginModel
}

func (u RealtimeUpdate) toModel() (instagram.RealtimeUpdate, error) {
	rtModel := instagram.RealtimeUpdate{
		ThreadID:     u.ThreadID,
		ThreadItemID: u.ThreadItemID,
	}

	itemModel, err := u.ThreadItem.toModel()
	if err != nil {
		return instagram.RealtimeUpdate{}, err
	}

	rtModel.ThreadItem = itemModel

	return rtModel, nil
}
