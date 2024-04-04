package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {

    // Capture connection properties.
	cfg := "user= postgres password= postgres dbname= recordings host= localhost port=5432 sslmode=disable"

    // Get a database handle.
    var err error
    db, err = sql.Open("postgres", cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

    if err := createAlbumTable(); err != nil {
        log.Fatal(err)
    }

    albums, err := albumsByArtist("Betty Carter")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Albums found: %v\n", albums)

    // Hard-code ID 2 here to test the query.
    alb, err := albumByID(2)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Album found: %v\n", alb)

    albID, err := addAlbum(Album{
        Title:  "The Modern Sound of Betty Carter",
        Artist: "Betty Carter",
        Price:  49.99,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ID of added album: %v\n", albID)
}

func createAlbumTable() error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS album (
        ID SERIAL PRIMARY KEY,
        Title VARCHAR(255),
        Artist VARCHAR(255),
        Price NUMERIC(10, 2)
    );`)
    if err != nil {
        return err
    }
    return nil
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
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

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
    // An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
    var id int64
    err := db.QueryRow("SELECT id FROM album WHERE title = $1 AND artist = $2", alb.Title, alb.Artist).Scan(&id)
    switch {
    case err == sql.ErrNoRows:
        result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", alb.Title, alb.Artist, alb.Price)
        if err != nil {
            return 0, fmt.Errorf("addAlbum: failed to insert album: %v", err)
        }
        id, err := result.LastInsertId()
        if err != nil {
            return 0, fmt.Errorf("addAlbum: failed to get last insert ID: %v", err)
        }
        fmt.Printf("addAlbum: inserted album with ID %d\n", id)
        return id, nil
    case err != nil:
        return 0, fmt.Errorf("addAlbum: failed to check for existing album: %v", err)
    default:
        return 0, fmt.Errorf("addAlbum: album already exists with ID %d", id)
    }
}