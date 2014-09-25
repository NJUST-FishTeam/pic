package picfs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func SavePic(data []byte, storePath string) (string, error) {
	md5_hash := fmt.Sprintf("%x", md5.Sum(data))
	suffix := ""
	if IsJPEG(data) {
		suffix = ".jpg"
	} else if IsPNG(data) {
		suffix = ".png"
	} else {
		return "", errors.New("The file is not a picture")
	}

	first := md5_hash[0:2]
	second := md5_hash[2:4]
	fileName := md5_hash[4:] + suffix
	filePath := path.Join(first, second, fileName)

	err := os.MkdirAll(path.Join(storePath, first, second), os.ModePerm)
	if err != nil {
		log.Fatal("Mkdir: ", err.Error())
		return "", errors.New("Can not make dir")
	}

	err = ioutil.WriteFile(path.Join(storePath, filePath), data, 0644)
	if err != nil {
		log.Fatal("WriteFile: ", err.Error())
		return "", errors.New("Can not write the file")
	}

	return filePath, nil
}
