package bot

type Plan struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	IsAllDay    bool   `json:"is_all_day"`
	CreatedAt   string `json:"created_at"`
}

type PlanRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	IsAllDay    bool   `json:"is_all_day"`
}
