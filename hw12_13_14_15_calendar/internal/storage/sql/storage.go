package sqlstorage

import (
	"context"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"os"
	"time"

	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4"
)

const (
	//For server
	dbInsert       = "select new_event(userName:=$1, title:=$2, description:=$3, notify:=$4, startEvent:=$5, endEvent:=$6);"
	dbUpdate       = "update events set title=$1, description=$2, notify=$3, startEvent=$4, endEvent=$5 where events.id=$6;"
	dbSelect       = "select events.id, users.id, name, title, description, notify, startEvent, endEvent from users, events where userID = (select users.id where name=$1);"
	dbSelectByTime = "select events.id, users.id, name, title, description, notify, startEvent, endEvent from users, events where userID = (select users.id where name=$1) and startEvent between $2 and $3;"
	dbDelete       = "delete from events where id=$1;"

	//For scheduler
	//I don`t use between in this because notify send in queue twice
	dbSelectByNotify = "select * from events where notify >= $1 and notify < $2;"
	dbClear          = "delete from events where endEvent < $1;"
)

type Storage struct {
	db *pgx.Conn
}

var logger app.Logger

type DBInfo struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func initDB() DBInfo {
	return DBInfo{
		user:     os.Getenv("USER"),
		password: os.Getenv("PASSWORD"),
		host:     os.Getenv("HOST"),
		port:     os.Getenv("PORT"),
		name:     os.Getenv("NAME"),
	}
}

func (db DBInfo) getLink() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.host,
		db.port,
		db.user,
		db.password,
		db.name,
	)
}

func (s *Storage) Add(ctx context.Context, event *storage.Event) error {
	rows := s.db.QueryRow(ctx, dbInsert,
		event.User.Name,
		event.Title,
		event.Description,
		event.GetNotifyTime().Round(time.Minute).UTC(),
		event.Start.Round(time.Minute).UTC(),
		event.End.Round(time.Minute).UTC(),
	)
	if err := rows.Scan(&event.ID); err != nil {
		logger.Error("[ERR] While scan event id: ", err)
		return err
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, event *storage.Event) (err error) {
	_, err = s.db.Exec(ctx, dbUpdate,
		event.Title,
		event.Description,
		event.GetNotifyTime().Round(time.Minute).UTC(),
		event.Start.Round(time.Minute).UTC(),
		event.End.Round(time.Minute).UTC(),
		event.ID,
	)

	return err
}

func (s *Storage) Delete(ctx context.Context, id int64) (err error) {
	_, err = s.db.Exec(ctx, dbDelete, id)
	return err
}

func New(logg app.Logger) *Storage {
	logger = logg
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = pgx.Connect(ctx, initDB().getLink())
	return err
}

func (s *Storage) Close(ctx context.Context) (err error) {
	return s.db.Close(ctx)
}

func getEventList(rows pgx.Rows) ([]storage.Event, error) {
	var (
		event storage.Event
		date  time.Time
		err   error
	)
	events := make([]storage.Event, 0, 1)

	for rows.Next() {
		err = rows.Scan(
			&event.ID,
			&event.User.ID,
			&event.User.Name,
			&event.Title,
			&event.Description,
			&date,
			&event.Start,
			&event.End,
		)
		if err != nil {
			logger.Error("[ERR] While scaning event: ", err)
			return nil, err
		}

		event.Notify = int32(event.Start.Sub(date).Minutes())

		events = append(events, event)
	}

	return events, err
}

func (s *Storage) Clear(ctx context.Context) error {
	_, err := s.db.Exec(ctx, dbClear, time.Now().Add(-time.Hour*24*30*12).UTC())
	return err
}

func (s *Storage) ListUpcoming(ctx context.Context, userName string, until time.Duration) ([]storage.Event, error) {
	now := time.Now().UTC().Round(time.Minute)

	now, err := time.ParseInLocation(time.RFC3339Nano, now.Format(time.RFC3339Nano), time.UTC)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, dbSelectByTime, userName, now, now.Add(until))
	if err != nil {
		logger.Error("[ERR] DB select query: ", err)
		return nil, err
	}
	defer rows.Close()

	return getEventList(rows)
}

func (s *Storage) List(ctx context.Context, userName string) ([]storage.Event, error) {
	rows, err := s.db.Query(ctx, dbSelect, userName)
	if err != nil {
		logger.Error("[ERR] DB select query: ", err)
		return nil, err
	}
	defer rows.Close()

	return getEventList(rows)
}

func (s *Storage) ListByNotify(ctx context.Context, until time.Duration) ([]storage.Event, error) {
	current := time.Now().Round(time.Minute).UTC()
	end := current.Add(until).UTC()

	rows, err := s.db.Query(ctx, dbSelectByNotify, current, end)
	if err != nil {
		logger.Error("[ERR] DB select by notify query: ", err)
		return nil, err
	}
	defer rows.Close()

	return getEventList(rows)
}
