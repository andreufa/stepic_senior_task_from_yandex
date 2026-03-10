package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

// Дано холодное и горячее хранилище. Необходимо реализовать метод Handle, после вызова которого
// должно быть гарантировано,
// что эти данные могут быть сразу получены. Запись
// в холодное хранилище должно давать гарантию At most
// once
type HotClient interface {
	// 10 мс
	Append(key, value string) error
	Close()
}

type HotClientImpl struct{}

func NewHotClient() HotClientImpl {
	return HotClientImpl{}
}

func (c HotClientImpl) Append(key, value string) error {
	fmt.Printf("hot client | append %v : %v\n", key, value)
	return nil
}
func (c HotClientImpl) Close() {}

type EventAction struct {
	key   string
	value string
}

type FreezeClient interface {
	// 100 мс
	Send(actions []EventAction) error
	Close()
}

type FreezeClientImpl struct{}

func (f FreezeClientImpl) Send(actions []EventAction) error {
	fmt.Printf("FreezeClient send: %v\n", len(actions))
	return nil
}
func (f FreezeClientImpl) Close() {}

func NewFreezeClient() FreezeClientImpl {
	return FreezeClientImpl{}
}

type Handler struct {
	hc        HotClient
	fc        FreezeClient
	batchSize int
	buffer    []EventAction
	in        chan EventAction
	close     chan struct{}
}

func NewHandler(batch int) Handler {
	handler := Handler{
		hc:        NewHotClient(),
		fc:        NewFreezeClient(),
		batchSize: batch,
		buffer:    make([]EventAction, 0, batch),
		in:        make(chan EventAction),
		close:     make(chan struct{}),
	}
	go handler.Worker()
	return handler
}
func (h *Handler) Worker() {
	for {
		select {
		case <-h.close:
			err := h.fc.Send(h.buffer)
			if err != nil {
				log.Println("func Worker |", err)
			}
			return
		case event := <-h.in:
			h.buffer = append(h.buffer, event)
		}
		if len(h.buffer) == h.batchSize {
			err := h.fc.Send(h.buffer)
			if err != nil {
				log.Println("func Worker |", err)
			}
			h.buffer = make([]EventAction, 0, h.batchSize)
		}
	}
}

func (h *Handler) Close() {
	h.hc.Close()
	h.fc.Close()
	close(h.close)
}

func (h *Handler) Handle(key, value string) {
	err := h.hc.Append(key, value)
	if err != nil {
		log.Println("func Handle Append to hot storage |", err)
	}
	h.in <- EventAction{
		key:   key,
		value: value,
	}
}

func main() {
	const (
		BATCH_SIZE = 10
	)
	handle := NewHandler(BATCH_SIZE)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := range 10 {
			handle.Handle(strconv.Itoa(i), "yes")
		}
	}()
	go func() {
		defer wg.Done()
		for i := range 15 {
			handle.Handle(strconv.Itoa(i*100), "yes")
		}
	}()
	wg.Wait()
	handle.Close()

	// handler := NewHandler(BATCH_SIZE)
	// server := http.NewServeMux()
	// server.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost{
	// 		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	key := r.URL.Query().Get("key")
	// 	value := r.URL.Query().Get("value")

	// 	if key == "" || value == "" {
	// 		http.Error(w, "key and value required", http.StatusBadRequest)
	// 		return
	// 	}

	// 	handler.Handle(key,value)
	// 	w.Header().Set("Content-Type", "text/plain")
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprintf(w, "stored: %s = %s", key, value)
	// })
}
