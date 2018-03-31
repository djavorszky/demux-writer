package demux

import (
	"fmt"
	"io"
)

// AddDevice adds a writing device to the User in question. Returns an error
// if name is empty or writer is nil
func (u *User) AddDevice(name string, w io.Writer) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	if w == nil {
		return fmt.Errorf("writer is nil")
	}

	if u.deviceExists(name) {
		return fmt.Errorf("device with name %q already exists", name)
	}

	d := &Device{
		Name:   name,
		writer: w,
	}

	u.doAddDevice(d)

	return nil
}

func (u *User) doAddDevice(device *Device) {
	u.rw.Lock()
	defer u.rw.Unlock()

	u.devices = append(u.devices, device)
}

func (u *User) deviceExists(name string) bool {
	u.rw.RLock()
	defer u.rw.RUnlock()

	for _, d := range u.devices {
		if d.Name == name {
			return true
		}
	}

	return false
}
