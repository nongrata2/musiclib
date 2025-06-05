package repositories

import (
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nongrata2/musiclib/internal/models"
	"github.com/nongrata2/musiclib/pkg/errors"
)

type DB struct {
	log  *slog.Logger
	conn *pgxpool.Pool
}

func addCondition(conditions *[]string, args *[]any, field string, value any, index *int) {
	switch v := value.(type) {
	case string:
		if v != "" {
			*conditions = append(*conditions, fmt.Sprintf("%s = $%d", field, *index))
			*args = append(*args, value)
			*index++
		}
	case time.Time:
		if !v.IsZero() {
			*conditions = append(*conditions, fmt.Sprintf("%s = $%d", field, *index))
			*args = append(*args, value)
			*index++
		}
	}
}

func New(log *slog.Logger, address string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Error("failed to ping database", "error", err)
		return nil, err
	}

	log.Info("successfully connected to database", "address", address)

	return &DB{
		log:  log,
		conn: pool,
	}, nil
}

func (db *DB) Add(ctx context.Context, song models.Song) error {

	db.log.Debug("started adding song DB")

	var groupID int
	query := `
        INSERT INTO groups (group_name)
        VALUES ($1)
        ON CONFLICT (group_name) DO NOTHING
        RETURNING id
    `
	err := db.conn.QueryRow(ctx, query, song.Group).Scan(&groupID)
	if err == pgx.ErrNoRows {
		query = `SELECT id FROM groups WHERE group_name = $1`
		err = db.conn.QueryRow(ctx, query, song.Group).Scan(&groupID)
		if err != nil {
			db.log.Error("failed to get group ID", "error", err)
			return err
		}
	} else if err != nil {
		db.log.Error("failed to add or check group", "error", err)
		return err
	}

	query = `
        INSERT INTO songs (group_id, song_name, release_date, text, link)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = db.conn.Exec(ctx, query,
		groupID,
		song.Songname,
		song.ReleaseDate,
		song.Text,
		song.Link,
	)

	if err != nil {
		db.log.Error("failed to add song", "error", err)
		return err
	}
	db.log.Debug("ended adding song DB")

	return nil
}

func (db *DB) GetSongs(ctx context.Context, filters models.SongFilter, page, limit int) ([]models.Song, error) {
	db.log.Debug("started getting song list DB")
	var songs []models.Song
	var maxLimit = 20

	query := `
        SELECT s.id, g.group_name, s.song_name, s.release_date, s.text, s.link
        FROM songs s
        JOIN groups g ON s.group_id = g.id
    `

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

	db.log.Debug("executing query", "query", query, "args", args)

	rows, err := db.conn.Query(ctx, query, args...)
	if err != nil {
		db.log.Error("failed to fetch songs", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID,
			&song.Group,
			&song.Songname,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		)
		if err != nil {
			db.log.Error("failed to scan song row", "error", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		db.log.Error("error while iterating over rows", "error", err)
		return nil, err
	}

	db.log.Debug("ended getting song list DB")
	return songs, nil
}

func (db *DB) Delete(ctx context.Context, songID string) error {
	db.log.Debug("started deleting song DB")

	query := `DELETE FROM songs WHERE id = $1`

	result, err := db.conn.Exec(ctx, query, songID)
	if err != nil {
		db.log.Error("failed to delete song", "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		db.log.Warn("no song found with the given id", "id", songID)
		return errors.NotFoundErr
	}
	db.log.Debug("ended deleting song DB")
	return nil
}

func (db *DB) GetLyrics(ctx context.Context, songID string, page, limit int) (string, error) {
	db.log.Debug("started getting lyrics DB")
	var songLyrics string

	query := `SELECT text FROM songs WHERE id = $1`
	err := db.conn.QueryRow(ctx, query, songID).Scan(&songLyrics)
	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			db.log.Error("no song found with the given ID", "id", songID)
			return "", errors.NotFoundErr
		}
		db.log.Error("failed to get lyrics of the song", "error", err)
		return "", err
	}

	if limit == 0 || page == 0 {
		db.log.Debug("pagination not used, returning full lyrics")
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

	db.log.Debug("ended getting lyrics DB")
	return result, nil
}

func (db *DB) Update(ctx context.Context, id int, song models.Song) (*models.Song, error) {
	db.log.Debug("started updating song DB")

	var groupID int
	query := `
        INSERT INTO groups (group_name)
        VALUES ($1)
        ON CONFLICT (group_name) DO NOTHING
        RETURNING id
    `
	err := db.conn.QueryRow(ctx, query, song.Group).Scan(&groupID)
	if err == pgx.ErrNoRows {
		query = `SELECT id FROM groups WHERE group_name = $1`
		err = db.conn.QueryRow(ctx, query, song.Group).Scan(&groupID)
		if err != nil {
			db.log.Error("failed to get group ID", "error", err)
			return nil, err
		}
	} else if err != nil {
		db.log.Error("failed to add or check group", "error", err)
		return nil, err
	}

	query = `
        UPDATE songs
        SET group_id = $1,
            song_name = $2,
            release_date = $3,
            text = $4,
            link = $5
        WHERE id = $6
        RETURNING id, (SELECT group_name FROM groups WHERE id = $1), song_name, release_date, text, link
    `

	var updatedSong models.Song
	err = db.conn.QueryRow(ctx, query,
		groupID,
		song.Songname,
		song.ReleaseDate,
		song.Text,
		song.Link,
		id,
	).Scan(
		&updatedSong.ID,
		&updatedSong.Group,
		&updatedSong.Songname,
		&updatedSong.ReleaseDate,
		&updatedSong.Text,
		&updatedSong.Link,
	)

	if err != nil {
		if stdErrors.Is(err, pgx.ErrNoRows) {
			db.log.Error("no song found with the given ID", "id", id)
			return nil, errors.NotFoundErr
		}
		db.log.Error("failed to update song", "error", err)
		return nil, err
	}

	db.log.Debug("end updating song DB")
	return &updatedSong, nil
}
