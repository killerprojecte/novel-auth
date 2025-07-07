package repository

import (
	"auth/.gen/auth/public/model"
	. "auth/.gen/auth/public/table"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
)

type Event = model.AuthEvent

type EventFilter struct {
	Action        *string
	App           *string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
}

type EventRepository interface {
	List(filter EventFilter, pageNumber, pageSize int64) ([]*Event, error)
	Save(event *Event) error
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) List(filter EventFilter, pageNumber, pageSize int64) ([]*Event, error) {
	stmt := SELECT(AuthEvent.AllColumns).
		FROM(AuthEvent)

	if filter.Action != nil {
		stmt = stmt.WHERE(AuthEvent.Action.EQ(String(*filter.Action)))
	}
	if filter.App != nil {
		stmt = stmt.WHERE(AuthEvent.App.EQ(String(*filter.App)))
	}
	if filter.CreatedAfter != nil {
		stmt = stmt.WHERE(AuthEvent.CreatedAt.GT(TimestampzT(*filter.CreatedAfter)))
	}
	if filter.CreatedBefore != nil {
		stmt = stmt.WHERE(AuthEvent.CreatedAt.LT(TimestampzT(*filter.CreatedBefore)))
	}

	stmt = stmt.
		ORDER_BY(AuthUser.ID.ASC()).
		LIMIT(pageSize).
		OFFSET(pageNumber * pageSize)

	var dest []*Event
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (r *eventRepository) Save(event *Event) error {
	stmt := AuthEvent.INSERT(AuthUser.MutableColumns).
		MODEL(event)

	_, err := stmt.Exec(r.db)
	return err
}
