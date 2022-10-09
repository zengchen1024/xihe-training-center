package domain

import (
	"errors"
	"strings"
)

const (
	resourceProject = "project"
	resourceDataset = "dataset"
	resourceModel   = "model"
)

var (
	ResourceTypeProject ResourceType = resourceType(resourceProject)
	ResourceTypeModel   ResourceType = resourceType(resourceModel)
	ResourceTypeDataset ResourceType = resourceType(resourceDataset)
)

// ResourceType
type ResourceType interface {
	ResourceType() string
}

func NewResourceType(v string) (ResourceType, error) {
	if v != resourceProject && v != resourceModel && v != resourceDataset {
		return nil, errors.New("invalid resource type")
	}

	return resourceType(v), nil
}

type resourceType string

func (s resourceType) ResourceType() string {
	return string(s)
}

// ProjectName
type ProjectName interface {
	ProjectName() string
}

func NewProjectName(v string) (ProjectName, error) {
	if v == "" || !strings.HasPrefix(v, resourceProject) {
		return nil, errors.New("invalid project name")
	}

	return projectName(v), nil
}

type projectName string

func (r projectName) ProjectName() string {
	return string(r)
}
