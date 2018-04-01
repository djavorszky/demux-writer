package demux

import (
	"fmt"
)

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

// UnregisterDevice removes the device from the User.
func (t *Topic) UnregisterDevice(userID, deviceID string) {
	t.getUser(userID).deleteDevice(deviceID)
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
