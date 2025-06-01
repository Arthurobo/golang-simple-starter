package posts

import (
	"api/pkg/middleware"
	"api/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, db *sql.DB) {
	// r.HandleFunc("/posts", getAllPostsHandler(db)).Methods("GET")
	r.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(getAllPostsHandler(db)))).Methods("GET")
	r.HandleFunc("/posts", createPostHandler(db)).Methods("POST")
	r.HandleFunc("/posts/{id}", getPostHandler(db)).Methods("GET")
	r.HandleFunc("/posts/{id}", updatePostHandler(db)).Methods("PUT")
	r.HandleFunc("/posts/{id}/delete", deletePostHandler(db)).Methods("PUT")
}

func getAllPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := middleware.GetAuthenticatedUserID(w, r)
		fmt.Println("Authenticated User ID:", userID)

		posts, err := GetAllPosts(db)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching posts", err.Error())
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Posts fetched successfully", posts)
	}
}

func getPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post, err := GetPostByID(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found", nil)
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Post fetched successfully", post)
	}
}

func createPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p CreatePostModel
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		if errs := p.Validate(); len(errs) > 0 {
			utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed", errs)
			return
		}

		if err := CreatePost(db, &p); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to create post", err.Error())
			return
		}
		utils.WriteJSONSuccess(w, http.StatusCreated, "Post created successfully", p)

	}
}

func updatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var post UpdatePostModel

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid payload", nil)
			return
		}

		ok, err := UpdatePost(db, id, &post)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Update failed", err.Error())
			return
		}
		if !ok {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found or already deleted", nil)
			return
		}

		// ✅ Assign the correct ID back to the post model before returning
		if postID, err := strconv.Atoi(id); err == nil {
			post.ID = postID
		}

		utils.WriteJSONSuccess(w, http.StatusOK, "Post updated successfully", post)
	}
}

func deletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		_, err := DeletePost(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found", nil)
			return
		}

		utils.WriteJSONSuccess(w, http.StatusOK, "Post deleted successfully", nil)
	}
}
