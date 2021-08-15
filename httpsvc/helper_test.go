package httpsvc

import (
	"testing"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/stretchr/testify/require"
)

func TestJWT_Valid(t *testing.T) {
	jwtKey = []byte("secret")
	user := &model.User{
		Base: model.Base{ID: "11"},
	}
	expiredAt := time.Now().Add(1 * time.Hour)
	time.Sleep(3 * time.Microsecond)
	token, err := generateAccessToken(user, expiredAt)
	require.NoError(t, err)
	_, err = parseJWTToken(token)
	require.NoError(t, err)
}

func TestJWT_Expired(t *testing.T) {
	jwtKey = []byte("secret")
	user := &model.User{
		Base: model.Base{ID: "11"},
	}
	expiredAt := time.Now().Add(-1 * time.Second)
	time.Sleep(3 * time.Microsecond)
	token, err := generateAccessToken(user, expiredAt)
	require.NoError(t, err)
	_, err = parseJWTToken(token)
	require.Error(t, err)
}

func TestJWT_Invalid(t *testing.T) {
	jwtKey = []byte("secret")
	time.Sleep(3 * time.Microsecond)
	_, err := parseJWTToken("randoms")
	require.Error(t, err)
}
