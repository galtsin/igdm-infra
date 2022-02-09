package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/service/instagram_api"
)

const (
	RefreshSlotsTimer = 5 * time.Minute
)

type factory struct {
	ctx      context.Context
	mx       *sync.Mutex
	logger   domain.Logger
	slotsURI string
	slots    []domain.SlotContainer
}

func Factory(ctx context.Context, logger domain.Logger, slotsURI string) (domain.Service, error) {
	f := &factory{
		ctx:      ctx,
		slotsURI: slotsURI,
		logger:   logger,
		mx:       &sync.Mutex{},
	}

	if err := f.RefreshSlots(); err != nil {
		return nil, fmt.Errorf("ScanSlots: Error %s", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.NewTicker(RefreshSlotsTimer).C:
				f.logger.Debug("RefreshSlots: Processing", nil)

				if err := f.RefreshSlots(); err != nil {
					f.logger.Error(fmt.Sprintf("RefreshSlots: ScanSlots, Error %s", err), nil)
				}
			}
		}
	}()

	return f, nil
}

func (f *factory) RefreshSlots() error {
	f.mx.Lock()
	defer f.mx.Unlock()

	hosts, err := f.readSlots()
	if err != nil {
		return err
	}

	if len(hosts) == 0 {
		f.logger.Info("RefreshSlots: Hosts is empty", nil)
	}

	if err := f.scanSlots(hosts); err != nil {
		return err
	}

	f.discoverySlots()

	return nil
}

func (f *factory) readSlots() ([]string, error) {
	u, err := url.Parse(f.slotsURI)
	if err != nil {
		return []string{}, err
	}

	var reader io.Reader

	switch u.Scheme {
	case "http", "https":
		resp, err := http.Get(u.String())
		if err != nil {
			return []string{}, err
		}

		reader = resp.Body
	default:
		file, err := os.Open(u.String())
		if err != nil {
			return []string{}, err
		}

		reader = file
	}

	if reader == nil {
		return []string{}, nil
	}

	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				f.logger.Error(fmt.Sprintf("ReadSlots: Error %s", err), nil)
			}
		}
	}()

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return []string{}, err
	}

	if len(bs) == 0 {
		return []string{}, nil
	}

	data := string(bs[:])
	rows := strings.Split(data, "\n")

	hosts := make([]string, 0, len(rows))

	for _, row := range rows {
		if row == "" {
			continue
		}

		hosts = append(hosts, strings.TrimSpace(row))
	}

	return hosts, nil
}

// Attention: Операция должна использоваться только в атомарном вызове f.mx.Lock()
func (f *factory) scanSlots(hosts []string) error {
	slots := make([]domain.Slot, 0, len(hosts))

	for _, host := range hosts {
		slots = append(slots, domain.Slot{
			Host:   host,
			Status: domain.SlotStatusFree,
		})
	}

	slotContainers := make([]domain.SlotContainer, 0, len(slots))

	// Удаляем неиспользованные слоты
REMOVE_LIST:
	for _, sc := range f.slots {
		for _, s := range slots {
			// Копируем совпадающие слоты
			if sc.Slot.Host == s.Host {
				slotContainers = append(slotContainers, sc)
				continue REMOVE_LIST
			}
		}

		// Удаляем, если слот свободен
		if sc.Slot.Status == domain.SlotStatusFree {
			f.logger.Info(fmt.Sprintf("ScanSlots: Slot %s was revoked", sc.Slot.Host), nil)
			continue
		}

		// Удаляем, если слот недоступен
		if sc.Slot.Status == domain.SlotStatusUnavailable {
			f.logger.Info(fmt.Sprintf("ScanSlots: Slot %s was revoked", sc.Slot.Host), nil)
			continue
		}

		// Не трогаем занятые слоты
		if sc.Slot.Status == domain.SlotStatusBusy {
			slotContainers = append(slotContainers, sc)
			continue
		}
	}

APPEND_LIST:
	for _, s := range slots {
		for _, sc := range f.slots {
			if s.Host == sc.Slot.Host {
				continue APPEND_LIST
			}
		}

		// Добавляем новые слоты
		slotContainers = append(slotContainers, domain.SlotContainer{
			Slot:     s,
			Metadata: domain.SlotMetadata{},
			Username: "",
			Service:  nil,
		})

		f.logger.Info(fmt.Sprintf("ScanSlots: Slot %s was appended", s.Host), nil)
	}

	sort.Slice(slotContainers, func(i, j int) bool {
		return len(slotContainers[i].Metadata.Users) < len(slotContainers[j].Metadata.Users)
	})

	f.slots = slotContainers

	return nil
}

