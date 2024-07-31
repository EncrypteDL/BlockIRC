package pkg

import "sync"

type WhoIsList struct {
	sync.RWMutex
	buffer []*WhoIs
	start  int
	end    int
}

type WhoIs struct {
	nickname Name
	username Name
	hostname Name
	hostmask Name
	realname Text
}

func NewWhoIsList(size uint) *WhoIsList {
	return &WhoIsList{
		buffer: make([]*WhoIs, size),
	}
}

func (list *WhoIsList) Append(client *Client) {
	list.Lock()
	defer list.Unlock()
	list.buffer[list.end] = &WhoIs{
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

func (list *WhoIsList) Find(nickname Name, limit int64) []*WhoIs {
	list.RLock()
	defer list.RUnlock()
	results := make([]*WhoIs, 0)
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

func (list *WhoIsList) prev(index int) int {
	list.RLock()
	defer list.RUnlock()
	index -= 1
	if index < 0 {
		index += len(list.buffer)
	}
	return index
}

// Iterate the buffer in reverse.
func (list *WhoIsList) Each() <-chan *WhoIs {
	ch := make(chan *WhoIs)
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
