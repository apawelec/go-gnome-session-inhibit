package inhibit

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Inhibit_Acquire(t *testing.T) {
	is := assert.New(t)

	// when
	i, err := Acquire(generateAppId(), "Library testing", Idle|Logout)

	// then
	if is.NoError(err) {
		is.NotZero(i.cookie)
		is.NotNil(i)
	}
}

func Test_Inhibit_Release(t *testing.T) {
	is := assert.New(t)

	// given
	i, err := Acquire(generateAppId(), "Library testing", Idle|Logout)
	is.NoError(err)

	// when
	err = i.Release()

	// then
	if is.NoError(err) {
		is.Zero(i.cookie)
	}
}

func Test_No_Error_For_Multiple_Release(t *testing.T) {
	is := assert.New(t)

	// given
	i, err := Acquire(generateAppId(), "Library testing", Idle|Logout)
	is.NoError(err)

	err = i.Release()
	is.NoError(err)

	// when
	err = i.Release()

	// then
	is.NoError(err)
}

func generateAppId() string {
	return fmt.Sprintf("test-app-%03d", rand.Int31n(1000))
}
