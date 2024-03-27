package utils

import (
	"go.uber.org/zap"
	"io"
	"os"
)

func CopyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
func WriteStringToFile(content, filePath string) error {
	// 将字符串转换为字节数组
	data := []byte(content)

	// 使用 os.WriteFile 将字节数组写入文件
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	zap.L().Info(content + "has been written to file： " + filePath)
	return nil
}
