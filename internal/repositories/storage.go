package repositories

import (
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"
	"strings"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/nongrata2/musiclib/internal/models"
	"github.com/nongrata2/musiclib/pkg/errors"
)

type DB struct {
	Log  *slog.Logger
	Conn *sqlx.DB
}

func addCondition(conditions *[]string, args *[]any, field string, value string, index *int) {
	if value != "" {
		*conditions = append(*conditions, fmt.Sprintf("%s = $%d", field, *index))
		*args = append(*args, value)
		*index++
	}
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

func (db *DB) Add(ctx context.Context, song models.Song) error {

	db.Log.Debug("started adding song DB")
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
	db.Log.Debug("ended adding song DB")

	return nil
}

func (db *DB) GetSongs(ctx context.Context, filters models.SongFilter, page, limit int) ([]models.Song, error) {
	db.Log.Debug("started getting song list DB")
	var songs []models.Song
	var maxLimit = 20

	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs`

	var conditions []string
	var args []any
	i := 1

	addCondition(&conditions, &args, "group_name", filters.Group, &i)
	addCondition(&conditions, &args, "song_name", filters.Songname, &i)
	addCondition(&conditions, &args, "release_date", filters.ReleaseDate, &i)
	addCondition(&conditions, &args, "text", filters.Text, &i)
	addCondition(&conditions, &args, "link", filters.Link, &i)

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limit != 0 && page != 0 {
		if limit > maxLimit {
			limit = maxLimit
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	err := db.Conn.SelectContext(ctx, &songs, query, args...)
	if err != nil {
		db.Log.Error("failed to fetch songs", "error", err)
		return nil, err
	}

	db.Log.Debug("ended getting song list DB")
	return songs, nil
}

func (db *DB) Delete(ctx context.Context, songID string) error {
	db.Log.Debug("started deleting song DB")

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

		return errors.NotFoundErr
	}
	db.Log.Debug("ended deleting song DB")
	return nil
}

func (db *DB) GetLyrics(ctx context.Context, songID string, page, limit int) (string, error) {
	db.Log.Debug("started getting lyrics DB")
	var songLyrics string

	query := `SELECT text FROM songs WHERE id = $1`
	err := db.Conn.GetContext(ctx, &songLyrics, query, songID)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			db.Log.Error("no song found with the given ID", "id", songID)
			return "", errors.NotFoundErr
		}
		db.Log.Error("failed to get lyrics of the song", "error", err)
		return "", err
	}

	if limit == 0 || page == 0 {
		db.Log.Debug("pagination not used, returning full lyrics")
		return songLyrics, nil
	}

	verses := strings.Split(songLyrics, "\n\n")

	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return "", errors.OutOfRangeErr
	}
	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := verses[start:end]
	result := strings.Join(paginatedVerses, "\n\n")

	db.Log.Debug("ended getting lyrics DB")
	return result, nil
}

func (db *DB) Update(ctx context.Context, id int, song models.Song) (*models.Song, error) {
	db.Log.Debug("started updating song DB")

	query := `
        UPDATE songs
        SET group_name = $1,
            song_name = $2,
            release_date = $3,
            text = $4,
            link = $5
        WHERE id = $6
        RETURNING *
    `

	var updatedSong models.Song
	err := db.Conn.QueryRowContext(ctx, query, song.Group, song.Songname, song.ReleaseDate, song.Text, song.Link, id).Scan(
		&updatedSong.ID,
		&updatedSong.Group,
		&updatedSong.Songname,
		&updatedSong.ReleaseDate,
		&updatedSong.Text,
		&updatedSong.Link,
	)

	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			db.Log.Error("no song found with the given ID", "id", id)
			return nil, errors.NotFoundErr
		}
		db.Log.Error("failed to update song", "error", err)
		return nil, err
	}

	db.Log.Debug("end updating song DB")
	return &updatedSong, nil
}
