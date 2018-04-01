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

// WriteToUser writes the message to all the devices a specific user has.
// Returns an error if the user does not exist.
func (t *Topic) WriteToUser(userID string, message []byte) error {
	if userID == "" {
		return fmt.Errorf("userID is empty")
	}

	if !t.userExists(userID) {
		return fmt.Errorf("user %q does not exist", userID)
	}

	u := t.getUser(userID)
	for _, d := range u.devices {
		fmt.Println(d)
		d.writer.Write(message)
	}

	return nil
}

// WriteToDevice writes the message to a specific user's specific device. Returns
// an error if the user does not exist or does not have a device specified by
// the deviceID
func (t *Topic) WriteToDevice(userID, deviceID string, message []byte) error {
	d, err := t.getDevice(userID, deviceID)
	if err != nil {
		return fmt.Errorf("couldn't get device: %v", err)
	}

	d.writer.Write(message)

	return nil
}

func (t *Topic) getDevice(userID, deviceID string) (*Device, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	if deviceID == "" {
		return nil, fmt.Errorf("deviceID is empty")
	}

	if !t.userExists(userID) {
		return nil, fmt.Errorf("user %q does not exist", userID)
	}

	u := t.getUser(userID)
	if !u.deviceExists(deviceID) {
		return nil, fmt.Errorf("device %q does not exist for user %q", deviceID, userID)
	}

	device := u.doGetDevice(deviceID)

	return device, nil
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

func (u *User) doGetDevice(deviceID string) *Device {
	u.rw.RLock()
	defer u.rw.RUnlock()

	return u.devices[deviceID]
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
