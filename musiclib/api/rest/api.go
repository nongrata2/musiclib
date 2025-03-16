package rest
	
import (
	"encoding/json"
	"log/slog"
	"net/http"
	"fmt"
	"strconv"

	"musiclib/db"
	"musiclib/api/core"
)


func AddSongHandler(log *slog.Logger, db *db.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var request struct {
            Group string `json:"group"`
            Songname  string `json:"song"`
			// Text string `json:"text"`
			// ReleaseDate string `json:"release_date"`
			// Link string `json:"link"`
        }

        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            log.Error("failed to decode request body", "error", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        newSong := core.Song{
            Group: request.Group,
            Songname:  request.Songname,
			Text: "i love you",
			ReleaseDate: "23.10.209",
			Link: "google.com",
			// Text:  request.Text,
			// ReleaseDate: request.ReleaseDate,
			// Link: request.Link,
        }

        if err := db.Add(r.Context(), newSong); err != nil {
            log.Error("failed to add song", "error", err)
            http.Error(w, "Failed to add song", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("Song added successfully\n"))
    }
}

func GetLibDataHandler(log *slog.Logger, db *db.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        filters := map[string]string{
            "group_name":   r.URL.Query().Get("group_name"),
            "song_name":    r.URL.Query().Get("song_name"),
            "release_date": r.URL.Query().Get("release_date"),
            "text":         r.URL.Query().Get("text"),
            "link":         r.URL.Query().Get("link"),
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
		}else{
			page = 0
		}
		
		if limitstr != "" {
			limit, err = strconv.Atoi(limitstr)
			if err != nil || limit < 1 {
				log.Error("wrong limit number", "error", err)
				http.Error(w, "wrong limit number", http.StatusInternalServerError)
				return
			}
		}else {
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

		if len(songs) == 0{
			w.Write([]byte("No songs was found\n")) 
		}else{
			w.Header().Set("Content-Type", "application/json")

			jsonData, err := json.MarshalIndent(songs, "", "  ")
			if err != nil {
				log.Error("failed to encode songs to JSON", "error", err)
				http.Error(w, "Failed to encode songs", http.StatusInternalServerError)
				return
			}
	
			w.Write(jsonData) 
		}
    }
}

func DeleteSongHandler(log *slog.Logger, db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID := r.PathValue("songID")

		if err := db.Delete(r.Context(), songID); err != nil {
            log.Error("failed to delete song", "error", err)
            http.Error(w, "Failed to delete song", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
		outstr := fmt.Sprintf("song with id %v deleted successfully\n", songID)
        w.Write([]byte(outstr))
	}
}

func GetLyricsHandler(log *slog.Logger, db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID := r.PathValue("songID")
		if songLyrics, err := db.GetLyrics(r.Context(), songID); err != nil {
            log.Error("failed to get lyrics of the song", "error", err)
            http.Error(w, "Failed to get lyrics of the song", http.StatusInternalServerError)
            return
        }else{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(songLyrics))
		}
	}
}

func EditSongHandler(log *slog.Logger, db *db.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        songIDstr := r.PathValue("songID")

        songID, err := strconv.Atoi(songIDstr)
        if err != nil {
            log.Error("invalid song ID", "error", err)
            http.Error(w, "Invalid song ID", http.StatusBadRequest)
            return
        }

        var song core.Song
        if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
            log.Error("failed to decode request body", "error", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        updatedSong, err := db.Update(r.Context(), songID, song)
        if err != nil {
            log.Error("failed to update song", "error", err)
            http.Error(w, "Failed to update song", http.StatusInternalServerError)
            return
        }

		w.Write([]byte("song was updated successfully\n"))
		w.WriteHeader(http.StatusOK)

        w.Header().Set("Content-Type", "application/json")

        jsonData, err := json.MarshalIndent(updatedSong, "", "  ")
        if err != nil {
            log.Error("failed to encode song to JSON", "error", err)
            http.Error(w, "Failed to encode song", http.StatusInternalServerError)
            return
        }
        w.Write(jsonData)
    }
}