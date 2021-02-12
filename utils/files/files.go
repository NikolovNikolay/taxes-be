package files

import (
	"fmt"
	"github.com/h2non/filetype"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"strings"
)

func GetExtension(file multipart.File, fileName string) (string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		logrus.Errorf("error reading file: %s", fileName)
		return "", err
	}
	kind, _ := filetype.Match(buff)
	if kind == filetype.Unknown {
		return "", err
	}

	return kind.Extension, nil
}

func GetMultipartFileExtension(fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("could not process file header")
	}
	defer file.Close()
	return GetExtension(file, fh.Filename)
}

func GetExtensionFromName(name string) string {
	components := strings.Split(name, ".")
	return components[len(components)-1]
}
