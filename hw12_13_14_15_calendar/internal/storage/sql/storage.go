package sqlstorage

import (
	"context"
	"fmt"
	"os"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4"
)

const (
	dbInsert = "insert into events (userID, Title, Description, EndEvent, StartEvent) values ($1, $2, $3, $4, $5) returning id;"
	dbUpdate = "update events set Title=$1, Description=$2, StartEvent=$3, EndEvent=$4 where ID=$5;"
	dbSelect = "select * from events where userID=$1;"
	dbDelete = "delete from events where ID=$1;"
)

type Storage struct {
	db *pgx.Conn
}

type DBInfo struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func initDB() DBInfo {
	return DBInfo{
		user:     os.Getenv("DATABASE_USER"),
		password: os.Getenv("DATABASE_PASSWORD"),
		host:     os.Getenv("DATABASE_HOST"),
		port:     os.Getenv("DATABASE_PORT"),
		name:     os.Getenv("DATABASE_NAME"),
	}
}

func (db DBInfo) getLink() string {
	info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.host,
		db.port,
		db.user,
		db.password,
		db.name,
	)

	fmt.Println(info)
	return info
}

func (s *Storage) Add(ctx context.Context, event storage.Event) (err error) {
	row := s.db.QueryRow(ctx, dbInsert, event.UserID, event.Title, event.Description, event.End, event.Start)

	return row.Scan(&event.ID)
}

func (s *Storage) List(ctx context.Context, idUser int) (result []storage.Event) {
	rows, err := s.db.Query(ctx, dbSelect, idUser)
	if err != nil {
		fmt.Println(err)
	}

	event := storage.Event{}
	for rows.Next() {
		err = rows.Scan(&event.ID, &event.UserID, &event.Title, &event.Description, &event.End, &event.Start)
		if err == nil {
			result = append(result, event)
		} else {
			fmt.Println(err)
		}
	}

	return result
}

func (s *Storage) Update(ctx context.Context, event storage.Event) (err error) {
	_, err = s.db.Exec(ctx, dbUpdate, event.Title, event.Description, event.Start, event.End, event.ID)

	return err
}

func (s *Storage) Delete(ctx context.Context, id int) (err error) {
	_, err = s.db.Exec(ctx, dbDelete, id)
	return err
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	info := initDB().getLink()

	s.db, err = pgx.Connect(ctx, info)

	return err
}

func (s *Storage) Close(ctx context.Context) (err error) {
	return s.db.Close(ctx)
}
