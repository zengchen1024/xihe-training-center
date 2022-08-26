package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reName      = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_-/]+$")
	reFilePath  = regexp.MustCompile("^[a-zA-Z0-9_-/.]+$")

	config = Config{}
)

func Init(cfg TrainingConfig) {
	config = Config{cfg}
}

type Config struct {
	Training TrainingConfig
}

type TrainingConfig struct {
	MaxNameLength int `json:"max_name_length"`
	MinNameLength int `json:"min_name_length"`
	MaxDescLength int `json:"max_desc_length"`
}

func (r *TrainingConfig) Setdefault() {
	if r.MaxNameLength == 0 {
		r.MaxNameLength = 50
	}

	if r.MinNameLength == 0 {
		r.MinNameLength = 5
	}

	if r.MaxDescLength == 0 {
		r.MaxDescLength = 100
	}
}

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !reName.MatchString(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// TrainingName
type TrainingName interface {
	TrainingName() string
}

func NewTrainingName(v string) (TrainingName, error) {
	max := config.Training.MaxNameLength
	min := config.Training.MinNameLength

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return trainingName(v), nil
}

type trainingName string

func (r trainingName) TrainingName() string {
	return string(r)
}

// TrainingDesc
type TrainingDesc interface {
	TrainingDesc() string
}

func NewTrainingDesc(v string) (TrainingDesc, error) {
	if v == "" {
		return nil, nil
	}

	max := config.Training.MaxDescLength
	if len(v) > max {
		return nil, fmt.Errorf("the length of desc should be less than %d", max)
	}

	return trainingDesc(v), nil
}

type trainingDesc string

func (r trainingDesc) TrainingDesc() string {
	return string(r)
}

// Directory
type Directory interface {
	Directory() string
}

func NewDirectory(v string) (Directory, error) {
	if !reDirectory.MatchString(v) {
		return nil, errors.New("invalid directory")
	}

	return directory(v), nil
}

type directory string

func (r directory) Directory() string {
	return string(r)
}

// FilePath
type FilePath interface {
	FilePath() string
}

func NewFilePath(v string) (FilePath, error) {
	if !reFilePath.MatchString(v) {
		return nil, errors.New("invalid filePath")
	}

	return filePath(v), nil
}

type filePath string

func (r filePath) FilePath() string {
	return string(r)
}

// ComputeType
type ComputeType interface {
	ComputeType() string
}

func NewComputeType(v string) (ComputeType, error) {
	if v == "" {
		return nil, errors.New("invalid compute type")
	}

	return computeType(v), nil
}

type computeType string

func (r computeType) ComputeType() string {
	return string(r)
}

// ComputeVersion
type ComputeVersion interface {
	ComputeVersion() string
}

func NewComputeVersion(v string) (ComputeVersion, error) {
	if v == "" {
		return nil, errors.New("invalid compute version")
	}

	return computeVersion(v), nil
}

type computeVersion string

func (r computeVersion) ComputeVersion() string {
	return string(r)
}

// ComputeFlavor
type ComputeFlavor interface {
	ComputeFlavor() string
}

func NewComputeFlavor(v string) (ComputeFlavor, error) {
	if v == "" {
		return nil, errors.New("invalid compute flavor")
	}

	return computeFlavor(v), nil
}

type computeFlavor string

func (r computeFlavor) ComputeFlavor() string {
	return string(r)
}

// CustomizedKey
type CustomizedKey interface {
	CustomizedKey() string
}

func NewCustomizedKey(v string) (CustomizedKey, error) {
	if v == "" {
		return nil, errors.New("invalid key")
	}

	return customizedKey(v), nil
}

type customizedKey string

func (r customizedKey) CustomizedKey() string {
	return string(r)
}

// CustomizedValue
type CustomizedValue interface {
	CustomizedValue() string
}

func NewCustomizedValue(v string) (CustomizedValue, error) {
	return customizedValue(v), nil
}

type customizedValue string

func (r customizedValue) CustomizedValue() string {
	return string(r)
}

// TrainingRegion
type TrainingRegion interface {
	TrainingRegion() string
}

func NewTrainingRegion(v string) (TrainingRegion, error) {
	if v == "" {
		return nil, errors.New("invalid key")
	}

	return trainingRegion(v), nil
}

type trainingRegion string

func (r trainingRegion) TrainingRegion() string {
	return string(r)
}

// TrainingStatus
type TrainingStatus interface {
	TrainingStatus() string
}

func NewStatusCreating() TrainingStatus {
	return trainingStatus("Creating")
}

func NewStatusPending() TrainingStatus {
	return trainingStatus("Pending")
}

func NewStatusRunning() TrainingStatus {
	return trainingStatus("Running")
}

func NewStatusFailed() TrainingStatus {
	return trainingStatus("Failed")
}

func NewStatusCompleted() TrainingStatus {
	return trainingStatus("Completed")
}

func NewStatusTerminating() TrainingStatus {
	return trainingStatus("Terminating")
}

func NewStatusTerminated() TrainingStatus {
	return trainingStatus("Terminated")
}

func NewStatusAbnormal() TrainingStatus {
	return trainingStatus("Abnormal")
}

type trainingStatus string

func (s trainingStatus) TrainingStatus() string {
	return string(s)
}

// TrainingDuration
type TrainingDuration interface {
	TrainingDuration() int
}

func NewTrainingDuration(t int) (TrainingDuration, error) {
	if t < 0 {
		return nil, errors.New("invalid training time")
	}

	return trainingDuration(t), nil
}

type trainingDuration int

func (t trainingDuration) TrainingDuration() int {
	return int(t)
}
