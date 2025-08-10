package queries

import (
	"context"
	db "github.com/OfficialEvsty/aa-data/db/interface"
	"github.com/OfficialEvsty/aa-data/domain/usecase"
	"github.com/google/uuid"
)

type GetRaidParticipantsInfoQuery struct {
	exec db.ISqlExecutor
}

func NewGetRaidParticipantsInfoQuery(exec db.ISqlExecutor) *GetRaidParticipantsInfoQuery {
	return &GetRaidParticipantsInfoQuery{exec}
}

func (q *GetRaidParticipantsInfoQuery) Handle(ctx context.Context, publishID uuid.UUID) (*usecase.RaidParticipantsWithS3Data, error) {
	var dto usecase.RaidParticipantsWithS3Data
	query := `SELECT
				p.s3,
				jsonb_set(
				  jsonb_set(
					fp.result::jsonb,
					ARRAY['nickname_ids'],
					CASE 
					  WHEN jsonb_typeof(fp.result::jsonb->'nickname_ids') = 'array' THEN
						(
						  SELECT jsonb_agg(
							jsonb_build_object(
							  'id', nick_id,
							  'name', n.name
							)
						  )
						  FROM jsonb_array_elements_text(fp.result::jsonb->'nickname_ids') nick_id
						  LEFT JOIN aa_nicknames n ON n.id::text = nick_id
						)
					  ELSE '[]'::jsonb
					END
				  ),
				  ARRAY['conflicts'],
				  CASE
					WHEN jsonb_typeof(fp.result::jsonb->'conflicts') = 'array' THEN
					  (
						SELECT jsonb_agg(
						  jsonb_build_object(
							'box', c->'box',
							'similar', (
							  SELECT jsonb_agg(
								jsonb_build_object(
								  'id', sim_id,
								  'name', n.name
								)
							  )
							  FROM jsonb_array_elements_text(c->'similar') sim_id
							  LEFT JOIN aa_nicknames n ON n.id::text = sim_id
							)
						  )
						)
						FROM jsonb_array_elements(fp.result::jsonb->'conflicts') c
					  )
					ELSE '[]'::jsonb
				  END
				) AS new_data
			FROM finished_publishes fp
			JOIN publishes p ON p.id = fp.publish_id
			WHERE p.id = $1;`
	err := q.exec.QueryRowContext(
		ctx,
		query,
		publishID,
	).Scan(
		&dto.Snapshot,
		&dto.IssuedParticipants,
	)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}
