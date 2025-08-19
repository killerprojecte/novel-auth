package repository

import (
	"auth/.gen/auth/public/model"
	. "auth/.gen/auth/public/table"
	"database/sql"
	"encoding/json"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

type Event = model.AuthEvent

type EventFilter struct {
	ActorUser     string
	TargetUser    string
	Action        string
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

type EventRepository interface {
	List(filter EventFilter, pageNumber, pageSize int64) ([]*Event, error)
	Save(action string, detail interface{}) error
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

	if filter.ActorUser != "" {
		stmt = stmt.WHERE(RawBool("detail ->> 'actor_user' = $user",
			map[string]interface{}{"$user": filter.ActorUser}))
	}
	if filter.TargetUser != "" {
		stmt = stmt.WHERE(RawBool("detail ->> 'target_user' = $user",
			map[string]interface{}{"$user": filter.TargetUser}))
	}
	if filter.Action != "" {
		stmt = stmt.WHERE(AuthEvent.Action.EQ(String(filter.Action)))
	}
	if !filter.CreatedAfter.IsZero() {
		stmt = stmt.WHERE(AuthEvent.CreatedAt.GT(TimestampzT(filter.CreatedAfter)))
	}
	if !filter.CreatedBefore.IsZero() {
		stmt = stmt.WHERE(AuthEvent.CreatedAt.LT(TimestampzT(filter.CreatedBefore)))
	}

	stmt = stmt.
		ORDER_BY(AuthUser.ID.ASC()).
		LIMIT(pageSize).
		OFFSET(pageNumber * pageSize)

	var dest []*Event
	err := stmt.Query(r.db, &dest)
	if err == qrm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return dest, nil
}

func (r *eventRepository) Save(action string, detail interface{}) error {
	detailEncoded, _ := json.Marshal(detail)
	event := &Event{
		Action:    action,
		Detail:    string(detailEncoded),
		CreatedAt: time.Now(),
	}
	stmt := AuthEvent.INSERT(AuthEvent.MutableColumns).
		MODEL(event)

	_, err := stmt.Exec(r.db)
	return err
}
