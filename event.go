package main

import (
	"channels-instagram-dm/domain"
)

type eventBus struct {
	subscribers []subscriberOn
}

type subscriberOn struct {
	channel interface{}
	filter  domain.EventFilter
}

func EventBus() domain.EventBus {
	return &eventBus{
		subscribers: make([]subscriberOn, 0, 100),
	}
}

func (e *eventBus) PublishAccountCreated(event domain.EventAccountCreated) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventAccountCreated); ok {
				go func(ch chan domain.EventAccountCreated) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishAccountResumed(event domain.EventAccountResumed) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventAccountResumed); ok {
				go func(ch chan domain.EventAccountResumed) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishAccountSuspended(event domain.EventAccountSuspended) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventAccountSuspended); ok {
				go func(ch chan domain.EventAccountSuspended) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishAccountDeleted(event domain.EventAccountDeleted) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventAccountDeleted); ok {
				go func(ch chan domain.EventAccountDeleted) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishAccountLogout(event domain.EventAccountLogout) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventAccountLogout); ok {
				go func(ch chan domain.EventAccountLogout) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishSuspendAccount(event domain.EventSuspendAccount) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventSuspendAccount); ok {
				go func(ch chan domain.EventSuspendAccount) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishInboxHasChanges(event domain.EventInboxHasChanges) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventInboxHasChanges); ok {
				go func(ch chan domain.EventInboxHasChanges) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) PublishLoginAccount(event domain.EventLoginAccount) {
	for _, s := range e.subscribers {
		if s.filter == nil || s.filter(event) {
			if channel, ok := s.channel.(chan domain.EventLoginAccount); ok {
				go func(ch chan domain.EventLoginAccount) {
					ch <- event
				}(channel)
			}
		}
	}
}

func (e *eventBus) SubscribeOn(channel interface{}, f domain.EventFilter) {
	e.subscribers = append(e.subscribers, subscriberOn{
		channel: channel,
		filter:  f,
	})
}

func (e *eventBus) SubscribeOnAccountCreated(f domain.EventFilter) chan domain.EventAccountCreated {
	ch := make(chan domain.EventAccountCreated)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnAccountResumed(f domain.EventFilter) chan domain.EventAccountResumed {
	ch := make(chan domain.EventAccountResumed)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnAccountSuspended(f domain.EventFilter) chan domain.EventAccountSuspended {
	ch := make(chan domain.EventAccountSuspended)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnAccountDeleted(f domain.EventFilter) chan domain.EventAccountDeleted {
	ch := make(chan domain.EventAccountDeleted)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnAccountLogout(f domain.EventFilter) chan domain.EventAccountLogout {
	ch := make(chan domain.EventAccountLogout)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnSuspendAccount(f domain.EventFilter) chan domain.EventSuspendAccount {
	ch := make(chan domain.EventSuspendAccount)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnInboxHasChanges(f domain.EventFilter) chan domain.EventInboxHasChanges {
	ch := make(chan domain.EventInboxHasChanges)
	e.SubscribeOn(ch, f)

	return ch
}

func (e *eventBus) SubscribeOnLoginAccount(f domain.EventFilter) chan domain.EventLoginAccount {
	ch := make(chan domain.EventLoginAccount)
	e.SubscribeOn(ch, f)

	return ch
}
