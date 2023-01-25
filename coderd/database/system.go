package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type systemQuerier interface {
	System() System
}

var _ System = (sqlcQuerier)(nil)

// System is a subset of sqlcQuerier that should only be called by system-related functions and never by a user.
type System interface {
	UpdateUserLinkedID(ctx context.Context, arg UpdateUserLinkedIDParams) (UserLink, error)
	UpdateUserLink(ctx context.Context, arg UpdateUserLinkParams) (UserLink, error)
	GetUserLinkByLinkedID(ctx context.Context, linkedID string) (UserLink, error)
	GetUserLinkByUserIDLoginType(ctx context.Context, arg GetUserLinkByUserIDLoginTypeParams) (UserLink, error)
	GetLatestWorkspaceBuilds(ctx context.Context) ([]WorkspaceBuild, error)
	// GetWorkspaceAgentByAuthToken is used in http middleware to get the workspace agent.
	// This should only be used by a system user in that middleware.
	GetWorkspaceAgentByAuthToken(ctx context.Context, authToken uuid.UUID) (WorkspaceAgent, error)
	GetActiveUserCount(ctx context.Context) (int64, error)
	GetAuthorizationUserRoles(ctx context.Context, userID uuid.UUID) (GetAuthorizationUserRolesRow, error)
	GetDERPMeshKey(ctx context.Context) (string, error)
	InsertDERPMeshKey(ctx context.Context, value string) error
	InsertDeploymentID(ctx context.Context, value string) error
	InsertReplica(ctx context.Context, arg InsertReplicaParams) (Replica, error)
	UpdateReplica(ctx context.Context, arg UpdateReplicaParams) (Replica, error)
	DeleteReplicasUpdatedBefore(ctx context.Context, updatedAt time.Time) error
	GetReplicasUpdatedAfter(ctx context.Context, updatedAt time.Time) ([]Replica, error)
	GetTemplates(ctx context.Context) ([]Template, error)
	// UpdateWorkspaceBuildCostByID is used by the provisioning system to update the cost of a workspace build.
	UpdateWorkspaceBuildCostByID(ctx context.Context, arg UpdateWorkspaceBuildCostByIDParams) (WorkspaceBuild, error)
	InsertOrUpdateLastUpdateCheck(ctx context.Context, value string) error
	GetLastUpdateCheck(ctx context.Context) (string, error)
	// Telemetry related functions. These functions are system functions for returning
	// telemetry data. Never called by a user.
	GetWorkspaceBuildsCreatedAfter(ctx context.Context, createdAt time.Time) ([]WorkspaceBuild, error)
	GetWorkspaceAgentsCreatedAfter(ctx context.Context, createdAt time.Time) ([]WorkspaceAgent, error)
	GetWorkspaceAppsCreatedAfter(ctx context.Context, createdAt time.Time) ([]WorkspaceApp, error)
	GetWorkspaceResourcesCreatedAfter(ctx context.Context, createdAt time.Time) ([]WorkspaceResource, error)
	GetWorkspaceResourceMetadataCreatedAfter(ctx context.Context, createdAt time.Time) ([]WorkspaceResourceMetadatum, error)
	DeleteOldAgentStats(ctx context.Context) error
	GetParameterSchemasCreatedAfter(ctx context.Context, createdAt time.Time) ([]ParameterSchema, error)
	GetProvisionerJobsCreatedAfter(ctx context.Context, createdAt time.Time) ([]ProvisionerJob, error)
	// Provisionerd server functions
	InsertWorkspaceAgent(ctx context.Context, arg InsertWorkspaceAgentParams) (WorkspaceAgent, error)
	InsertWorkspaceApp(ctx context.Context, arg InsertWorkspaceAppParams) (WorkspaceApp, error)
	InsertWorkspaceResourceMetadata(ctx context.Context, arg InsertWorkspaceResourceMetadataParams) ([]WorkspaceResourceMetadatum, error)
	AcquireProvisionerJob(ctx context.Context, arg AcquireProvisionerJobParams) (ProvisionerJob, error)
	UpdateProvisionerJobWithCompleteByID(ctx context.Context, arg UpdateProvisionerJobWithCompleteByIDParams) error
	UpdateProvisionerJobByID(ctx context.Context, arg UpdateProvisionerJobByIDParams) error
	InsertProvisionerJob(ctx context.Context, arg InsertProvisionerJobParams) (ProvisionerJob, error)
	InsertProvisionerJobLogs(ctx context.Context, arg InsertProvisionerJobLogsParams) ([]ProvisionerJobLog, error)
	InsertProvisionerDaemon(ctx context.Context, arg InsertProvisionerDaemonParams) (ProvisionerDaemon, error)
	InsertTemplateVersionParameter(ctx context.Context, arg InsertTemplateVersionParameterParams) (TemplateVersionParameter, error)
	InsertWorkspaceResource(ctx context.Context, arg InsertWorkspaceResourceParams) (WorkspaceResource, error)
	InsertParameterSchema(ctx context.Context, arg InsertParameterSchemaParams) (ParameterSchema, error)
}
