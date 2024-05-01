package tests

import (
	pb "aut_reg/proto/gen"
	"aut_reg/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	secret         = "test-secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	username := gofakeit.Username()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, username, claims["username"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	username := gofakeit.Username()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "Pleas enter Password",
		},
		{
			name:        "Register with Empty Username",
			username:    "",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			expectedErr: "Pleas enter Username",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			username:    gofakeit.Username(),
			password:    randomFakePassword(),
			expectedErr: "Pleas enter Email",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			username:    "",
			password:    "",
			expectedErr: "Pleas enter Email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
				Email:    tt.email,
				Username: tt.username,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "Pleas enter Password",
		},
		{
			name:        "Login with Empty Username",
			username:    "",
			password:    randomFakePassword(),
			expectedErr: "Pleas enter Username",
		},
		{
			name:        "Login with Both Empty Username and Password",
			username:    "",
			password:    "",
			expectedErr: "Pleas enter Username",
		},
		{
			name:        "Login with Non-Matching Password",
			username:    gofakeit.Username(),
			password:    randomFakePassword(),
			expectedErr: "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
				Email:    gofakeit.Email(),
				Username: gofakeit.Username(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &pb.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
