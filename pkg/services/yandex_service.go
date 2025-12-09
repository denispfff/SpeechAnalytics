package services

import (
	"SpeechAnalytics/pkg/logger"
	"SpeechAnalytics/pkg/models"
	"SpeechAnalytics/pkg/repositories"
	"SpeechAnalytics/proto/yandex/cloud/ai/stt/v3"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	yaApi       = "stt.api.cloud.yandex.net:443"
	textSchema  = `{"type":"object","properties":{"jsonSchema":{"type":"array","items":[{"type":"object","properties":{"calls":{"type":"object","properties":{"type":{"type":"string"},"call_duration":{"type":"integer"},"result":{"type":"string"},"satisfaction_level":{"type":"integer"}},"required":["type","call_duration","result","satisfaction_level"]}},"required":["calls"]},{"type":"object","properties":{"linguistic_data":{"type":"object","properties":{"pauses_count":{"type":"integer"},"pauses_duration":{"type":"integer"},"topics":{"type":"string"},"obscene_language_operator":{"type":"string"},"obscene_language_client":{"type":"string"},"rate_of_speech":{"type":"string"},"characteristic":{"type":"string"},"fillers":{"type":"string"},"speech_rate":{"type":"string"},"emotional_tone":{"type":"string"},"communication_techniques":{"type":"string"},"interruptions_client":{"type":"integer"},"interruptions_operator":{"type":"integer"}},"required":["pauses_count","pauses_duration","topics","obscene_language_operator","obscene_language_client","rate_of_speech","characteristic","fillers","speech_rate","emotional_tone","communication_techniques","interruptions_client","interruptions_operator"]}},"required":["linguistic_data"]},{"type":"object","properties":{"script_compliance":{"type":"object","properties":{"speech_analysis":{"type":"object","properties":{"greetings_message":{"type":"boolean"},"farewell_message":{"type":"boolean"},"offer_of_assistance":{"type":"boolean"},"company_name":{"type":"boolean"},"decision_confirmation":{"type":"boolean"},"apology_text":{"type":"boolean"}},"required":["greetings_message","farewell_message","offer_of_assistance","company_name","decision_confirmation","apology_text"]}},"required":["speech_analysis"]}},"required":["script_compliance"]}]}},"required":["jsonSchema"]}`
	instruction = `
Проанализируй текст, полученный из аудиозаписи, по следующим критериям:
- Звонок входящий или исходящий,
- Длительность звонка,
- Общий результат звонка,
- Уровень удовлетворенности клиента в процентах,
- Количество пауз в диалоге,
- Длительность пауз в диалоге,
- Темы диалога,
- Выдели нецензурную речь оператора,
- Выдели нецензурную речь клиента,
- Темп речи,
- Ключевые фразы и слова диалога,
- Список слов-паразитов от оператора,
- Скорость речи,
- Эмоциональный окрас разговора,
- Коммуникационные техники используемые оператором,
- Количество перебиваний речи клиентом,
- Количество перебиваний речи оператором,
- Поздоровался ли оператор,
- Попрощался ли оператор,
- Была ли предложена помощь оператором,
- Было ли произнесено наименовании компании,
- Подтверждено ли что проблема решена,
- Извинился ли оператор в диалоге 
Результат выдай в формате json с соответствующими наименованиями полей`
)

var rpcClient *grpc.ClientConn
var protobufSchema *structpb.Struct
var auth string
var modelURI string
var limiter = rate.NewLimiter(rate.Every(time.Second/5), 5) // 5 запросов в секунду

func waitForRateLimit() {
	limiter.Wait(context.Background()) // Ждёт, пока разрешён следующий запрос
}

// Инициализируем общее gRPC-соединение и готовим схему для запросов
func InitConnect(apiKey, URI string) {
	auth = apiKey
	modelURI = URI
	var err error
	var schema map[string]interface{}
	err = json.Unmarshal([]byte(textSchema), &schema)
	if err != nil {
		logger.Logger.Fatalln(err)
	}
	protobufSchema, err = structpb.NewStruct(schema)
	if err != nil {
		logger.Logger.Fatalln(err)
	}

	tlsCreds := credentials.NewClientTLSFromCert(nil, "")
	// Создаем клиентское соединение
	rpcClient, err = grpc.NewClient(
		yaApi,
		grpc.WithTransportCredentials(tlsCreds),
	)
	if err != nil {
		log.Fatalln("Не удалось подключиться к Ya")
	}
}

