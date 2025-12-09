package handlers

import "net/http"

func callHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCalls(w, r)
	default:
		errText := "method not allowed"
		http.Error(w, errText, http.StatusMethodNotAllowed)
	}
}

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/calls/", callHandler)

	return mux
}
