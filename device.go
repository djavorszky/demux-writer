package demux

import (
	"fmt"
	"io"
)

// addDevice adds a writing device to the User in question. Returns an error
// if name is empty or writer is nil
func (u *User) addDevice(name string, w io.Writer) error {
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
		UserID:   u.Name,
		DeviceID: name,
		writer:   w,
	}

	u.doAddDevice(d)

	return nil
}

func (d *Device) validate() error {
	if d.UserID == "" {
		return fmt.Errorf("userID required")
	}

	if d.DeviceID == "" {
		return fmt.Errorf("deviceID required")
	}

	if d.writer == nil {
		return fmt.Errorf("writer must be non-nil")
	}

	return nil
}

func (u *User) doAddDevice(device *Device) {
	u.rw.Lock()
	defer u.rw.Unlock()

	u.devices[device.DeviceID] = device
}

func (u *User) deviceExists(deviceID string) bool {
	u.rw.RLock()
	defer u.rw.RUnlock()

	_, ok := u.devices[deviceID]

	return ok
}

func (u *User) deleteDevice(deviceID string) {
	u.rw.RLock()
	_, ok := u.devices[deviceID]
	u.rw.RUnlock()

	if !ok {
		return
	}

	delete(u.devices, deviceID)
}

/*
	We only need RegisterDevice, but Device should have a user
	or some other identifier to it, so that when registering,
	we can register the device with an identifier.

	Meaning, that we can WriteTo(identifier) would write to
	all devices. But we also need a deviceID of some sort to
	be able to write to a specific device.

	Also deviceLabels, so we can write to a lot of devices
	irrespective of the "owners".
*/
