package osprovider

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/ulikunitz/xz"
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

	return uncompressXZ(resp.Filename, destDir)
}

func uncompressXZ(fileName string, targetDir string) (string, error) {
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader, err := xz.NewReader(file)
	if err != nil {
		return "", err
	}

	xzCutName, _ := strings.CutSuffix(filepath.Base(file.Name()), ".xz")
	outPath := filepath.Join(targetDir, xzCutName)
	out, err := os.Create(outPath)
	if err != nil {
		return "", err
	}

	bufferedWriter := bufio.NewWriter(out)
	defer bufferedWriter.Flush()

	_, err = io.Copy(bufferedWriter, reader)
	if err != nil {
		return "", err
	}

	return outPath, nil
}
