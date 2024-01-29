package db

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/Streamfair/streamfair-idp-svc/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) IdpSvcUser {
	byteHash, err := util.HashPassword(util.RandomPassword())
	hashedPassword := base64.StdEncoding.EncodeToString(byteHash.Hash)
	passwordSalt := base64.StdEncoding.EncodeToString(byteHash.Salt)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:     util.RandomUsername(),
		FullName:     util.RandomString(12),
		Email:        util.RandomEmail(),
		PasswordHash: hashedPassword,
		PasswordSalt: passwordSalt,
		CountryCode:  util.RandomCountryCode(),
		RoleID:       util.ConvertToInt8(util.RandomInt(1, 3)),
		Status:       util.ConvertToText(util.RandomString(12)),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotEmpty(t, user.PasswordHash)
	require.NotEmpty(t, user.PasswordSalt)
	require.Equal(t, arg.CountryCode, user.CountryCode)
	require.Equal(t, arg.RoleID, user.RoleID)
	require.Equal(t, arg.Status, user.Status)
	require.True(t, user.LastLoginAt.Time.IsZero())
	require.True(t, user.UsernameChangedAt.Time.IsZero())
	require.True(t, user.EmailChangedAt.Time.IsZero())
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.WithinDuration(t, time.Now(), user.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), user.UpdatedAt.Time, time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user := createRandomUser(t)

	fetchedUser, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)
	require.Equal(t, user.ID, fetchedUser.ID)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.FullName, fetchedUser.FullName)
	require.Equal(t, user.Email, fetchedUser.Email)
	require.NotEmpty(t, fetchedUser.PasswordHash)
	require.NotEmpty(t, fetchedUser.PasswordSalt)
	require.Equal(t, user.CountryCode, fetchedUser.CountryCode)
	require.Equal(t, user.RoleID, fetchedUser.RoleID)
	require.Equal(t, user.Status, fetchedUser.Status)
	require.True(t, fetchedUser.LastLoginAt.Time.IsZero())
	require.True(t, fetchedUser.UsernameChangedAt.Time.IsZero())
	require.True(t, fetchedUser.EmailChangedAt.Time.IsZero())
	require.True(t, fetchedUser.PasswordChangedAt.Time.IsZero())
	require.WithinDuration(t, user.CreatedAt.Time, fetchedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), fetchedUser.UpdatedAt.Time, time.Second)
}

func TestGetUserByUsername(t *testing.T) {
	user := createRandomUser(t)

	fetchedUser, err := testQueries.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)
	require.Equal(t, user.ID, fetchedUser.ID)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.FullName, fetchedUser.FullName)
	require.Equal(t, user.Email, fetchedUser.Email)
	require.NotEmpty(t, fetchedUser.PasswordHash)
	require.NotEmpty(t, fetchedUser.PasswordSalt)
	require.Equal(t, user.CountryCode, fetchedUser.CountryCode)
	require.Equal(t, user.RoleID, fetchedUser.RoleID)
	require.Equal(t, user.Status, fetchedUser.Status)
	require.True(t, fetchedUser.LastLoginAt.Time.IsZero())
	require.True(t, fetchedUser.UsernameChangedAt.Time.IsZero())
	require.True(t, fetchedUser.EmailChangedAt.Time.IsZero())
	require.True(t, fetchedUser.PasswordChangedAt.Time.IsZero())
	require.WithinDuration(t, user.CreatedAt.Time, fetchedUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), fetchedUser.UpdatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	fetchedUser, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, fetchedUser)
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
		require.NotZero(t, user.ID)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.FullName)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.CountryCode)
		require.NotZero(t, user.RoleID)
		require.NotEmpty(t, user.Status)
		require.NotZero(t, user.LastLoginAt.Time)
		require.NotZero(t, user.CreatedAt.Time)
		require.NotZero(t, user.UpdatedAt.Time)
	}
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserParams{
		ID:          user.ID,
		Username:    util.RandomUsername(),
		FullName:    util.RandomString(12),
		Status:      util.ConvertToText(util.RandomString(12)),
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Username, updatedUser.Username)
	require.NotEqual(t, user.FullName, updatedUser.FullName)
	require.NotEqual(t, user.Status, updatedUser.Status)
	require.True(t, user.LastLoginAt.Time.IsZero())
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Minute)
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	byteHash, err := util.HashPassword(util.RandomPassword())
	hashedPassword := base64.StdEncoding.EncodeToString(byteHash.Hash)
	passwordSalt := base64.StdEncoding.EncodeToString(byteHash.Salt)
	require.NoError(t, err)

	arg := UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: hashedPassword,
		PasswordSalt: passwordSalt,
	}

	updatedUser, err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEmpty(t, updatedUser.PasswordHash)
	require.NotEmpty(t, updatedUser.PasswordSalt)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Minute)
	require.WithinDuration(t, time.Now(), updatedUser.PasswordChangedAt.Time, time.Minute)
}

func TestUpdateUserEmail(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserEmailParams{
		ID:    user.ID,
		Email: util.RandomEmail(),
	}

	updatedUser, err := testQueries.UpdateUserEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Minute)
	require.WithinDuration(t, time.Now(), updatedUser.EmailChangedAt.Time, time.Minute)
}

func TestUpdateUsername(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUsernameParams{
		ID:       user.ID,
		Username: util.RandomString(12),
	}

	updatedUser, err := testQueries.UpdateUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, user.Username, updatedUser.Username)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Minute)
	require.WithinDuration(t, time.Now(), updatedUser.UsernameChangedAt.Time, time.Minute)
}
