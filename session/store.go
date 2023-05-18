package session

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"
	"sync"

	gin_sessions "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Store interface {
	gin_sessions.Store
}

func NewStore() Store {
	return &store{
		options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		data: make(map[string]valueType),
	}
}

type valueType map[interface{}]interface{}

type store struct {
	options *sessions.Options
	data    map[string]valueType
	mutex   sync.RWMutex
}

func (m *store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(m, name)
}

func (m *store) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	options := *m.options
	session.Options = &options
	session.IsNew = true

	c, err := r.Cookie(name)
	if err == nil {
		session.ID = c.Value
	} else {
		// Cookie not found, check for the bearer token
		auth := r.Header.Get("Authorization")
		if auth == "" {
			return session, nil
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			return session, nil
		}

		session.ID = token
	}

	v, ok := m.value(session.ID)
	if !ok {
		// No value found in cache, don't set any values in session object,
		// consider a new session
		session.ID = "" // reset session ID to generate a new one
		return session, nil
	}

	// Values found in session, this is not a new session
	session.Values = m.copy(v)
	session.IsNew = false
	return session, nil
}

func (m *store) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	var cookieValue string
	if s.Options.MaxAge < 0 {
		cookieValue = ""
		m.delete(s.ID)
		for k := range s.Values {
			delete(s.Values, k)
		}
	} else {
		if s.ID == "" {
			s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		cookieValue = s.ID
		m.setValue(s.ID, m.copy(s.Values))
	}
	http.SetCookie(w, sessions.NewCookie(s.Name(), cookieValue, s.Options))
	return nil
}

func (m *store) Options(options gin_sessions.Options) {
	m.options = options.ToGorillaOptions()
}

func (m *store) copy(v valueType) valueType {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Encoding to gob failed: %v", err))
	}
	var value valueType
	err = dec.Decode(&value)
	if err != nil {
		panic(fmt.Errorf("could not copy memstore value. Decoding from gob failed: %v", err))
	}
	return value
}

func (m *store) value(name string) (valueType, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	v, ok := m.data[name]
	return v, ok
}

func (m *store) setValue(name string, value valueType) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[name] = value
}

func (m *store) delete(name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, name)
}
