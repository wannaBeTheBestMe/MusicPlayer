package data_access

import (
	"database/sql"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/go-sql-driver/mysql"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// https://go.dev/doc/tutorial/database-access

var DB *sql.DB

func EstablishConnection() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "MusicPlayer",
	}
	// Get a database handle.
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

type Track struct {
	ID                 int64
	Format             string
	FileType           string
	Title              string
	Album              string
	Artist             string
	AlbumArtist        string
	Composer           string
	Year               int
	Genre              string
	TrackNum           int
	TrackTotal         int
	DiscNum            int
	DiscTotal          int
	PictureExt         string
	PictureMIMEType    string
	PictureType        string
	PictureDescription string
	PictureData        []byte
	Lyrics             string
	Comment            string
	Filepath           string
	Freq               int
	Valid              bool
}

var MusicDir = "C:\\Users\\Asus\\Music\\MusicPlayer"

func GetAlbumDirectories() {
	albumsDir, err := os.ReadDir(MusicDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, albumDir := range albumsDir {
		fmt.Println(albumDir)
	}
}

func GetFulltextSearchResults(term string) ([]Track, error) {
	var tracks []Track

	queryString := `
		SELECT * FROM tracks WHERE MATCH(Title, Album, Artist, AlbumArtist, Composer, Lyrics) AGAINST(? IN NATURAL LANGUAGE MODE);
	`

	rows, err := DB.Query(queryString, term)
	if err != nil {
		return []Track{}, fmt.Errorf("GetFulltextSearchResults: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		err := rows.Scan(&track.ID, &track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist,
			&track.AlbumArtist, &track.Composer, &track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal,
			&track.DiscNum, &track.DiscTotal, &track.PictureExt, &track.PictureMIMEType, &track.PictureType,
			&track.PictureDescription, &track.PictureData, &track.Lyrics, &track.Comment, &track.Filepath, &track.Freq)
		if err != nil {
			return []Track{}, fmt.Errorf("GetFulltextSearchResults: %v", err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

type DataPoint struct {
	Field string
	Freq  int
}

func GetTopByGenreFreq() (Track, error) {
	var genre string
	var freq float64

	queryString := `
		SELECT Genre, SUM(Freq) AS total
		FROM tracks
		WHERE Genre IS NOT NULL AND Genre != '' AND Genre != ' '
		GROUP BY Genre
		ORDER BY total DESC
		LIMIT 1;
	`

	row := DB.QueryRow(queryString)
	err := row.Scan(&genre, &freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByGenreFreq: %v", err)
	}

	var topID string

	queryString = `
		SELECT ID
		FROM tracks
		WHERE Genre = ?
		GROUP BY ID
		ORDER BY SUM(Freq) DESC
		LIMIT 1;
	`

	row = DB.QueryRow(queryString, genre)
	err = row.Scan(&topID)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByGenreFreq: %v", err)
	}

	var track Track

	queryString = `
		SELECT *
		FROM tracks
		WHERE ID = ?
	`

	row = DB.QueryRow(queryString, topID)
	err = row.Scan(&track.ID, &track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist,
		&track.AlbumArtist, &track.Composer, &track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal,
		&track.DiscNum, &track.DiscTotal, &track.PictureExt, &track.PictureMIMEType, &track.PictureType,
		&track.PictureDescription, &track.PictureData, &track.Lyrics, &track.Comment, &track.Filepath, &track.Freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByGenreFreq: %v", err)
	}

	return track, nil
}

func GetTopByArtistFreq() (Track, error) {
	var artist string
	var freq float64

	queryString := `
		SELECT Artist, SUM(Freq) AS total
		FROM tracks
		WHERE Artist IS NOT NULL AND Artist != '' AND Artist != ' '
		GROUP BY Artist
		ORDER BY total DESC
		LIMIT 1;
	`

	row := DB.QueryRow(queryString)
	err := row.Scan(&artist, &freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByArtistFreq: %v", err)
	}

	var topID string

	queryString = `
		SELECT ID
		FROM tracks
		WHERE Artist = ?
		GROUP BY ID
		ORDER BY SUM(Freq) DESC
		LIMIT 1;
	`

	row = DB.QueryRow(queryString, artist)
	err = row.Scan(&topID)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByArtistFreq: %v", err)
	}

	var track Track

	queryString = `
		SELECT *
		FROM tracks
		WHERE ID = ?
	`

	row = DB.QueryRow(queryString, topID)
	err = row.Scan(&track.ID, &track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist,
		&track.AlbumArtist, &track.Composer, &track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal,
		&track.DiscNum, &track.DiscTotal, &track.PictureExt, &track.PictureMIMEType, &track.PictureType,
		&track.PictureDescription, &track.PictureData, &track.Lyrics, &track.Comment, &track.Filepath, &track.Freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByArtistFreq: %v", err)
	}

	return track, nil
}

func GetTopByAlbumFreq() (Track, error) {
	var album string
	var freq float64

	queryString := `
		SELECT Album, SUM(Freq) AS total
		FROM tracks
		WHERE Album IS NOT NULL AND Album != '' AND Album != ' '
		GROUP BY Album
		ORDER BY total DESC
		LIMIT 1;
	`

	row := DB.QueryRow(queryString)
	err := row.Scan(&album, &freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByAlbumFreq: %v", err)
	}

	var topID string

	queryString = `
		SELECT ID
		FROM tracks
		WHERE Album = ?
		GROUP BY ID
		ORDER BY SUM(Freq) DESC
		LIMIT 1;
	`

	row = DB.QueryRow(queryString, album)
	err = row.Scan(&topID)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByAlbumFreq: %v", err)
	}

	var track Track

	queryString = `
		SELECT *
		FROM tracks
		WHERE ID = ?
	`

	row = DB.QueryRow(queryString, topID)
	err = row.Scan(&track.ID, &track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist,
		&track.AlbumArtist, &track.Composer, &track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal,
		&track.DiscNum, &track.DiscTotal, &track.PictureExt, &track.PictureMIMEType, &track.PictureType,
		&track.PictureDescription, &track.PictureData, &track.Lyrics, &track.Comment, &track.Filepath, &track.Freq)
	if err != nil {
		return Track{}, fmt.Errorf("GetTopByAlbumFreq: %v", err)
	}

	return track, nil
}

func GetFreqByGenre() ([]string, []float64, error) {
	var genre []string
	var freq []float64

	queryString := `
		SELECT Genre, SUM(Freq) AS total
		FROM tracks
		WHERE Genre IS NOT NULL AND Genre != '' AND Genre != ' '
		GROUP BY Genre
		ORDER BY total DESC
		LIMIT 10;
	`

	rows, err := DB.Query(queryString)
	if err != nil {
		return nil, nil, fmt.Errorf("GetFreqByGenre: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var g string
		var f float64
		err := rows.Scan(&g, &f)
		if err != nil {
			return nil, nil, fmt.Errorf("GetFreqByGenre: %v", err)
		}
		genre = append(genre, g)
		freq = append(freq, f)
	}
	return genre, freq, nil
}

func GetFreqByArtist() ([]string, []float64, error) {
	var artist []string
	var freq []float64

	queryString := `
		SELECT Artist, SUM(Freq) AS total
		FROM tracks
		WHERE Artist IS NOT NULL AND Artist != '' AND Artist != ' '
		GROUP BY Artist
		ORDER BY total DESC
		LIMIT 10;
	`

	rows, err := DB.Query(queryString)
	if err != nil {
		return nil, nil, fmt.Errorf("GetFreqByArtist: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a string
		var f float64
		err := rows.Scan(&a, &f)
		if err != nil {
			return nil, nil, fmt.Errorf("GetFreqByArtist: %v", err)
		}
		artist = append(artist, a)
		freq = append(freq, f)
	}
	return artist, freq, nil
}

func GetFreqByAlbum() ([]string, []float64, error) {
	var album []string
	var freq []float64

	queryString := `
		SELECT Album, SUM(Freq) AS total
		FROM tracks
		WHERE Album IS NOT NULL AND Album != '' AND Album != ' '
		GROUP BY Album
		ORDER BY total DESC
		LIMIT 10;
	`

	rows, err := DB.Query(queryString)
	if err != nil {
		return nil, nil, fmt.Errorf("GetFreqByAlbum: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a string
		var f float64
		err := rows.Scan(&a, &f)
		if err != nil {
			return nil, nil, fmt.Errorf("GetFreqByAlbum: %v", err)
		}
		album = append(album, a)
		freq = append(freq, f)
	}
	return album, freq, nil
}

func IncrementTrackFreq(track Track) error {
	queryString := `UPDATE tracks SET Freq = Freq + 1 WHERE id = ?`

	_, err := DB.Exec(queryString, &track.ID)

	if err != nil {
		return fmt.Errorf("IncrementTrackFreq: %v", err)
	}

	return nil
}

func BatchAddTracks(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			BatchAddTracks(filepath.Join(dir, file.Name()))
		} else {
			func(filename string) {
				pathToTrack := filepath.Join(dir, filename)
				currTrack := CreateTrackFromFile(pathToTrack, dir)
				_, err := AddTrack(currTrack)
				if err != nil {
					fmt.Println(err)
					return
				}
			}(file.Name())
		}
	}
}

func AddTrack(track Track) (int64, error) {
	if track.Valid == false {
		return 0, fmt.Errorf("AddTrack: %v is not a valid track, skipping", track)
	}

	queryString := `
INSERT INTO Tracks
	(Format, FileType, Title, Album, Artist, AlbumArtist, Composer, Year, Genre, TrackNum, TrackTotal, DiscNum, DiscTotal,
	 PictureExt, PictureMIMEType, PictureType, PictureDescription, PictureData, Lyrics, Comment, Filepath)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(queryString,
		&track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist, &track.AlbumArtist, &track.Composer,
		&track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal, &track.DiscNum, &track.DiscTotal, &track.PictureExt,
		&track.PictureMIMEType, &track.PictureType, &track.PictureDescription, &track.PictureData, &track.Lyrics,
		&track.Comment, &track.Filepath)

	if err != nil {
		return 0, fmt.Errorf("AddTrack: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddTrack: %v", err)
	}
	return id, nil
}

// isImageFile checks if a file is an image based on its extension.
func isImageFile(filename string) bool {
	// List of image file extensions
	imageExtensions := []string{".jpg", ".jpeg", ".JPG", ".JPEG"}

	for _, ext := range imageExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

// getImage checks if a directory contains at least one image file.
func getImage(directory string) (*tag.Picture, error) {
	var picName string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("getImage: %v", err) // Propagate the error
		}
		if !d.IsDir() && isImageFile(d.Name()) {
			picName = d.Name()
			return fs.SkipDir // No need to continue once an image is found
		}
		return nil
	})

	picBytes, _ := os.ReadFile(filepath.Join(directory, picName))
	pic := tag.Picture{
		Ext:         filepath.Ext(picName),
		MIMEType:    "image/jpeg",
		Type:        "",
		Description: "",
		Data:        picBytes,
	}

	return &pic, err
}

func CreateTrackFromFile(filePath string, dir string) Track {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return Track{Valid: false}
	}
	defer f.Close()

	tags, err := ReadTags(f, &filePath)
	if err != nil {
		log.Printf("Error reading tags from file %s: %v", filePath, err)
		return Track{Valid: false}
	}

	// Process tags here...
	//log.Printf("Tags for %s: %v", filePath, tags)

	trackNum, trackTotal := tags.Track()
	discNum, discTotal := tags.Disc()

	track := Track{
		Format:      string(tags.Format()),
		FileType:    string(tags.FileType()),
		Title:       tags.Title(),
		Album:       tags.Album(),
		Artist:      tags.Artist(),
		AlbumArtist: tags.AlbumArtist(),
		Composer:    tags.Composer(),
		Year:        tags.Year(),
		Genre:       tags.Genre(),
		TrackNum:    trackNum,
		TrackTotal:  trackTotal,
		DiscNum:     discNum,
		DiscTotal:   discTotal,
		Lyrics:      tags.Lyrics(),
		Comment:     tags.Comment(),
		Filepath:    filePath,
		Freq:        0,
		Valid:       true,
	}

	pic := tags.Picture()
	if pic != nil {
		track.PictureExt = pic.Ext
		track.PictureMIMEType = pic.MIMEType
		track.PictureType = pic.Type
		track.PictureDescription = pic.Description
		track.PictureData = pic.Data

		return track
	}

	folderPic, err := getImage(dir)
	if err != nil {
		log.Fatal(err)
	}
	track.PictureExt = folderPic.Ext
	track.PictureMIMEType = folderPic.MIMEType
	track.PictureType = folderPic.Type
	track.PictureDescription = folderPic.Description
	track.PictureData = folderPic.Data

	return track
}

func ReadTags(file *os.File, path *string) (tag.Metadata, error) {
	var meta tag.Metadata
	var err error
	switch MFileExt := filepath.Ext(*path); MFileExt {
	case ".flac":
		meta, err = tag.ReadFLACTags(file)
	case ".ogg":
		meta, err = tag.ReadOGGTags(file)
	case ".dsf":
		meta, err = tag.ReadDSFTags(file)
	case ".mp4":
		meta, err = tag.ReadAtoms(file)
	default:
		meta, err = tag.ReadFrom(file)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	return meta, err
}

func GetMusicDirSize() string {
	var sizeBytes int64
	err := filepath.Walk(MusicDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			sizeBytes += info.Size()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}

	sizeGibi := ((float64(sizeBytes) / 1024) / 1024) / 1024
	size := int(math.Round(sizeGibi))
	return fmt.Sprintf("%v GB on disk", size)
}

type Album struct {
	ID    int64
	Title string
	//Artist      string
	AlbumArtist string
	PictureData []byte
}

func GetTracksInAlbum(album Album) ([]Track, error) {
	var tracks []Track

	queryString := `SELECT * FROM tracks WHERE Album = ? ORDER BY TrackNum;`
	rows, err := DB.Query(queryString, &album.Title)
	if err != nil {
		return nil, fmt.Errorf("GetTracksInAlbum: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var track Track
		err := rows.Scan(&track.ID, &track.Format, &track.FileType, &track.Title, &track.Album, &track.Artist,
			&track.AlbumArtist, &track.Composer, &track.Year, &track.Genre, &track.TrackNum, &track.TrackTotal,
			&track.DiscNum, &track.DiscTotal, &track.PictureExt, &track.PictureMIMEType, &track.PictureType,
			&track.PictureDescription, &track.PictureData, &track.Lyrics, &track.Comment, &track.Filepath, &track.Freq)
		if err != nil {
			return nil, fmt.Errorf("GetTracksInAlbum: %v", err)
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

var HomeAlbumsLoaded = false
var HomeAlbumsArr []Album

func LoadHomeAlbums() error {
	albums, err := HomeAlbums()
	if err != nil {
		return err
	}

	HomeAlbumsArr = albums
	HomeAlbumsLoaded = true
	return nil
}

func HomeAlbums() ([]Album, error) {
	var albums []Album

	queryString := `
SELECT t.id AS FirstTrackId, t.Album, t.AlbumArtist, t.PictureData
FROM tracks t
INNER JOIN (
    SELECT MIN(id) AS MinId, Album
    FROM tracks
    GROUP BY Album
) AS sub ON t.id = sub.MinId
ORDER BY t.Album;`
	rows, err := DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("HomeAlbums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.AlbumArtist, &alb.PictureData); err != nil {
			return nil, fmt.Errorf("HomeAlbums %v", err)
		}
		albums = append(albums, alb)
	}

	return albums, nil
}

//// AlbumsByArtist queries for albums that have the specified artist name.
//func AlbumsByArtist(name string) ([]Album, error) {
//	var albums []Album
//
//	rows, err := DB.Query("SELECT * FROM album WHERE artist = ?", name)
//	if err != nil {
//		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
//	}
//	defer rows.Close()
//	// Loop through rows, using Scan to assign column data to struct fields.
//	for rows.Next() {
//		var alb Album
//		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
//			return nil, fmt.Errorf("AlbumsbyArtist %q: %v", name, err)
//		}
//		albums = append(albums, alb)
//	}
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
//	}
//	return albums, nil
//}
//
//// AlbumByID queries for the album with the specified ID.
//func AlbumByID(id int64) (Album, error) {
//	var alb Album
//
//	row := DB.QueryRow("SELECT * FROM album where id = ?", id)
//	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
//		if err == sql.ErrNoRows {
//			return alb, fmt.Errorf("AlbumByID %d: no such album", id)
//		}
//		return alb, fmt.Errorf("AlbumByID %d: %v", id, err)
//	}
//	return alb, nil
//}
//
//// AddAlbum adds the specified album to the database,
//// returning the album ID of the new entry
//func AddAlbum(alb Album) (int64, error) {
//	result, err := DB.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
//	if err != nil {
//		return 0, fmt.Errorf("AddAlbum: %v", err)
//	}
//	id, err := result.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("AddAlbum: %v", err)
//	}
//	return id, nil
//}
