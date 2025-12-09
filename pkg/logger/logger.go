package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

// func Init() {
// 	Logger = log.New(
// 		os.Stdout, "",
// 		// "server: ",
// 		log.LstdFlags|log.Lshortfile,
// 	)
// }

func Init() {
	log.SetOutput(os.Stdout)                     // Стандартный вывод
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Отображать дату-время и короткий путь к файлу
	log.SetPrefix("[GLOBAL]: ")                  // Префикс логов
}
