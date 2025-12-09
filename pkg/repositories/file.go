package repositories

import (
	"SpeechAnalytics/pkg/models"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	rootPath      = "./tmp/"
	newPath       = "./tmp/new/"
	inProcessPath = "./tmp/in_process/"
	donePath      = "./tmp/done/"

	postfix = ".mp3"
)

func InitFilePaths() {
	paths := []string{rootPath, newPath, inProcessPath, donePath}
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				log.Printf("ошибка создания директории '%s': %v\n", path, err)
			}
		}
	}
}

func ProcessNewFiles() error {
	err := filepath.WalkDir(newPath, getNewFile)

	return err
}

func getNewFile(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() || !strings.HasSuffix(d.Name(), postfix) {
		return nil
	}

	_, err = CreateCall(d.Name())

	if err != nil {
		log.Printf("ошибка добавления нового файла в бд: %v", err)
		return err
	}

	err = os.Rename(path, filepath.Join(filepath.Dir(inProcessPath), d.Name()))
	if err != nil {
		log.Printf("ошибка перемещения файла: %v", err)
		return err
	}
	return nil
}

func GetFileByte(f *models.Call) ([]byte, error) {
	path := filepath.Join(filepath.Dir(inProcessPath), f.FileName)

	fileByte, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileByte.Close()

	return io.ReadAll(fileByte)
}

func MoveDone(c *models.Call) error {
	oldPath := filepath.Join(inProcessPath, c.FileName)
	newPath := filepath.Join(donePath, c.FileName)

	return os.Rename(oldPath, newPath)
}
