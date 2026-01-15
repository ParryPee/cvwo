package routers

import (
	"backend/handlers"
	"backend/middleware"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB, jwtkey []byte) http.Handler {
	r := mux.NewRouter()

	topicsHandler := &handlers.TopicHandler{DB: db}
	postHandler := &handlers.PostHandler{DB: db}
	commentHandler := &handlers.CommentHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db, JWTKey: jwtkey}
	authMiddleware := &middleware.AuthMiddleware{JWTKey: jwtkey}

	//Public routes

	//User routes
	r.HandleFunc("/api/users/login", userHandler.Login).Methods("POST")     // User login
	r.HandleFunc("/api/users/register", userHandler.Create).Methods("POST") // Create new user
	r.HandleFunc("/api/users/logout", userHandler.Logout).Methods("POST")   // User logout

	// Public routes that can optionally read user context
	optionalAuth := r.PathPrefix("/api").Subrouter()
	optionalAuth.Use(authMiddleware.OptionalAuthMiddleware)

	//Topic routes
	optionalAuth.HandleFunc("/topics", topicsHandler.GetAllTopics).Methods("GET")   // Get all topics
	optionalAuth.HandleFunc("/topics/{topic_id}", topicsHandler.Get).Methods("GET") // Get topic by ID
	//Comment routes
	optionalAuth.HandleFunc("/posts/{post_id}/comments", commentHandler.GetAllPostComments).Methods("GET") // Get all comments for a post
	optionalAuth.HandleFunc("/comments/{comment_id}", commentHandler.GetCommentByID).Methods("GET")        // Get comment by ID
	// Post routes
	optionalAuth.HandleFunc("/topics/{topic_id}/posts", postHandler.GetAllTopicPosts).Methods("GET") // Get all posts for a topic
	optionalAuth.HandleFunc("/posts/{post_id}", postHandler.GetPostByID).Methods("GET")              // Get Post by ID
	optionalAuth.HandleFunc("/search", postHandler.SearchPost).Methods("GET")
	optionalAuth.HandleFunc("/posts", postHandler.GetAllPosts).Methods("GET")
	//Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.ValidateToken)

	// User routes
	protected.HandleFunc("/users", userHandler.Delete).Methods("DELETE") // Delete user
	protected.HandleFunc("/users/me", userHandler.GetMe).Methods("GET")  // Get current user info

	//Topic routes
	protected.HandleFunc("/topics", topicsHandler.CreateTopic).Methods("POST")              // Create new topic
	protected.HandleFunc("/topics/{topic_id}", topicsHandler.DeleteTopic).Methods("DELETE") // Delete topic by ID
	protected.HandleFunc("/topics/{topic_id}", topicsHandler.UpdateTopic).Methods("PUT")    // Update topic by ID

	//Comment routes
	protected.HandleFunc("/comments", commentHandler.Create).Methods("POST")                        // Create new comment
	protected.HandleFunc("/comments/{comment_id}", commentHandler.Delete).Methods("DELETE")         // Delete comment by ID
	protected.HandleFunc("/comments/{comment_id}", commentHandler.Update).Methods("PUT")            // Update comment by ID
	protected.HandleFunc("/comments/{comment_id}/like", commentHandler.LikeComment).Methods("POST") // Like a comment

	// Post routes
	protected.HandleFunc("/posts/{post_id}", postHandler.Delete).Methods("DELETE") // Delete post by ID
	protected.HandleFunc("/posts", postHandler.Create).Methods("POST")             //Create a new post
	protected.HandleFunc("/posts/{post_id}", postHandler.Update).Methods("PUT")    // Update a post by ID
	protected.HandleFunc("/posts/{post_id}/like", postHandler.LikePost).Methods("POST")

	return r
}
