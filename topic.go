package demux

import (
	"fmt"
	"io"
	"sync"
)

var (
	rw sync.RWMutex

	topics = make(map[string]*Topic)
)

// Topic is used as a collection of users. You can add Users and Devices to
// Users. Topic has to have a unique name.
type Topic struct {
	Name  string
	users map[string]*User
	rw    sync.RWMutex
}

// User represents a grouping of a list of Devices. User has to have a
// unique name.
type User struct {
	Name    string
	devices map[string]*Device
	rw      sync.RWMutex
}

// Device is a unique writing interface. It has a name and an io.Writer onto
// which messages can be written.
type Device struct {
	UserID   string
	DeviceID string
	Writer   io.Writer
}

// NewTopic creates a topic with empty initialized Users and Devices. Calling
// NewTopic with a name that already exists will result in an error.
func NewTopic(name string) (*Topic, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if topicExists(name) {
		return nil, fmt.Errorf("topic with name %q already exists", name)
	}

	t := &Topic{
		Name:  name,
		users: make(map[string]*User),
	}

	doAddTopic(t)

	return t, nil
}

func topicExists(name string) bool {
	rw.RLock()
	defer rw.RUnlock()
	_, ok := topics[name]

	return ok
}

func doAddTopic(topic *Topic) {
	rw.Lock()
	topics[topic.Name] = topic
	rw.Unlock()
}

// addUser adds a user to the topic. If a user with the name already exists, an error
// is returned.
func (t *Topic) addUser(name string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if t.userExists(name) {
		return nil, fmt.Errorf("user with name %q already exists", name)
	}

	user := &User{
		Name:    name,
		devices: make(map[string]*Device, 0),
	}

	t.doAddUser(user)

	return user, nil
}

func (t *Topic) getUser(name string) *User {
	t.rw.RLock()
	user, ok := t.users[name]
	t.rw.RUnlock()

	if !ok {
		user = &User{
			Name:    name,
			devices: make(map[string]*Device, 0),
		}

		t.doAddUser(user)
	}

	return user
}

func (t *Topic) userExists(name string) bool {
	t.rw.RLock()
	defer t.rw.RUnlock()
	_, ok := t.users[name]

	return ok
}

func (t *Topic) doAddUser(user *User) {
	t.rw.Lock()
	t.users[user.Name] = user
	t.rw.Unlock()
}
