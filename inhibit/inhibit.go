package inhibit

import (
	"fmt"
)

type InhibitFlag uint

const (
	Logout     InhibitFlag = 1
	SwitchUser InhibitFlag = 2
	Suspend    InhibitFlag = 4
	Idle       InhibitFlag = 8
	Automount  InhibitFlag = 16
)

type AcquiredInhibit struct {
	cookie uint
}

func Acquire(appId string, reason string, flags InhibitFlag, opts ...Option) (*AcquiredInhibit, error) {
	options, err := readOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("error while reading options: %w", err)
	}

	gnomeSession := options.gnomeSessionObject()

	var cookie uint
	if err = gnomeSession.Call("org.gnome.SessionManager.Inhibit", 0, appId, uint(0), reason, flags).Store(&cookie); err != nil {
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

	if result := gnomeSession.Call("org.gnome.SessionManager.Uninhibit", 0, i.cookie); result.Err != nil {
		return fmt.Errorf("error while calling Uninhibit method: %w", result.Err)
	}
	i.cookie = 0
	return nil
}
