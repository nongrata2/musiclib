package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"musiclib/api/core"
	"musiclib/db"
	"log/slog"
)

func TestCreateListDeleteList(t *testing.T) {
	mockDB, mock, err := db.NewMockDB()
	assert.NoError(t, err)
	defer mockDB.Conn.Close()

	apiBaseURL := "http://localhost:8082" 
	handler := AddSongHandler(slog.Default(), mockDB, apiBaseURL)

	// 1. adding song
	song := core.Song{
		Group:       "Muse",
		Songname:    "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	mock.ExpectExec("INSERT INTO songs").
		WithArgs(song.Group, song.Songname, song.ReleaseDate, song.Text, song.Link).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(song)
	req := httptest.NewRequest(http.MethodPost, "/songs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. song list
	rows := sqlmock.NewRows([]string{"id", "group_name", "song_name", "release_date", "text", "link"}).
		AddRow(1, song.Group, song.Songname, song.ReleaseDate, song.Text, song.Link)

	mock.ExpectQuery("SELECT id, group_name, song_name, release_date, text, link FROM songs").
		WillReturnRows(rows)

	listHandler := GetLibDataHandler(slog.Default(), mockDB)
	req = httptest.NewRequest(http.MethodGet, "/songs", nil)
	w = httptest.NewRecorder()
	listHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var songs []core.Song
	err = json.Unmarshal(w.Body.Bytes(), &songs)
	assert.NoError(t, err)
	assert.Len(t, songs, 1)
	assert.Equal(t, song.Group, songs[0].Group)

	// 3. delete song
	mock.ExpectExec("DELETE FROM songs WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	deleteHandler := DeleteSongHandler(slog.Default(), mockDB)
	req = httptest.NewRequest(http.MethodDelete, "/songs/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "songID", "1"))
	w = httptest.NewRecorder()
	deleteHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// // 4. Вывод списка после удаления
	// mock.ExpectQuery("SELECT id, group_name, song_name, release_date, text, link FROM songs").
	// 	WillReturnRows(sqlmock.NewRows([]string{"id", "group_name", "song_name", "release_date", "text", "link"}))

	// req = httptest.NewRequest(http.MethodGet, "/songs", nil)
	// w = httptest.NewRecorder()
	// listHandler(w, req)

	// assert.Equal(t, http.StatusOK, w.Code)
	// err = json.Unmarshal(w.Body.Bytes(), &songs)
	// assert.NoError(t, err)
	// assert.Len(t, songs, 0)

	// // Проверяем, что все ожидаемые запросы были выполнены
	// assert.NoError(t, mock.ExpectationsWereMet())
}
