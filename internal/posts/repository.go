package posts

import "database/sql"

func GetAllPosts(db *sql.DB) ([]AllPostsModel, error) {
	rows, err := db.Query(`
		SELECT id, title, content, user_id FROM posts
		WHERE is_active = TRUE AND is_deleted = FALSE
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []AllPostsModel
	for rows.Next() {
		var u AllPostsModel
		if err := rows.Scan(&u.ID, &u.Title, &u.Content, &u.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, u)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, id string) (IndividualPostModel, error) {
	var p IndividualPostModel
	err := db.QueryRow("SELECT id, title, content, user_id FROM posts WHERE id = $1 AND is_active = TRUE AND is_deleted = FALSE", id).Scan(&p.ID, &p.Title, &p.Content, &p.UserID)
	if err != nil {
		return p, err
	}
	return p, nil
}

func CreatePost(db *sql.DB, p *CreatePostModel) error {
	return db.QueryRow("INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", p.Title, p.Content, p.UserID).Scan(&p.ID)
}

func UpdatePost(db *sql.DB, id string, p *UpdatePostModel) (bool, error) {
	result, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3 AND is_active = TRUE AND is_deleted = FALSE", p.Title, p.Content, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func DeletePost(db *sql.DB, id string) (IndividualPostModel, error) {
	p, err := GetPostByID(db, id)
	if err != nil {
		return p, err
	}

	_, err = db.Exec(`
		UPDATE posts
		SET is_active = FALSE,
		    is_deleted = TRUE,
		    last_updated = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = TRUE AND is_deleted = FALSE
	`, id)

	return p, err
}
