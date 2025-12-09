package models

import (
	"SpeechAnalytics/proto/yandex/cloud/ai/stt/v3"

	"gorm.io/gorm"
)

type RoleType string

// Возможные роли
const (
	RoleOperator RoleType = "operator"
	RoleClient   RoleType = "client"
)

// SpeakerStatistics представляет сводную статистику по речи говорящего
type SpeakerStatistics struct {
	gorm.Model // Авто-ID, CreatedAt, UpdatedAt, DeletedAt

	CallID uint `gorm:"column:call_id;not null"` // Внешний ключ, ссылающийся на Call

	SpeakerTag string   `gorm:"column:speaker_tag;"`  // Идентификационный тег говорящего (например, "1", "2")
	Role       RoleType `gorm:"column:role;not null"` // Роль говорящего (operator/client)
	// Основные метрики
	TotalSpeechMs  int64   `gorm:"column:total_speech_ms;"`  // Всего произнесённого времени (в миллисекундах)
	SpeechRatio    float64 `gorm:"column:speech_ratio;"`     // Соотношение речи ко всему времени
	TotalSilenceMs int64   `gorm:"column:total_silence_ms;"` // Общее время тишины (в миллисекундах)
	SilenceRatio   float64 `gorm:"column:silence_ratio;"`    // Соотношение тишины ко всему времени

	// Информационная статистика по словам и буква
	WordsCount       int64   `gorm:"column:words_count;"`        // Общее количество слов
	LettersCount     int64   `gorm:"column:letters_count;"`      // Общее количество букв
	WordsPerSecond   float64 `gorm:"column:words_per_second;"`   // Среднее количество слов в секунду
	LettersPerSecond float64 `gorm:"column:letters_per_second;"` // Средняя скорость набора буквы в секунду

}

func CreateSpeakerStatistics(analysis *stt.SpeakerAnalysis) SpeakerStatistics {
	// Инициализация пустой структуры
	var newStats SpeakerStatistics

	// Присвоим Speaker Tag
	newStats.SpeakerTag = analysis.GetSpeakerTag()

	// Заполняем поля на основе анализа
	newStats.TotalSpeechMs = analysis.GetTotalSpeechMs()
	newStats.SpeechRatio = analysis.GetSpeechRatio()
	newStats.TotalSilenceMs = analysis.GetTotalSilenceMs()
	newStats.SilenceRatio = analysis.GetSilenceRatio()

	// Слова и буквы
	newStats.WordsCount = analysis.GetWordsCount()
	newStats.LettersCount = analysis.GetLettersCount()
	newStats.WordsPerSecond = analysis.GetWordsPerSecond().Mean
	newStats.LettersPerSecond = analysis.GetLettersPerSecond().Mean

	return newStats
}
