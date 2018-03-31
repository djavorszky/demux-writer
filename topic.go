package demux

import (
	"fmt"
)

// AddUser adds a user to the topic. If a user with the name already exists, an error
// is returned.
func (t *Topic) AddUser(name string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if t.userExists(name) {
		return nil, fmt.Errorf("user with name %q already exists", name)
	}

	user := &User{
		Name:    name,
		devices: make([]*Device, 0),
	}

	t.doAddUser(user)

	return user, nil
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

func (t *Topic) doDeleteUser(user *User) {
	t.rw.Lock()
	delete(t.users, user.Name)
	t.rw.Unlock()
}
