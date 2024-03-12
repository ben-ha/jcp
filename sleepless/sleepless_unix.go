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
	err = conn.BusObject().Call("org.gnome.SessionManager.Inhibit", 0, appName, 0, reason, 0).Store(&cookie)
	if err != nil {
		return func() {}, err
	}

	return func() {
		conn.BusObject().Call("org.gnome.SessionManager.Uninhibit", 0, cookie)
		conn.Close()
	}, nil
}
