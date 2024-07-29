package mercure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"sync"
)

func init() {
	RegisterTransportFactory("mysql", NewMySQLTransport)
}

type MySQLTransport struct {
	sync.RWMutex
	subscribers *SubscriberList
	db          *sql.DB
	closed      chan struct{}
	lastEventID string
}

func NewMySQLTransport(u *url.URL, l Logger) (Transport, error) {
	dsn := u.Query().Get("dsn")
	if dsn == "" {
		return nil, fmt.Errorf("MySQL DSN is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLTransport{
		db:          db,
		subscribers: NewSubscriberList(1e5),
		closed:      make(chan struct{}),
	}, nil
}

func (m *MySQLTransport) Dispatch(update *Update) error {
	select {
	case <-m.closed:
		return ErrClosedTransport
	default:
	}
	serializedUpdate, err := json.Marshal(update)
	if err != nil {
		return err
	}
	_, err = m.db.Exec("INSERT INTO updates (data) VALUES (?)", serializedUpdate)
	if err != nil {
		return err
	}

	m.RLock()
	defer m.RUnlock()

	for _, s := range m.subscribers.MatchAny(update) {
		s.Dispatch(update, false)
	}

	return nil
}

func (m *MySQLTransport) AddSubscriber(s *Subscriber) error {
	select {
	case <-m.closed:
		return ErrClosedTransport
	default:
	}

	m.Lock()
	defer m.Unlock()

	m.subscribers.Add(s)

	s.Ready()

	return nil
}

func (m *MySQLTransport) RemoveSubscriber(s *Subscriber) error {
	select {
	case <-m.closed:
		return ErrClosedTransport
	default:
	}

	m.Lock()
	defer m.Unlock()

	m.subscribers.Remove(s)

	return nil
}

func (m *MySQLTransport) GetSubscribers() (string, []*Subscriber, error) {
	m.RLock()
	defer m.RUnlock()

	return m.lastEventID, getSubscribers(m.subscribers), nil
}

func (m *MySQLTransport) Close() error {
	return m.db.Close()
}
