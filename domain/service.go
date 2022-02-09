package domain

import (
	"channels-instagram-dm/domain/model"
	"channels-instagram-dm/domain/model/instagram"
)

const (
	SlotStatusBusy        = "busy"        // Слот занят
	SlotStatusFree        = "free"        // Слот освобожден
	SlotStatusUnavailable = "unavailable" // Слот недоступен, пока не будет проведен discovery (ошибка или logout)
)

type SlotStatus string

type SlotContainer struct {
	Slot     Slot
	Username string
	Metadata SlotMetadata
	Service  InstagramAPI
}

type Slot struct {
	Host   string
	Status SlotStatus
}

type SlotMetadata struct {
	Users      []string
	ActiveUser string
}

type DiscoveryRow struct {
	Host       string   `json:"host"`
	Active     bool     `json:"active"`
	Users      []string `json:"users"`
	ActiveUser string   `json:"active_user"`
}

type Service interface {
	InstagramAPI(username string) (InstagramAPI, error)
	RefreshSlots() error
	Slots() []SlotContainer
}

type InstagramAPI interface {
	Discovery() (DiscoveryRow, error)
	DirectInbox(cursor string, limit int) (instagram.InboxWithThreads, error) // TODO: DirectInbox
	DirectInboxPending(cursor string) ([]instagram.ThreadWithItems, error)
	DirectAcceptInboxPending(threadIDs []string) error
	DirectThread(string, string) (instagram.ThreadWithItems, error) // Add sleep duration
	DirectSendText(username string, text model.Message) error
	RealtimeSendText(threadID string, text model.Message) error
	Login(credentials instagram.Credentials) (instagram.Required, error)
	Login2F(credentials instagram.Credentials, required instagram.Required) error
	Challenge(required instagram.Required) error
	Logout() error
	ListenThreadUpdates(chan instagram.RealtimeUpdate) (chan error, error)
	Close()
	IsClosed() bool
}
