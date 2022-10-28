package watch

const (
	workAim    = "aim"
	workGetLog = "log"
	workOutput = "output"
)

type TrainingInfo struct {
	User       string
	ProjectId  string
	TrainingId string

	JobId     string
	LogDir    string
	AimDir    string
	OutputDir string

	// TODO if timeout, ignore this work and set status to timeout
	//timeout      int

	result       Data
	notifyFailed bool
}

type JobDetail struct {
	Status   string
	Duration int
}

type TrainingService interface {
	IsDone(string) bool

	IsSucess(string) bool

	GetDetail(string) (JobDetail, error)

	// GetLogFilePath return the obs path of log
	GetLogFilePath(logDir string) (string, error)

	// GenOutput generates the zip file of output dir and
	// return the obs path of that file.
	GenOutput(outputDir string) (string, error)

	// GenAim generates the zip file of aim dir
	// and return the obs path of that file.
	GenAim(aimDir string) (string, error)
}

type WatchService interface {
	AddTraining(*TrainingInfo) error
}
