package domain

import (
	"errors"
)

// ProjectName
type ProjectName interface {
	ProjectName() string
}

func NewProjectName(v string) (ProjectName, error) {
	if v == "" {
		return nil, errors.New("invalid project name")
	}

	return projectName(v), nil
}

type projectName string

func (r projectName) ProjectName() string {
	return string(r)
}
