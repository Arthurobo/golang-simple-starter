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
	r.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(getAllPostsHandler(db)))).Methods("GET")
	r.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(createPostHandler(db)))).Methods("POST")
	r.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(getPostDetailsHandler(db)))).Methods("GET")
	r.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(updatePostHandler(db)))).Methods("PUT")
	r.Handle("/posts/{id}/delete", middleware.AuthMiddleware(http.HandlerFunc(deletePostHandler(db)))).Methods("PUT")
}

func getAllPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		posts, err := GetAllPosts(db)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching posts", err.Error())
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Posts fetched successfully", posts)
	}
}

func getPostDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetAuthenticatedUserID(w, r)
		id := mux.Vars(r)["id"]
		post, err := GetPostByID(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found", nil)
			return
		}
		if post.UserID != userID {
			utils.WriteJSONError(w, http.StatusForbidden, "Forbidden", nil)
			return
		}
		utils.WriteJSONSuccess(w, http.StatusOK, "Post fetched successfully", post)
	}
}

func createPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetAuthenticatedUserID(w, r)

		var p CreatePostModel
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request payload", nil)
			return
		}

		p.UserID = userID

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
		// Extract the authenticated user ID from context
		userID := middleware.GetAuthenticatedUserID(w, r)
		fmt.Println("Authenticated User ID:", userID)

		// Get the post ID from URL params
		id := mux.Vars(r)["id"]

		// Decode the request payload into the update model
		var post UpdatePostModel
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid payload", nil)
			return
		}

		// Assign the userID from token to the post struct
		post.UserID = userID
		fmt.Println("Post User ID:", post.UserID)

		// Check ownership and update the post
		ok, err := UpdatePost(db, id, &post)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Update failed", err.Error())
			return
		}
		if !ok {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found or already deleted", nil)
			return
		}

		// Set the correct post ID in the response
		if postID, err := strconv.Atoi(id); err == nil {
			post.ID = postID
		}

		utils.WriteJSONSuccess(w, http.StatusOK, "Post updated successfully", post)
	}
}

func deletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the authenticated user ID from context
		userID := middleware.GetAuthenticatedUserID(w, r)
		fmt.Println("Authenticated User ID:", userID)

		// Get the post ID from the URL
		id := mux.Vars(r)["id"]

		// Fetch the post to check ownership
		post, err := GetPostByID(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found", nil)
			return
		}

		// Ensure only the post owner can delete it
		if post.UserID != userID {
			utils.WriteJSONError(w, http.StatusForbidden, "You are not authorized to delete this post", nil)
			return
		}

		// Proceed to delete
		_, err = DeletePost(db, id)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to delete post", err.Error())
			return
		}

		utils.WriteJSONSuccess(w, http.StatusOK, "Post deleted successfully", nil)
	}
}
