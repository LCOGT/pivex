package pivotal

import (
	"net/http"
	"encoding/json"
	"os"
	"strconv"
)

type Pivotal struct {
	pivUrl    string
	projUrl   string
	projectId int
	apiToken  string
	Intervals []Interval
}

type Interval struct {
	Number       int     `json:"number"`
	ProjectId    int     `json:"project_id"`
	Length       int     `json:"length"`
	TeamStrength int     `json:"team_strength"`
	Stories      []Story `json:"stories"`
	Start        string  `json:"start"`
	Finish       string  `json:"finish"`
	Kind         string  `json:"kind"`
}

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

var client = http.Client{

}

func New() *Pivotal {
	pivUrl := "https://www.pivotaltracker.com/services/v5/projects"
	projectId := 1314272

	piv := Pivotal{
		pivUrl:    pivUrl,
		projUrl:   pivUrl + "/" + strconv.Itoa(projectId),
		projectId: projectId,
		// TODO: Load from file and arg instead
		apiToken: os.Getenv("PIVOTAL_API_TOKEN"),
	}

	return &piv
}

func (piv *Pivotal) GetStories() {
	req, err := http.NewRequest("GET", piv.projUrl+"/iterations?scope=current", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-TrackerToken", piv.apiToken)

	resp, err := client.Do(req)

	json.NewDecoder(resp.Body).Decode(&piv.Intervals)

	if err != nil {
		panic(err)
	}
}
