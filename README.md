# SpeechAnalytics
Сервис умеет обрабатывать .mp3 файлы и обрабатывать в текстовую аналитику по звонкам.

backlog:
ручка для получения данных
___

___
### Инструкция по запуску кода локально
#### Для запуска из исходников:
Добавить .env файл в корень проекта с содержимым:
```
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
API_KEY=ключ к сервисам яндекс аналитики с доступом к SpeechKit
MODEL_URI=url модели для распознавания, например "gpt://b1gj44runl5qb5a5l65n/yandexgpt-lite/latest"
```

`go run ./cmd/main.go`


<!-- дополнительные флаги, примеры .env и так далее. Напишите, какой адрес следует указывать в браузере. -->

Пример для сборки запуска через докер:
```
docker build --tag ya_tasks:v1
docker run -e DB_HOST=192.168.88.6 -e DB_USER=questions -e DB_PASSWORD=0505 -e DB_NAME=steamdb -e API_KEY=AQVNxfr3c5Y7RS58Ueh7OW2xpbb4nHBs7WvEzunY -e MODEL_URI='gpt://b1gj44runl5qb5a5l65n/yandexgpt-lite/latest' speech_analytics
```

Пример для сборки запуска через docker-compose:
```
docker-compose up -d
```