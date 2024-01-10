package data_access

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// https://go.dev/doc/tutorial/database-access

var db *sql.DB

// Making db a global variable simplifies this example.
// In production, you’d avoid the global variable,
// such as by passing the variable to functions that
// need it or by wrapping it in a struct.

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// AlbumsByArtist queries for albums that have the specified artist name.
func AlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("AlbumsbyArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// AlbumByID queries for the album with the specified ID.
func AlbumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album where id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("AlbumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("AlbumByID %d: %v", id, err)
	}
	return alb, nil
}

// AddAlbum adds the specified album to the database,
// returning the album ID of the new entry
func AddAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	return id, nil
}

func EstablishConnection() {
	// Capture connection properties
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "MusicPlayer",
	}
	// Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	// To simplify the code, you’re calling log.Fatal to end
	// execution and print the error to the console. In
	// production code, you’ll want to handle errors in a
	// more graceful way.

	fmt.Println("Connected!")
}
