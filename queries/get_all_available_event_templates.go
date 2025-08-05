package queries

import (
	"context"
	"database/sql"
	db "github.com/OfficialEvsty/aa-data/db/interface"
)

type EventTemplateDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetAllAvailableEventTemplates struct {
	exec db.ISqlExecutor
}

func NewGetAllAvailableEventTemplates(exec db.ISqlExecutor) *GetAllAvailableEventTemplates {
	return &GetAllAvailableEventTemplates{
		exec: exec,
	}
}

func (q *GetAllAvailableEventTemplates) Handle(ctx context.Context) ([]*EventTemplateDTO, error) {
	var eventTemplates []*EventTemplateDTO
	query := `SELECT id, name FROM aa_template_events`
	rows, err := q.exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var eventTemplate EventTemplateDTO
		if err = rows.Scan(&eventTemplate.ID, &eventTemplate.Name); err != nil {
			return nil, err
		}
		eventTemplates = append(eventTemplates, &eventTemplate)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return eventTemplates, nil
}

func (q *GetAllAvailableEventTemplates) WithTx(tx *sql.Tx) *GetAllAvailableEventTemplates {
	return &GetAllAvailableEventTemplates{
		tx,
	}
}
