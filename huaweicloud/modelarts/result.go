package modelarts

type Job struct {
	Metadata JobMetadata `json:"metadata"`
	Status   JobStatus   `json:"status"`
}

type JobMetadata struct {
	Id string `json:"id"`
}

type JobStatus struct {
	Phase     string `json:"phase"`
	Duration  int    `json:"duration"`
	StartTime int    `json:"start_time"`
}
