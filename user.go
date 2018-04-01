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

// RegisterDevice registers a Device to a user. Creates the user if it doesn't
// exist yet. Fails if User already has a device with the same DeviceID.
func (t *Topic) RegisterDevice(d *Device) error {
	if err := d.validate(); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	user := t.getUser(d.UserID)

	if user.deviceExists(d.DeviceID) {
		return fmt.Errorf("device %q already exists for user", d.DeviceID)
	}

	user.doAddDevice(d)

	return nil
}

func (t *Topic) getUser(name string) *User {
	t.rw.RLock()
	user, ok := t.users[name]
	t.rw.RUnlock()

	if !ok {
		user = &User{
			Name:    name,
			devices: make([]*Device, 0),
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

func (t *Topic) doDeleteUser(user *User) {
	t.rw.Lock()
	delete(t.users, user.Name)
	t.rw.Unlock()
}
