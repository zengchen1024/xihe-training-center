package mysql

const (
	fieldStatus     = "status"
	fieldVersion    = "version"
	fieldLastCommit = "last_commit"
)

var projectTableName string

type ProjectRepoSyncLock struct {
	Id         int    `json:"-"            gorm:"column:id"`
	Owner      string `json:"-"            gorm:"column:owner"`
	RepoId     string `json:"-"            gorm:"column:repo_id"`
	Status     string `json:"status"       gorm:"column:status"`
	Version    int    `json:"-"            gorm:"column:version"`
	LastCommit string `json:"last_commit"  gorm:"column:last_commit"`
}

func (r *ProjectRepoSyncLock) TableName() string {
	return projectTableName
}
