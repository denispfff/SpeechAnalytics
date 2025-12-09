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
	rootPath = "./tmp/"
	postfix  = ".mp3"
)

func InitFilePaths() {
	for _, status := range []models.StatusType{models.New, models.Processing, models.Success, models.Error} {
		path := filepath.Join(rootPath, string(status))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				log.Printf("ошибка создания директории '%s': %v\n", path, err)
			}
		}
	}
}

func ProcessNewFiles() error {
	newPath := filepath.Join(rootPath, string(models.New))
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

	call, err := CreateCall(d.Name())

	if err != nil {
		log.Printf("ошибка добавления нового файла в бд: %v", err)
		return err
	}

	err = MoveByStatus(&call, models.New, models.Processing)
	if err != nil {
		log.Printf("ошибка перемещения файла: %v", err)
		return err
	}
	return nil
}

func GetFileByte(f *models.Call) ([]byte, error) {
	path := filepath.Join(rootPath, string(f.Status), f.FileName)

	fileByte, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileByte.Close()

	return io.ReadAll(fileByte)
}

func MoveByStatus(c *models.Call, oldStatus, newStatus models.StatusType) error {

	oldPath := filepath.Join(rootPath, string(oldStatus), c.FileName)
	newPath := filepath.Join(rootPath, string(newStatus), c.FileName)

	return os.Rename(oldPath, newPath)
}
