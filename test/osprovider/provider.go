package osprovider

import (
	"runtime"

	"github.com/cavaliergopher/grab/v3"
)

type OsProvider interface {
	Fetch(destDir string) error
}

func kernelArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return "invalid"
	}
}

func downloadOS(destDir, url string) (string, error) {
	// https://github.com/cavaliergopher/grab/issues/104
	grab.DefaultClient.UserAgent = "macadam"
	resp, err := grab.Get(destDir, url)
	if err != nil {
		return "", err
	}

	return resp.Filename, nil
}
