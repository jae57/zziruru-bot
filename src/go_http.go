package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Request struct {
	Object string
	Entry []Messaging
}

type Messaging struct {
	Messaging []Message
}

type Message struct {
	Message string
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		const verifyToken = "<YOUR_VERIFY_TOKEN>"

		mode := r.FormValue("hub.mode")
		token := r.FormValue("hub.verify_token")
		challenge := r.FormValue("hub.challenge")

		if (mode == "subscribe") && (token == verifyToken) {
			fmt.Println("WEBHOOK_VERIFIED")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, challenge)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	case "POST":
		req := Request{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Object == "page" {
			webhookEvent := req.Entry[0].Messaging[0].Message
			fmt.Println(webhookEvent)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "EVENT_RECEIVED")
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

func main() {
	// 기본 Url 핸들러 메서드 지정
	http.HandleFunc("/webhook", webhookHandler)
	// 서버 시작
	err := http.ListenAndServe(":1337", nil)
	// 예외 처리
	if err == nil {
		fmt.Println("ListenAndServe Started! -> Port(1337)")
	} else {
		log.Fatal("ListenAndServe: ",err)
	}
}