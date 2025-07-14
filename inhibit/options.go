package inhibit

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type options struct {
	bus *dbus.Conn
}

type Option func(o *options)

func WithBus(bus *dbus.Conn) Option {
	return func(o *options) {
		o.bus = bus
	}
}

func readOptions(opts []Option) (options, error) {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.bus == nil {
		var err error
		options.bus, err = dbus.SessionBus()
		if err != nil {
			return options, fmt.Errorf("failed to retrieve default session bus: %w", err)
		}
	}
	return options, nil
}

func (o options) gnomeSessionObject() dbus.BusObject {
	return o.bus.Object("org.gnome.SessionManager", "/org/gnome/SessionManager")
}
