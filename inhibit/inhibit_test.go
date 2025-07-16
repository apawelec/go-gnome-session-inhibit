package inhibit

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/assert"
)

func Test_Inhibit_Acquire(t *testing.T) {
	is := assert.New(t)

	// given
	appId := generateAppId()
	flags := Idle | Logout

	// when
	i, err := Acquire(appId, "Library testing", flags)
	t.Cleanup(func() { i.Release() })

	// then
	if is.NoError(err) {
		is.NotZero(i.cookie)
		is.NotNil(i)
	}

	if inhibitors, err := readActiveInhibitors(); is.NoError(err) {
		predicate := func(inhibitorPath string) (bool, error) {
			return inhibitorMatches(inhibitorPath, appId, flags)
		}
		if match, err := some(inhibitors, predicate); is.NoError(err) {
			is.True(match)
		}
	}
}

func Test_Inhibit_Release(t *testing.T) {
	is := assert.New(t)

	// given
	appId := generateAppId()
	flags := Idle | Logout

	i, err := Acquire(appId, "Library testing", flags)
	is.NoError(err)

	// when
	err = i.Release()

	// then
	if is.NoError(err) {
		is.Zero(i.cookie)
	}

	if inhibitors, err := readActiveInhibitors(); is.NoError(err) {
		predicate := func(inhibitorPath string) (bool, error) {
			return inhibitorMatches(inhibitorPath, appId, flags)
		}
		if match, err := some(inhibitors, predicate); is.NoError(err) {
			is.False(match)
		}
	}
}

func Test_No_Error_For_Multiple_Release(t *testing.T) {
	is := assert.New(t)

	// given
	appId := generateAppId()
	flags := Idle | Logout

	i, err := Acquire(appId, "Library testing", flags)
	is.NoError(err)

	err = i.Release()
	is.NoError(err)

	// when
	err = i.Release()

	// then
	is.NoError(err)
}

func generateAppId() string {
	return fmt.Sprintf("go-gnome-session-inhibit-test-%04d", rand.Int31n(10000))
}

func readActiveInhibitors() ([]string, error) {
	bus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	var inhibitors []string
	err = bus.Object("org.gnome.SessionManager", "/org/gnome/SessionManager").
		Call("org.gnome.SessionManager.GetInhibitors", 0).
		Store(&inhibitors)
	return inhibitors, err

}

func inhibitorMatches(inhibitorPath string, expectedAppId string, expectedFlags InhibitFlag) (bool, error) {
	bus, err := dbus.SessionBus()
	if err != nil {
		return false, nil
	}
	object := bus.Object("org.gnome.SessionManager", dbus.ObjectPath(inhibitorPath))

	var appId string
	if err := object.Call("org.gnome.SessionManager.Inhibitor.GetAppId", 0).Store(&appId); err != nil {
		return false, err
	}
	var flags uint32
	if err := object.Call("org.gnome.SessionManager.Inhibitor.GetFlags", 0).Store(&flags); err != nil {
		return false, err
	}
	return appId == expectedAppId && flags == uint32(expectedFlags), nil
}

func some[T any](seq []T, predicate func(T) (bool, error)) (bool, error) {
	for _, el := range seq {
		if match, err := predicate(el); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}
