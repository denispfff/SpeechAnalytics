package repositories

import (
	"SpeechAnalytics/pkg/database"
	"SpeechAnalytics/pkg/models"
	"strings"
	"time"
)

func CreateCall(fileName string) (models.Call, error) {
	var newCall models.Call

	newCall.FileName = fileName
	newCall.Status = models.NewProcessingStatus

	created, err := time.Parse("2006-01-02__15-04-05", fileName[:20])
	if err != nil {
		created = time.Now()
	}

	newCall.Date = created

	tmp := strings.Split(fileName, "__")
	newCall.OperatorName = strings.Split(tmp[3], "@")[0]
	newCall.ClientNumber = tmp[2]

	newCallDB := database.DB.Db.Create(&newCall)
	return newCall, newCallDB.Error
}

func GetCallByStatus(status models.StatusType) ([]*models.Call, error) {
	var calls []*models.Call
	err := database.DB.Db.Where("status = ?", status).Find(&calls).Error
	return calls, err
}

func SaveCall(call *models.Call) error {
	return database.DB.Db.Save(call).Error

}
