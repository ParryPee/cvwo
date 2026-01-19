package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

type SearchResponse struct {
	Posts  []models.Post  `json:"posts"`
	Topics []models.Topic `json:"topics"`
}

type SearchHandler struct {
	DB *sql.DB
}

func (m *SearchHandler) SearchPostAndTopics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}
	posts, err := PostDB.SearchPost(query)
	if err != nil {
		http.Error(w, "Error searching for posts", http.StatusInternalServerError)
		return
	}
	TopicDB := models.TopicDB{DB: m.DB}
	topics, err := TopicDB.SearchTopic(query)
	if err != nil {
		http.Error(w, "Error searching for topics", http.StatusInternalServerError)
		return
	}
	response := SearchResponse{
		Posts:  posts,
		Topics: topics,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
