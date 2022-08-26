package domain

type UserTraining struct {
	User Account

	Training
}

type Training struct {
	ProjectId string

	Name TrainingName
	Desc TrainingDesc

	CodeDir  Directory
	BootFile FilePath
	LogDir   Directory

	Hypeparameters []KeyValue
	Env            []KeyValue
	Inputs         []KeyValue
	Outputs        []KeyValue

	Compute Compute
}

type Compute struct {
	Type    ComputeType
	Version ComputeVersion
	Flavor  ComputeFlavor
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type TrainingDetail struct {
	Status   TrainingStatus
	Duration TrainingDuration
}
