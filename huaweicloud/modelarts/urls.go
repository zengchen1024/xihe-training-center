package modelarts

import "github.com/chnsz/golangsdk"

const base = "training-jobs"

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(base)
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL(base, jobId)
}

func actionURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL(base, jobId, "actions")
}

func logURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL(base, jobId, "tasks/worker-0/logs/url")
}
