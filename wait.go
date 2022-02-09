package main

import (
	"sync"

	"channels-instagram-dm/domain"
)

type syncer struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func Syncer() domain.Syncer {
	return &syncer{
		ch: make(chan struct{}),
		wg: &sync.WaitGroup{},
	}
}

func (s *syncer) Add() {
	s.wg.Add(1)
}

func (s *syncer) Remove() {
	s.wg.Done()
}

func (s *syncer) Done() <-chan struct{} {
	go func() {
		s.wg.Wait()
		s.ch <- struct{}{}
	}()

	return s.ch
}
