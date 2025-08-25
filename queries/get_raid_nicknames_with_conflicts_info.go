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
					  COALESCE(fp.result::jsonb, '{}'::jsonb),
					  '{nickname_ids}',
					  CASE
						WHEN jsonb_typeof(fp.result::jsonb->'nickname_ids') = 'array' THEN
						  COALESCE(
							(
							  SELECT jsonb_agg(
									   jsonb_build_object('id', nick_id, 'name', n.name)
									 )
							  FROM jsonb_array_elements_text(
									 COALESCE(fp.result::jsonb->'nickname_ids', '[]'::jsonb)
								   ) AS nick_id
							  LEFT JOIN aa_nicknames n ON n.id::text = nick_id
							),
							'[]'::jsonb
						  )
						ELSE '[]'::jsonb
					  END,
					  true
					),
					'{conflicts}',
					CASE
					  WHEN jsonb_typeof(fp.result::jsonb->'conflicts') = 'array' THEN
						COALESCE(
						  (
							SELECT jsonb_agg(
									 jsonb_build_object(
									   'box', c->'box',
									   'similar',
										 COALESCE(
										   (
											 SELECT jsonb_agg(
													  jsonb_build_object('id', sim_id, 'name', n.name)
													)
											 FROM jsonb_array_elements_text(
													COALESCE(c->'similar', '[]'::jsonb)
												  ) AS sim_id
											 LEFT JOIN aa_nicknames n ON n.id::text = sim_id
										   ),
										   '[]'::jsonb
										 )
									 )
								   )
							FROM jsonb_array_elements(
								   COALESCE(fp.result::jsonb->'conflicts', '[]'::jsonb)
								 ) AS c
						  ),
						  '[]'::jsonb
						)
					  ELSE '[]'::jsonb
					END,
					true
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
