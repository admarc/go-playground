package users

import (
	"context"
	"database/sql"
	"testing"

	"github.com/admarc/users/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	db, err := sql.Open("sqlite3", "/home/adam/Projects/go-playground/users/tmp/db.sqlite3")

	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		s := Storage{
			db: db,
		}

		usr, err := s.Create(ctx, "mike")
		require.NoError(t, err)

		dbUser := models.User{}
		row := db.QueryRowContext(ctx, "SELECT id, name from users where id = ?", usr.ID)
		err = row.Scan(&dbUser.ID, &dbUser.Name)
		require.NoError(t, err)

		assert.Equal(t, usr, dbUser)
	})

}
