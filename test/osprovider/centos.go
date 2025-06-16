package osprovider

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type CentosProvider struct {
	diskImage string
}

func NewCentosProvider() *CentosProvider {
	return &CentosProvider{}
}

func (centos *CentosProvider) Fetch(destDir string) (string, error) {
	log.Infof("downloading centos to %s", destDir)
	arch := kernelArch()
	centosURL := fmt.Sprintf("https://cloud.centos.org/centos/10-stream/%s/images/CentOS-Stream-GenericCloud-10-20250324.0.%s.qcow2", arch, arch)
	file, err := downloadOS(destDir, centosURL)
	if err != nil {
		return "", err
	}

	centos.diskImage = file

	return file, nil
}
