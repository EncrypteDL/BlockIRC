package internal

import "sync"

type WhoWasList struct {
	sync.RWMutex
	buffer []*WhoWas
	start  int
	end    int
}

type WhoWas struct {
	nickname Name
	username Name
	hostname Name
	hostmask Name
	realname Text
}

func NewWhoWasList(size uint) *WhoWasList {
	return &WhoWasList{
		buffer: make([]*WhoWas, size),
	}
}

func (list *WhoWasList) Append(client *Client) {
	list.Lock()
	defer list.Unlock()
	list.buffer[list.end] = &WhoWas{
		nickname: client.Nick(),
		username: client.username,
		hostname: client.hostname,
		hostmask: client.hostmask,
		realname: client.realname,
	}
	list.end = (list.end + 1) % len(list.buffer)
	if list.end == list.start {
		list.start = (list.end + 1) % len(list.buffer)
	}
}

func (list *WhoWasList) Find(nickname Name, limit int64) []*WhoWas {
	list.RLock()
	defer list.RUnlock()
	results := make([]*WhoWas, 0)
	for whoWas := range list.Each() {
		if nickname != whoWas.nickname {
			continue
		}
		results = append(results, whoWas)
		if int64(len(results)) >= limit {
			break
		}
	}
	return results
}

func (list *WhoWasList) prev(index int) int {
	list.RLock()
	defer list.RUnlock()
	index -= 1
	if index < 0 {
		index += len(list.buffer)
	}
	return index
}

// Iterate the buffer in reverse.
func (list *WhoWasList) Each() <-chan *WhoWas {
	ch := make(chan *WhoWas)
	go func() {
		list.RLock()
		defer list.RUnlock()
		defer close(ch)
		if list.start == list.end {
			return
		}
		start := list.prev(list.end)
		end := list.prev(list.start)
		for start != end {
			ch <- list.buffer[start]
			start = list.prev(start)
		}
	}()
	return ch
}
