package handlers

import "net/http"

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	//	mux.HandleFunc("/questions/", questionsHandler)

	return mux
}
