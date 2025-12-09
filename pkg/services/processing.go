package services

import (
	"SpeechAnalytics/pkg/models"
	"SpeechAnalytics/pkg/repositories"
	"context"
	"log"
	"sync"
	"time"
)

func StartProcessing(ctx context.Context, auth, modelURI string) {
	InitConnect(auth, modelURI)
	defer rpcClient.Close()

	filesCheckTicker := time.NewTicker(30 * time.Second)
	sendTicker := time.NewTicker(60 * time.Second)
	chechStatusTicker := time.NewTicker(30 * time.Second)

	go RpcGetRecognition(&ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Завершение работы...")
			return

		case <-filesCheckTicker.C:
			log.Println("Проверка наличия новых файлов...")
			go repositories.ProcessNewFiles()

		case <-sendTicker.C:
			log.Println("Отправка файлов на обработку...")
			go RpcSendNew(&ctx)

		case <-chechStatusTicker.C:
			log.Println("Проверка файлов в обработке...")
			go RpcGetRecognition(&ctx)
		}
	}
}

func RpcSendNew(ctx *context.Context) {
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	var calls []*models.Call
	calls, err := repositories.GetCallByStatus(models.Processing)
	if err != nil {
		log.Printf("ошибка получения новых файлов")
		return
	}

	if len(calls) == 0 {
		log.Printf("новых файлов не найдено")
		return
	}

	for _, call := range calls {
		wg.Add(1)
		go func(call *models.Call) {
			defer wg.Done()

			call.Status = models.Processing

			if err := SendToRecognize(call); err != nil {
				call.Status = models.Error
				log.Println(err)
			}

			mutex.Lock()
			defer mutex.Unlock()

			err := repositories.SaveCall(call)
			if err != nil {
				log.Println("Ошибка обновления файла:", err)
			}
		}(call)
	}

	wg.Wait()
}

func RpcGetRecognition(ctx *context.Context) {
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	var calls []*models.Call
	calls, err := repositories.GetCallByStatus(models.Processing)
	if err != nil {
		log.Printf("ошибка получения новых запросов")
		return
	}

	if len(calls) == 0 {
		log.Printf("новых запросов не найдено")
		return
	}

	for _, call := range calls {
		wg.Add(1)
		go func(call *models.Call) {
			defer wg.Done()
			_, err := GetRecognition(call)

			if err != nil {
				call.Status = models.Error
				log.Println(err)
				err = repositories.MoveByStatus(call, models.Processing, models.Error)
				if err != nil {
					log.Println("Ошибка перемещения файла:", err)
				}
			}

			// В случае, если статус не изменился - пропускаем обработку
			if call.Status == models.Processing {
				return
			}

			mutex.Lock()
			defer mutex.Unlock()

			if err = repositories.SaveCall(call); err != nil {
				log.Println("Ошибка обновления файла:", err)
			}

			err = repositories.MoveByStatus(call, models.Processing, models.Success)
			if err != nil {
				log.Println("Ошибка перемещения файла:", err)
			}
		}(call)
		wg.Wait()
	}
}
