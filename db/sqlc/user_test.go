package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/chaojin101/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updateUserParams := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			Valid:  true,
			String: newFullName,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateUserParams)
	require.NoError(t, err)

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}
func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(8)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updateUserParams := UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			Valid:  true,
			String: newHashedPassword,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateUserParams)
	require.NoError(t, err)

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updateUserParams := UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			Valid:  true,
			String: newEmail,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateUserParams)
	require.NoError(t, err)

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(8)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updateUserParams := UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			Valid:  true,
			String: newFullName,
		},
		Email: sql.NullString{
			Valid:  true,
			String: newEmail,
		},
		HashedPassword: sql.NullString{
			Valid:  true,
			String: newHashedPassword,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateUserParams)
	require.NoError(t, err)

	require.Equal(t, oldUser.Username, updatedUser.Username)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
}
