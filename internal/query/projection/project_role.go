package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectRoleProjection struct {
	crdb.StatementHandler
}

const ProjectRoleProjectionTable = "zitadel.projections.project_roles"

func NewProjectRoleProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectRoleProjection {
	p := &ProjectRoleProjection{}
	config.ProjectionName = ProjectRoleProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectRoleProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.RoleAddedType,
					Reduce: p.reduceProjectRoleAdded,
				},
				{
					Event:  project.RoleChangedType,
					Reduce: p.reduceProjectRoleChanged,
				},
				{
					Event:  project.RoleRemovedType,
					Reduce: p.reduceProjectRoleRemoved,
				},
			},
		},
	}
}

const (
	ProjectRoleColumnProjectID     = "project_id"
	ProjectRoleColumnKey           = "role_key"
	ProjectRoleColumnCreationDate  = "creation_date"
	ProjectRoleColumnChangeDate    = "change_date"
	ProjectRoleColumnResourceOwner = "resource_owner"
	ProjectRoleColumnSequence      = "sequence"
	ProjectRoleColumnDisplayName   = "display_name"
	ProjectRoleColumnGroupName     = "group_name"
	ProjectRoleColumnCreator       = "creator_id"
)

func (p *ProjectRoleProjection) reduceProjectRoleAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Fmre5", "seq", event.Sequence(), "expectedType", project.RoleAddedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g92Fg", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectRoleColumnKey, e.Key),
			handler.NewCol(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectRoleColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectRoleColumnSequence, e.Sequence()),
			handler.NewCol(ProjectRoleColumnDisplayName, e.DisplayName),
			handler.NewCol(ProjectRoleColumnGroupName, e.Group),
			handler.NewCol(ProjectRoleColumnCreator, e.EditorUser()),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-M0fwg", "seq", event.Sequence(), "expectedType", project.GrantChangedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-sM0f", "reduce.wrong.event.type")
	}
	if e.DisplayName == nil && e.Group == nil {
		return crdb.NewNoOpStatement(e), nil
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnSequence, e.Sequence()),
			handler.NewCol(ProjectRoleColumnDisplayName, *e.DisplayName),
			handler.NewCol(ProjectRoleColumnGroupName, *e.Group),
		},
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MlokF", "seq", event.Sequence(), "expectedType", project.GrantRemovedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-L0fJf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}