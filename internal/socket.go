package internal

import (
	"bufio"
	"io"
	"net"
	"sync"

	"github.com/prometheus/common/log"
)

const (
	R = '→'
	W = '←'
)

type Socket struct {
	closed      bool
	closedMutex sync.RWMutex
	conn        net.Conn
	scanner     *bufio.Scanner
	writer      *bufio.Writer
}

func NewSocket(conn net.Conn) *Socket {
	return &Socket{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		writer:  bufio.NewWriter(conn),
	}
}

func (socket *Socket) String() string {
	return socket.conn.RemoteAddr().String()
}

func (socket *Socket) Close() {
	socket.closedMutex.Lock()
	defer socket.closedMutex.Unlock()

	if socket.closed {
		return
	}
	socket.closed = true
	socket.conn.Close()
	log.Debugf("%s closed", socket)
}

func (socket *Socket) Read() (line string, err error) {
	socket.closedMutex.RLock()
	defer socket.closedMutex.RUnlock()
	if socket.closed {
		err = io.EOF
		return
	}

	for socket.scanner.Scan() {
		line = socket.scanner.Text()
		if len(line) == 0 {
			continue
		}
		log.Debugf("%s → %s", socket, line)
		return
	}

	err = socket.scanner.Err()
	socket.isError(err, R)
	if err == nil {
		err = io.EOF
	}
	return
}

func (socket *Socket) Write(line string) (err error) {
	socket.closedMutex.RLock()
	defer socket.closedMutex.RUnlock()
	if socket.closed {
		err = io.EOF
		return
	}

	if _, err = socket.writer.WriteString(line); socket.isError(err, W) {
		return
	}

	if _, err = socket.writer.WriteString(CRLF); socket.isError(err, W) {
		return
	}

	if err = socket.writer.Flush(); socket.isError(err, W) {
		return
	}

	log.Debugf("%s ← %s", socket, line)
	return
}

func (socket *Socket) isError(err error, dir rune) bool {
	if err != nil {
		if err != io.EOF {
			log.Debugf("%s %c error: %s", socket, dir, err)
		}
		return true
	}
	return false
}