// Attention: Операция должна использоваться только в атомарном вызове f.mx.Lock()
func (f *factory) discoverySlots() {
	// Инорфмация из discovery дополняющая и не должна противоречить данным из takeService
	for i, sc := range f.slots {
		if sc.Slot.Status == domain.SlotStatusBusy {
			if sc.Service.IsClosed() {
				// Слот либо закрыт из-за ошибки, либо освобожден через logout
				sc.Slot.Status = domain.SlotStatusUnavailable
				sc.Service = nil
			}
		}

		if sc.Service == nil {
			service, err := f.createService(sc.Slot)
			if err != nil {
				sc.Slot.Status = domain.SlotStatusUnavailable
				sc.Service = nil

				// Фиксируем изменения
				f.slots[i] = sc

				f.logger.Error(fmt.Sprintf("CreateService: Slot %s, Error. %s", sc.Slot.Host, err), nil)
				continue
			}

			sc.Service = service
		}

		// Исследуются только слоты в статусе SlotStatusUnavailable и SlotStatusFree
		discovered, err := sc.Service.Discovery()
		if err != nil {
			sc.Service.Close()

			sc.Slot.Status = domain.SlotStatusUnavailable
			sc.Service = nil

			// Фиксируем изменения
			f.slots[i] = sc

			f.logger.Error(fmt.Sprintf("DiscoveryService: Slot %s, Error. %s", sc.Slot.Host, err), nil)
			continue
		}

		sc.Metadata = domain.SlotMetadata{
			Users:      discovered.Users,
			ActiveUser: discovered.ActiveUser,
		}

		if sc.Slot.Status == domain.SlotStatusFree {
			// Фиксируем изменения
			f.slots[i] = sc
			continue
		}

		if sc.Slot.Status == domain.SlotStatusUnavailable {
			// Слот был освобожден через logout или перезапущен. Зачищаем окончательно
			if discovered.ActiveUser == "" {
				sc.Username = ""
				sc.Slot.Status = domain.SlotStatusFree

				// Фиксируем изменения
				f.slots[i] = sc
				continue
			}

			// Слот был с ошибкой, т.к. аккаунт не разлогинило. Возвращаем его аккаунту
			if discovered.ActiveUser != "" {
				// Сбой. Не совпадает аккаунт владеющий слотом и физически занявший его
				if sc.Username != discovered.ActiveUser {
					sc.Slot.Status = domain.SlotStatusUnavailable

					// Фиксируем изменения
					f.slots[i] = sc

					f.logger.Error(fmt.Sprintf("DiscoveryService: Slot %s, Users is mismatch, want %s, got %s", sc.Slot.Host, sc.Username, discovered.ActiveUser), nil)
					continue
				}

				sc.Slot.Status = domain.SlotStatusBusy

				// Фиксируем изменения
				f.slots[i] = sc
			}
		}

		// Фиксируем изменения
		f.slots[i] = sc
	}
}

func (f *factory) Slots() []domain.SlotContainer {
	return f.slots
}

func (f *factory) InstagramAPI(username string) (domain.InstagramAPI, error) {
	service := f.takeService(username)
	if service == nil {
		return nil, fmt.Errorf("No available servers")
	}

	if service.IsClosed() {
		if err := f.RefreshSlots(); err != nil {
			f.logger.Error(fmt.Sprintf("RefreshSlots: Error %s", err), nil)
		}
	}

	return service, nil
}

func (f *factory) createService(slot domain.Slot) (domain.InstagramAPI, error) {
	return instagram_api.NewService(f.ctx, f.logger.Copy(fmt.Sprintf("(host=%s)", slot.Host)), slot.Host)
}

func (f *factory) takeService(username string) domain.InstagramAPI {
	f.mx.Lock()
	defer f.mx.Unlock()

	// Ищем свой активный слот
	for _, sc := range f.slots {
		if sc.Username == username {
			// Аккаунт не должен занимать соседних слотов, в случае проблем с текущим
			// Иначе может возникнуть ситуация, когда неисправный аккаунт способен забрать все свободные слоты
			// Решение проблемы: отозвать слот
			if sc.Slot.Status == domain.SlotStatusUnavailable {
				return nil
			}

			return sc.Service
		}
	}

	// Ищем наиболее подходящий слот
	for i, sc := range f.slots {
		if sc.Slot.Status == domain.SlotStatusUnavailable {
			continue
		}

		if sc.Slot.Status == domain.SlotStatusBusy {
			continue
		}

		for _, user := range sc.Metadata.Users {
			if user == username {
				f.slots[i].Slot.Status = domain.SlotStatusBusy
				f.slots[i].Username = username
				return sc.Service
			}
		}
	}

	// Ищем слот из совсех свободных
	for i, sc := range f.slots {
		if sc.Slot.Status == domain.SlotStatusUnavailable {
			continue
		}

		if sc.Slot.Status == domain.SlotStatusBusy {
			continue
		}

		if len(sc.Metadata.Users) == 0 {
			f.slots[i].Slot.Status = domain.SlotStatusBusy
			f.slots[i].Username = username

			return sc.Service
		}
	}

	// Ищем слот с минимальным количеством аккаунтов
	sort.Slice(f.slots, func(i, j int) bool {
		return len(f.slots[i].Metadata.Users) < len(f.slots[j].Metadata.Users)
	})

	for i, sc := range f.slots {
		if sc.Slot.Status == domain.SlotStatusUnavailable {
			continue
		}

		if sc.Slot.Status == domain.SlotStatusBusy {
			continue
		}

		f.slots[i].Slot.Status = domain.SlotStatusBusy
		f.slots[i].Username = username

		return sc.Service
	}

	return nil
}
