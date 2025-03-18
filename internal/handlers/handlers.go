package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nongrata2/musiclib/internal/externalapi"
	"github.com/nongrata2/musiclib/internal/models"
	"github.com/nongrata2/musiclib/internal/repositories"
	"github.com/nongrata2/musiclib/pkg/errors"
)

func AddSongHandler(log *slog.Logger, db *repositories.DB, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("adding song handler")
		log.Info("start adding song")
		var request struct {
			Group    string `json:"group"`
			Songname string `json:"song"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Error("failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		apiResponse, err := externalapi.GetDataFromExternalAPI(apiBaseURL, request.Group, request.Songname)
		if err != nil {
			log.Error("failed to get data from external API", "error", err)
			http.Error(w, "Failed to get data from external API", http.StatusInternalServerError)
			return
		}

		newSong := models.Song{
			Group:       request.Group,
			Songname:    request.Songname,
			Text:        apiResponse.Text,
			ReleaseDate: apiResponse.ReleaseDate,
			Link:        apiResponse.Link,
		}

		if err := db.Add(r.Context(), newSong); err != nil {
			log.Error("failed to add song", "error", err)
			http.Error(w, "Failed to add song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("Song was added successfully\n"))
		if err != nil {
			log.Error("error writing", "error", err)
		}

		log.Info("end adding song")
	}
}

func GetLibDataHandler(log *slog.Logger, db *repositories.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("getting library data handler")
		log.Info("start getting data from library")

		filters := models.SongFilter{
			Group:       r.URL.Query().Get("group_name"),
			Songname:    r.URL.Query().Get("song_name"),
			ReleaseDate: r.URL.Query().Get("release_date"),
			Text:        r.URL.Query().Get("text"),
			Link:        r.URL.Query().Get("link"),
		}

		pagestr := r.URL.Query().Get("page")
		limitstr := r.URL.Query().Get("limit")

		var page, limit int
		var err error

		if pagestr != "" {
			page, err = strconv.Atoi(pagestr)
			if err != nil || page < 1 {
				log.Error("wrong page number", "error", err)
				http.Error(w, "wrong page number", http.StatusInternalServerError)
				return
			}
		} else {
			page = 0
		}

		if limitstr != "" {
			limit, err = strconv.Atoi(limitstr)
			if err != nil || limit < 1 {
				log.Error("wrong limit number", "error", err)
				http.Error(w, "wrong limit number", http.StatusInternalServerError)
				return
			}
		} else {
			limit = 0
		}

		if limit == 0 && page != 0 || page == 0 && limit != 0 {
			log.Warn("not using pagynation. both limit and page parameters should be filled")
			limit, page = 0, 0
		}

		songs, err := db.GetSongs(r.Context(), filters, page, limit)
		if err != nil {
			log.Error("failed to fetch songs", "error", err)
			http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
			return
		}

		if len(songs) == 0 {
			_, err := w.Write([]byte("No songs was found\n"))
			if err != nil {
				log.Error("error writing", "error", err)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")

			jsonData, err := json.MarshalIndent(songs, "", "  ")
			if err != nil {
				log.Error("failed to encode songs to JSON", "error", err)
				http.Error(w, "Failed to encode songs", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(jsonData)
			if err != nil {
				log.Error("error writing", "error", err)
			}
		}
		log.Info("end getting data from library")
	}
}

func DeleteSongHandler(log *slog.Logger, db *repositories.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("deleting song handler")
		log.Info("start deleting song")
		songID := r.PathValue("songID")

		if err := db.Delete(r.Context(), songID); err != nil {
			log.Error("failed to delete song", "error", err)
			http.Error(w, "Failed to delete song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		outstr := fmt.Sprintf("song with id %v was deleted successfully\n", songID)
		_, err := w.Write([]byte(outstr))
		if err != nil {
			log.Error("error writing", "error", err)
		}

		log.Info("end deleting song")
	}
}

func GetLyricsHandler(log *slog.Logger, db *repositories.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("getting lyrics handler")
		log.Info("start getting lyrics")

		songID := r.PathValue("songID")
		pagestr := r.URL.Query().Get("page")
		limitstr := r.URL.Query().Get("limit")

		var page, limit int
		var err error

		if pagestr != "" {
			page, err = strconv.Atoi(pagestr)
			if err != nil || page < 1 {
				log.Error("wrong page number", "error", err)
				http.Error(w, "wrong page number", http.StatusInternalServerError)
				return
			}
		} else {
			page = 0
		}

		if limitstr != "" {
			limit, err = strconv.Atoi(limitstr)
			if err != nil || limit < 1 {
				log.Error("wrong limit number", "error", err)
				http.Error(w, "wrong limit number", http.StatusInternalServerError)
				return
			}
		} else {
			limit = 0
		}

		if limit == 0 && page != 0 || page == 0 && limit != 0 {
			log.Warn("not using pagynation. both limit and page parameters should be filled")
			limit, page = 0, 0
		}

		songLyrics, err := db.GetLyrics(r.Context(), songID, page, limit)
		if err != nil {
			if err == errors.NotFoundErr {
				log.Error("no song with the given ID", "error", err)
				http.Error(w, "no song with the given ID", http.StatusNotFound)
				return
			}
			log.Error("failed to get lyrics of the song", "error", err)
			http.Error(w, "Failed to get lyrics of the song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(songLyrics))
		if err != nil {
			log.Error("error writing", "error", err)
		}
		log.Info("end getting lyrics")
	}
}

func EditSongHandler(log *slog.Logger, db *repositories.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("editing song handler")
		log.Info("start editing song")

		songIDstr := r.PathValue("songID")

		songID, err := strconv.Atoi(songIDstr)
		if err != nil {
			log.Error("invalid song ID", "error", err)
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}

		var song models.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			log.Error("failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updatedSong, err := db.Update(r.Context(), songID, song)
		if err != nil {
			if err == errors.NotFoundErr {
				log.Error("no song with the given ID", "error", err)
				http.Error(w, "no song with the given ID", http.StatusNotFound)
				return
			}
			log.Error("failed to update song", "error", err)
			http.Error(w, "Failed to update song", http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte("song was updated successfully\n"))
		if err != nil {
			log.Error("error writing", "error", err)
		}

		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.MarshalIndent(updatedSong, "", "  ")
		if err != nil {
			log.Error("failed to encode song to JSON", "error", err)
			http.Error(w, "Failed to encode song", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonData)
		if err != nil {
			log.Error("error writing", "error", err)
		}

		log.Info("end editing song")
	}
}
