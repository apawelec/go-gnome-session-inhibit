package inhibit

import (
	"fmt"
)

type InhibitFlag uint32

const (
	Logout InhibitFlag = 1 << iota
	SwitchUser
	Suspend
	Idle
	Automount
)

type AcquiredInhibit struct {
	cookie uint32
}

func Acquire(appId string, reason string, flags InhibitFlag, opts ...Option) (*AcquiredInhibit, error) {
	options, err := readOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("error while reading options: %w", err)
	}

	gnomeSession := options.gnomeSessionObject()

	var cookie uint32
	if err = gnomeSession.CallWithContext(
		options.ctx,
		"org.gnome.SessionManager.Inhibit",
		0,
		appId,
		uint32(0),
		reason,
		flags,
	).Store(&cookie); err != nil {
		return nil, fmt.Errorf("error while calling Inhibit method: %w", err)
	}

	return &AcquiredInhibit{cookie: cookie}, nil
}

func (i *AcquiredInhibit) Release(opts ...Option) error {
	if i.cookie == 0 {
		return nil
	}

	options, err := readOptions(opts)
	if err != nil {
		return fmt.Errorf("error while reading options: %w", err)
	}

	gnomeSession := options.gnomeSessionObject()

	if result := gnomeSession.CallWithContext(
		options.ctx,
		"org.gnome.SessionManager.Uninhibit",
		0,
		i.cookie,
	); result.Err != nil {
		return fmt.Errorf("error while calling Uninhibit method: %w", result.Err)
	}
	i.cookie = 0
	return nil
}
