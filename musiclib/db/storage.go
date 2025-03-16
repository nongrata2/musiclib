package db

import (
	"context"
	"errors"
	"log/slog"
	"fmt"
	"strings"
    "os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"musiclib/api/core"
	"database/sql"
    "github.com/DATA-DOG/go-sqlmock"
)

type DB struct {
	Log  *slog.Logger
	Conn *sqlx.DB
}

// mock for tests
func NewMockDB() (*DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	return &DB{
		Log:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Conn: db,
	}, mock, nil
}

func New(log *slog.Logger, address string) (*DB, error) {

	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		Log:  log,
		Conn: db,
	}, nil
}

func (db *DB) Add(ctx context.Context, song core.Song) error {

    query := `
        INSERT INTO songs (group_name, song_name, release_date, text, link)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := db.Conn.ExecContext(ctx, query,
        song.Group,
        song.Songname,
        song.ReleaseDate,
        song.Text,
        song.Link,
    )
    if err != nil {
        db.Log.Error("failed to add song", "error", err)
        return err
    }

    return nil
}

func (db *DB) GetSongs(ctx context.Context, filters map[string]string, page, limit int) ([]core.Song, error) {
    var songs []core.Song

    query := `SELECT id, group_name, song_name, release_date, text, link FROM songs`

    var conditions []string
    var args []interface{}
    i := 1

    for key, value := range filters {
        if value != "" {
            conditions = append(conditions, fmt.Sprintf("%s = $%d", key, i))
            args = append(args, value)
            i++
        }
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

	if limit != 0 && page != 0 { // need pagy
        offset := (page - 1) * limit
        query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

    err := db.Conn.SelectContext(ctx, &songs, query, args...)
    if err != nil {
        db.Log.Error("failed to fetch songs", "error", err)
        return nil, err
    }

    return songs, nil
}

func (db *DB) Delete(ctx context.Context, songID string) error {
    query := `DELETE FROM songs WHERE id = $1`

    result, err := db.Conn.ExecContext(ctx, query, songID)
    if err != nil {
        db.Log.Error("failed to delete song", "error", err)
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        db.Log.Error("failed to get rows affected", "error", err)
        return err
    }

    if rowsAffected == 0 {
        db.Log.Warn("no song found with the given id", "id", songID)

        return errors.New("no song found with the given ID")
    }

    return nil
}

func (db *DB) GetLyrics(ctx context.Context, songID string, page, limit int) (string, error) {
	var songLyrics string

	query := `SELECT text FROM songs WHERE id = $1`
	err := db.Conn.GetContext(ctx, &songLyrics, query, songID)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Log.Error("no song found with the given ID", "id", songID)
			return "", errors.New("no song found with the given ID")
		}
		db.Log.Error("failed to get lyrics of the song", "error", err)
		return "", err
	}

	var result string

	if limit != 0 && page != 0{
		verses := strings.Split(songLyrics, "\n\n")

		start := (page - 1) * limit
		end := start + limit
	
		if start >= len(verses) {
			return "", errors.New("page out of range")
		}
		if end > len(verses) {
			end = len(verses)
		}
	
		paginatedVerses := verses[start:end]
	
		result = strings.Join(paginatedVerses, "\n\n")
	}

	return result, nil
}

func (db *DB) Update(ctx context.Context, id int, song core.Song) (*core.Song, error) {

    var exists bool
    checkQuery := `SELECT EXISTS(SELECT 1 FROM songs WHERE id = $1)`
    err := db.Conn.GetContext(ctx, &exists, checkQuery, id)
    if err != nil {
        db.Log.Error("failed to check if song exists", "error", err)
        return nil, err
    }

    if !exists {
        db.Log.Error("no song found with the given ID", "id", id)
        return nil, errors.New("no song found with the given ID")
    }

    query := `
        UPDATE songs
        SET group_name = :group_name,
            song_name = :song_name,
            release_date = :release_date,
            text = :text,
            link = :link
        WHERE id = :id
    `

    song.ID = id

    _, err = db.Conn.NamedExecContext(ctx, query, song)
    if err != nil {
        db.Log.Error("failed to update song", "error", err)
        return nil, err
    }

    var updatedSong core.Song
    getQuery := `SELECT * FROM songs WHERE id = $1`
    err = db.Conn.GetContext(ctx, &updatedSong, getQuery, id)
    if err != nil {
        db.Log.Error("failed to fetch updated song", "error", err)
        return nil, err
    }

    return &updatedSong, nil
}