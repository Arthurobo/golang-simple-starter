package posts

type CreatePostModel struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

func (p *CreatePostModel) Validate() []string {
	var errors []string
	if p.Title == "" {
		errors = append(errors, "title is required")
	}
	if p.Content == "" {
		errors = append(errors, "content is required")
	}
	return errors
}

type UpdatePostModel struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type AllPostsModel struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type IndividualPostModel struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type DeletePostModel struct {
	ID        int  `json:"id"`
	IsActive  bool `json:"is_active"`
	IsDeleted bool `json:"is_deleted"`
}
