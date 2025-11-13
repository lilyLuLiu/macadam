package osprovider

import (
	"runtime"
	"fmt"
	"os"
	"io"
	"archive/tar"
	"compress/gzip"
	"os/exec"

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

	return resp.Filename, nil
}

func ConvertTarXzToTarGz(inputPath, outputPath string) error {
	// 1. 打开输入文件 (.tar.xz)
	fmt.Printf("1. 打开输入文件: %s\n", inputPath)
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %w", err)
	}
	defer inputFile.Close()

	// 2. 创建 xz 解压器
	fmt.Println("2. 创建 XZ 解压器...")
	xzReader, err := xz.NewReader(inputFile)
	if err != nil {
		return fmt.Errorf("无法创建 XZ 解压器: %w", err)
	}

	// 3. 创建输出文件 (.tar.gz)
	fmt.Printf("3. 创建输出文件: %s\n", outputPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %w", err)
	}
	defer outputFile.Close()

	// 4. 创建 Gzip 压缩器 (写入到输出文件)
	fmt.Println("4. 创建 Gzip 压缩器...")
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close() // 确保在函数退出时关闭 Gzip writer

	// 5. 创建 Tar 读取器和写入器
	fmt.Println("5. 开始流式处理 Tar 文件内容...")
	tarReader := tar.NewReader(xzReader)
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close() // 确保在函数退出时关闭 Tar writer

	// 6. 遍历并复制每个文件
	for {
		// 读取下一个文件头 (Header)
		header, err := tarReader.Next()
		if err == io.EOF {
			// 遍历结束
			break
		}
		if err != nil {
			return fmt.Errorf("读取 tar 头错误: %w", err)
		}

		fmt.Printf("   -> 处理文件/目录: %s\n", header.Name)

		// 写入文件头到新的 tar 包
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("写入 tar 头错误: %w", err)
		}

		// 复制文件内容 (仅对常规文件或符号链接)
		if header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeSymlink {
			if _, err := io.Copy(tarWriter, tarReader); err != nil {
				return fmt.Errorf("复制文件内容错误: %w", err)
			}
		}
		// 其他类型（如目录）没有内容，只需写入头
	}

	fmt.Println("6. 转换完成!")
	return nil
}

func ConvertQcow2ToVDH(inputPath string, outputPath string) error {
	qemuImgPath := "qemu-img"
	args := []string{
		"convert",
		"-f", "qcow2",
		"-O", "vpc",
		"-o", "subformat=dynamic",
		inputPath,
		outputPath,
	}

	// Create the command
	cmd := exec.Command(qemuImgPath, args...)

	fmt.Printf("Starting conversion: %s %v\n", qemuImgPath, args)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Conversion failed: %v\nOutput: %s\n", err, output)
	}

	fmt.Printf("Conversion completed successfully!\nOutput: %s\n", output)
	return nil
}