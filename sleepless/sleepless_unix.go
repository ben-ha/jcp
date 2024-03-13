//go:build unix && !darwin

package sleepless

import (
	dbus "github.com/godbus/dbus/v5"
)

func PreventSleep(appName string, reason string) (func(), error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return func() {}, err
	}

	var cookie uint
	objectPath := dbus.ObjectPath("/org/gnome/SessionManager")

	conn.Object("org.gnome.SessionManager", objectPath).Call("org.gnome.SessionManager.Inhibit",
		0,
		"jcp",
		uint32(0),
		"Copy is in progress",
		uint32(1),
	).Store(&cookie)

	return func() {
		conn.Object("org.gnome.SessionManager", objectPath).Call(
			"org.gnome.SessionManager.Uninhibit",
			0,
			cookie)
		conn.Close()
	}, nil
}
