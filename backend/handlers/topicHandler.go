package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type TopicHandler struct {
	DB *sql.DB
}

func (m *TopicHandler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	TopicDB := models.TopicDB{DB: m.DB}
	topics, err := TopicDB.All()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching topics: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
func (m *TopicHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicIDParam := vars["topic_id"]
	var topicID int64
	_, err := fmt.Sscanf(topicIDParam, "%d", &topicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}
	TopicDB := models.TopicDB{DB: m.DB}
	topic, err := TopicDB.GetByID(topicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching topic: %v", err), http.StatusInternalServerError)
		return
	}
	if topic == nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

func (m *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   int64  `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Title == "" || reqBody.Description == "" || reqBody.CreatedBy == 0 {
		http.Error(w, "Title, Description, and CreatedBy are required", http.StatusBadRequest)
		return
	}
	TopicDB := models.TopicDB{DB: m.DB}
	topicID, err := TopicDB.Create(reqBody.Title, reqBody.Description, reqBody.CreatedBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating topic: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": topicID})
}
func (m *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicIDParam := vars["topic_id"]
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var topicID int64
	_, err := fmt.Sscanf(topicIDParam, "%d", &topicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}
	TopicDB := models.TopicDB{DB: m.DB}
	topic, err := TopicDB.GetByID(topicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching topic: %v", err), http.StatusInternalServerError)
		return
	}
	if topic == nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}
	if topic.CreatedBy != currentUserID {
		http.Error(w, "Forbidden: You can only delete your own topics", http.StatusForbidden)
		return
	}
	err = TopicDB.Delete(topicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting topic: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *TopicHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicIDParam := vars["topic_id"]
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var topicID int64
	_, err := fmt.Sscanf(topicIDParam, "%d", &topicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}
	TopicDB := models.TopicDB{DB: m.DB}
	topic, err := TopicDB.GetByID(topicID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching topic: %v", err), http.StatusInternalServerError)
		return
	}
	if topic == nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}
	if topic.CreatedBy != currentUserID {
		http.Error(w, "Forbidden: You can only update your own topics", http.StatusForbidden)
		return
	}
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Title == "" || reqBody.Description == "" {
		http.Error(w, "Title and Description are required", http.StatusBadRequest)
		return
	}
	err = TopicDB.Update(topicID, reqBody.Title, reqBody.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating topic: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
