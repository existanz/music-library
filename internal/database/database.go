package database

import (
	"database/sql"
	"fmt"
	"log"
	"music-library/internal/models"
	"music-library/internal/server/query"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error
	AddNewSong(song models.Song) error
	GetSongs(opts query.Options) ([]models.Song, error)
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = MigrateUp(db)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	dbInstance = &service{
		db: db,
	}
	FillTestData(dbInstance)
	return dbInstance
}

func (s *service) AddNewSong(song models.Song) error {
	id, err := s.AddNewArtist(song.Group)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO songs (artist_id, song, release_date, lirycs, link) VALUES ($1, $2, $3, $4, $5)", id, song.Song, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddNewArtist(artist string) (int, error) {
	var id int
	rows, err := s.db.Query("SELECT id FROM artists WHERE artist = $1 LIMIT 1", artist)
	if err == nil && rows.Next() {
		rows.Scan(&id)
		return id, nil
	}

	err = s.db.QueryRow("INSERT INTO artists (artist) VALUES ($1) RETURNING id", artist).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) GetSongs(opts query.Options) ([]models.Song, error) {
	var songs []models.Song
	query := "SELECT songs.id, artist, song, release_date, lirycs, link FROM songs LEFT JOIN artists ON songs.artist_id = artists.id" + getFilersString(opts.Filters) + getPaginatorString(opts.Paginator)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.Id, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func getFilersString(filters []query.Filter) string {
	if len(filters) == 0 {
		return ""
	}

	filtersString := " WHERE "
	for _, filter := range filters {
		filtersString += filter.Field + " = '" + filter.Value + "' AND "
	}
	filtersString = strings.TrimSuffix(filtersString, "AND ")
	return filtersString
}

func getPaginatorString(paginator query.Paginator) string {
	return " LIMIT " + fmt.Sprint(paginator.Limit) + " OFFSET " + fmt.Sprint(paginator.Offset)
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func MigrateUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Can't create driver for migration : %v\n", err)
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
