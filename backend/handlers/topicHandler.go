package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TopicHandler struct {
	DB *sql.DB
}

func (m *TopicHandler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	TopicDB := models.TopicDB{DB: m.DB}
	//Size specifies the number of search results to give
	size_str := r.URL.Query().Get("size")
	//Offset specifies the offset for the results, eg. we already served 1 to 10, we want 11 to 20. So the offset should be 10.
	offset_str := r.URL.Query().Get("offset")
	var topics []models.Topic

	var err error

	if size_str != "" && offset_str != "" {
		size, _ := strconv.Atoi(size_str)
		offset, _ := strconv.Atoi(offset_str)
		topics, err = TopicDB.GetByBatch(size, offset)
	} else {
		// If no size or offset is specified just return all the topics, THIS CAN BE VERY SLOW IT MIGHT BE BETTER TO SET A DEFAULT LIMIT
		//FOR FUTURE REFERENCE.
		topics, err = TopicDB.All()
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching topics: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
func (m *TopicHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicIDParam := vars["topic_id"] //Get the topic ID from the request
	//Convert topicID to integer
	topicID, err := strconv.ParseInt(topicIDParam, 10, 64)
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
	var reqBody struct { //Request body that we expect to receive
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   int64  `json:"created_by"`
	}
	//Validate the request inputs
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Title == "" || reqBody.Description == "" {
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
	//Return the newly created topic ID
	json.NewEncoder(w).Encode(map[string]int64{"id": topicID})
}
func (m *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicIDParam := vars["topic_id"]
	//Get the user ID from the context, if not throw an auth error.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}

	var topicID int64
	topicID, err := strconv.ParseInt(topicIDParam, 10, 64)
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
	//Verify that the requester is authorized to delete the topic
	if topic.UserID != currentUserID {
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
	//get the user ID from the context, if unable throw an auth error
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	var topicID int64
	topicID, err := strconv.ParseInt(topicIDParam, 10, 64)
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
	//Verify that the user is authorized to make the changes
	if topic.UserID != currentUserID {
		http.Error(w, "Forbidden: You can only update your own topics", http.StatusForbidden)
		return
	}
	var reqBody struct { //Request body that we expect to receive.
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	//validate the request body
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