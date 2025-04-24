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
	archInName := "-" + arch
	if archInName == "-x86_64" {
		archInName = ""
	}
	var centosURL = fmt.Sprintf("https://cloud.centos.org/centos/10-stream/%s/images/CentOS-Stream-ec2%s-10-20250324.0.%s.raw.xz", arch, archInName, arch)
	file, err := downloadOS(destDir, centosURL)
	if err != nil {
		return "", err
	}

	centos.diskImage = file

	return file, nil
}
