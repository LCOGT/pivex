package story

type Story struct {
	Kind          string  `json:"kind"`
	Id            int     `json:"id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Estimate      int     `json:"estimate"`
	StoryType     string  `json:"story_type"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	CurrentState  string  `json:"current_state"`
	RequestedById int     `json:"requested_by_id"`
	Url           string  `json:"url"`
	ProjectId     int     `json:"project_id"`
	OwnerIds      []int   `json:"owner_ids"`
	Labels        []Label `json:"labels"`
	OwnedById     int     `json:"owned_by_id"`
}

type Label struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
