package core

type Song struct {
    ID          int    `db:"id" json:"id"`
    Group   string `db:"group_name" json:"group"`
    Songname    string `db:"song_name" json:"song"`
    ReleaseDate string `db:"release_date" json:"release_date"`
    Text        string `db:"text" json:"text"`
	Link        string `db:"link" json:"link"`
}