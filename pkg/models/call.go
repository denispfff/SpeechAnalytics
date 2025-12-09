package models

import (
	"time"

	"gorm.io/gorm"
)

// StatusType определяет возможные статусы обработки файла
type StatusType string

// Возможные статусы обработки
const (
	NewProcessingStatus       StatusType = "new"
	ProcessingStatus          StatusType = "processing"
	CompletedProcessingStatus StatusType = "completed"
	FailedProcessingStatus    StatusType = "failed"
)

type Call struct {
	gorm.Model
	FileName     string     `db:"file_name" json:"file_name"`                                    // Название файла
	RequestGID   string     `db:"request_gid" json:"request_gid"`                                // Уникальный идентификатор запроса (GID)
	Status       StatusType `db:"status" json:"status"`                                          // Текущий статус обработки
	Date         time.Time  `gorm:"column:date;not null" json:"date"`                            // Дата звонка
	OperatorName string     `gorm:"column:operator_name;size:255;not null" json:"operator_name"` // Имя оператора
	ClientNumber string     `gorm:"column:client_number;size:20;not null" json:"client_number"`  // Номер телефона клиента
	// SummarisationID uint       `gorm:"column:summarisation_id;null" json:"summarisation_id"`        // Внешний ключ суммирования (если доступно)
	CallDuration uint `gorm:"column:call_duration;null" json:"call_duration"`

	SpeakerStatistics []SpeakerStatistics `gorm:"foreignKey:CallID;references:ID"` // Many-to-one relation
}
