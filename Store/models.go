package store

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"IsCompleted"`
	ListId      int    `json:"list_id"`
}

type List struct {
	Tasks  []Task `json:"tasks"`
	UserId int    `json:"user_id"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
