package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
)

type Status struct {
	ID          json.Number `json:"id"`
	Description string      `json:"description"`
}

type Data struct {
	StdOut  *string `json:"stdout"`
	Time    string  `json:"time"`
	Memory  int     `json:"memory"`
	StdErr  *string `json:"stderr"`
	Token   string  `json:"token"`
	Message *string `json:"message"`
	Status  Status  `json:"status"`
}

const TypeProcessSubmission = "submission:process"

func CallbackUrl(w http.ResponseWriter, r *http.Request, taskClient *asynq.Client) {
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println("Callback URL hit")

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data: ", err)
		http.Error(w, "Error queuing request", http.StatusInternalServerError)
		return
	}

	task := asynq.NewTask("submission:process", payload)
	info, err := taskClient.Enqueue(task)
	if err != nil {
		log.Println("Error enqueuing task: ", err)
		http.Error(w, "Error queuing request", http.StatusInternalServerError)
		return
	}
	log.Printf("Enqueued task: %+v, Queue: %s", info.ID, info.Queue)

	log.Println("Task enqueued successfully")
	w.WriteHeader(http.StatusOK)
}
