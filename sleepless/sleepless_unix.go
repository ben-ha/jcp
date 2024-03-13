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
	ifaceName := "org.gnome.SessionManager"

	conn.Object("org.gnome.SessionManager", objectPath).Call(
		ifaceName+".Inhibit",
		0,                     // Flags
		"jcp",                 // App Name
		uint32(0),             // Flags
		"Copy is in progress", // Reason
		uint32(1),             // Flags
	).Store(&cookie)

	return func() {
		conn.Object("org.gnome.SessionManager", objectPath).Call(
			ifaceName+".Uninhibit",
			0, // Flags
			cookie)
		conn.Close()
	}, nil
}
