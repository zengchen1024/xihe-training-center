package trainingimpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	libutils "github.com/opensourceways/community-robot-lib/utils"
)

func (s *helper) GetLogFilePath(logDir string) (p string, err error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = s.bucket
	input.Prefix = logDir // "src0/"

	output, err := s.obsClient.ListObjects(input)
	if err != nil {
		return
	}

	if v := output.Contents; len(v) > 0 {
		p = v[0].Key
	}

	return
}

func (s *helper) GenOutput(outputDir string) (string, error) {
	return s.uploadFolder(outputDir)
}

func (s *helper) GenAim(aimDir string) (string, error) {
	return s.uploadFolder(aimDir)
}

func (s *helper) uploadFolder(obsPath string) (string, error) {
	if obsPath == "" {
		return "", nil
	}

	tempDir, err := ioutil.TempDir(s.suc.UploadWorkDir, "upload")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(tempDir)

	params := []string{
		s.suc.UploadFolderShell, tempDir,
		s.suc.OBSUtilPath, s.bucket, obsPath,
	}

	v, err, _ := libutils.RunCmd(params...)
	if err != nil {
		err = fmt.Errorf(
			"run upload folder shell, err=%s, params=%v",
			err.Error(), params,
		)

		return "", err
	}

	return strings.TrimSuffix(string(v), "\n"), nil
}
