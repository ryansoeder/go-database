package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID int64
	Title string
	Artist string
	Price float32
}

func main () {
	// Capture connection properties.
	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "recordings",
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	albums, err := albumsByArtist("Betty Carter")
	if err != nil {
		log.Fatal("err")
	}
	fmt.Printf("ablums found: %v\n", albums)

	album, err := albumByID(5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("album found: %v\n", album)

	addId, err := addAlbum(Album{
		Title:  "Cool Album",
		Artist: "Dyontai Swan",
		Price:  499.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("album added: %v\n", addId)

// 	deleteID, err := deleteAlbum(6)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("album deleted: %v\n", deleteID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold dta from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumBy Id queries for a single album by its ID
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %q: No such album", id)
		}
		return alb, fmt.Errorf("albumByID %q: %v\n", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database, returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (Title, Artist, Price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// deleteAlbum deletes the specified album by id
func deleteAlbum(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	
	albId, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return albId, nil
}