func SendToRecognize(call *models.Call) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Api-Key "+auth)
	fileByte, err := repositories.GetFileByte(call)
	if err != nil {
		return err
	}
	req := prepareRpcBody(fileByte)

	client := stt.NewAsyncRecognizerClient(rpcClient)

	resp, err := client.RecognizeFile(ctx, req)
	if err == nil {
		call.RequestGID = resp.GetId()
	}

	return err
}

func prepareRpcBody(fileBytes []byte) *stt.RecognizeFileRequest {
	req := &stt.RecognizeFileRequest{
		AudioSource: &stt.RecognizeFileRequest_Content{
			Content: fileBytes,
		},
		RecognitionModel: &stt.RecognitionModelOptions{
			Model: "deferred-general",
			AudioFormat: &stt.AudioFormatOptions{
				AudioFormat: &stt.AudioFormatOptions_ContainerAudio{
					ContainerAudio: &stt.ContainerAudio{
						ContainerAudioType: stt.ContainerAudio_MP3, // Указываем тип контейнера (MP3)
					},
				},
			},
			TextNormalization: &stt.TextNormalizationOptions{
				TextNormalization: stt.TextNormalizationOptions_TEXT_NORMALIZATION_ENABLED,
			},
			LanguageRestriction: &stt.LanguageRestrictionOptions{
				RestrictionType: stt.LanguageRestrictionOptions_WHITELIST,
				LanguageCode:    []string{"ru-RU"},
			},
			AudioProcessingType: stt.RecognitionModelOptions_FULL_DATA,
		},
		SpeechAnalysis: &stt.SpeechAnalysisOptions{
			EnableSpeakerAnalysis:      true,
			EnableConversationAnalysis: true,
			// https://yandex.cloud/ru/docs/speechkit/stt/analysis#statistics
			DescriptiveStatisticsQuantiles: []float64{0.5, 0.9}, // задает процентили, по которым SpeechKit рассчитывает количественные показатели распределения характеристик речи.

		},
		SpeakerLabeling: &stt.SpeakerLabelingOptions{
			SpeakerLabeling: stt.SpeakerLabelingOptions_SPEAKER_LABELING_ENABLED,
		},
		Summarization: &stt.SummarizationOptions{
			ModelUri: modelURI,
			Properties: []*stt.SummarizationProperty{
				{
					Instruction: instruction,
					ResponseFormat: &stt.SummarizationProperty_JsonSchema{
						JsonSchema: &stt.JsonSchema{
							Schema: protobufSchema,
						},
					},
				},
			},
		},
	}
	return req
}

func GetRecognition(call *models.Call) (models.Responses, error) {
	var responses models.Responses

	stream, err := sendRecognitionRequest(call)
	if err != nil {
		return responses, err
	}
	defer stream.CloseSend()

	return processStream(call, stream)
}

func sendRecognitionRequest(file *models.Call) (grpc.ServerStreamingClient[stt.StreamingResponse], error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Api-Key "+auth)

	req := &stt.GetRecognitionRequest{
		OperationId: file.RequestGID,
	}

	client := stt.NewAsyncRecognizerClient(rpcClient)

	waitForRateLimit()
	stream, err := client.GetRecognition(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func processStream(call *models.Call, stream grpc.ServerStreamingClient[stt.StreamingResponse]) (models.Responses, error) {
	res := models.Responses{Status: false}
	count := 0
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Printf("СЧЕТЧИК СТАТИСТИКИ: %d", count)
			if count == 4 {
				call.Status = models.CompletedProcessingStatus
			}
			break
		}

		if err != nil {
			switch status.Code(err) {
			case codes.ResourceExhausted:
				log.Println(err)
				return res, nil
			case codes.NotFound:
				if strings.Contains(err.Error(), "is not ready") {
					return res, nil
				}
			default:
				return res, err
			}
		}

		if response.GetSpeakerAnalysis().GetWindowType() == stt.SpeakerAnalysis_TOTAL {
			// fmt.Println(response.GetSpeakerAnalysis(), "\n-")
			count++
			call.SpeakerStatistics = append(call.SpeakerStatistics, models.CreateSpeakerStatistics(response.GetSpeakerAnalysis()))
		}

		if analysis := response.GetConversationAnalysis(); analysis != nil {
			count++
			fmt.Println(analysis, "\n-")
		}
		if summ := response.GetSummarization(); summ != nil {
			count++
			summStr := summ.Results[0].GetResponse()
			fmt.Println(summStr, "\n-")
		}
	}

	return res, nil
}
