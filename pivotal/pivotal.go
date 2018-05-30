package pivotal

import (
	"net/http"
	"encoding/json"
	"strconv"
	"log"
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
)

type Pivotal struct {
	credsPath    string
	pivUrl       string
	projUrl      string
	projectId    int
	apiToken     string
	Intervals    []Interval
	logger       *log.Logger
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

const (
	pivUrl    = "https://www.pivotaltracker.com/services/v5/projects"
	projectId = 1314272
)

// TODO: See if Go has something like kwargs, doesn't look like New functions can be overloaded
func New(apiToken string, credsPath string, logger *log.Logger) *Pivotal {
	apiTokenFile := fmt.Sprintf("%s/pivotal-token", credsPath)

	if apiToken == "" {
		apiToken = readTokenFile(apiTokenFile)
	} else {
		writeTokenFile(apiTokenFile, apiToken)
	}

	piv := Pivotal{
		credsPath: credsPath,
		pivUrl:    pivUrl,
		projUrl:   pivUrl + "/" + strconv.Itoa(projectId),
		projectId: projectId,
		apiToken:  apiToken,
		logger:    logger,
	}

	return &piv
}

func readTokenFile(filePath string) (apiToken string) {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to read token file %s: %v", filePath, err)
	}

	scanner := bufio.NewScanner(f)
	// TODO: Only reads one line
	for scanner.Scan() {
		apiToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func writeTokenFile(filePath string, apiToken string) {
	fData := []byte(apiToken + "\n")
	ioutil.WriteFile(filePath, fData, 0600)
}

func (piv *Pivotal) GetStories() {
	req, err := http.NewRequest("GET", piv.projUrl+"/iterations?scope=current", nil)
	if err != nil {
		piv.logger.Fatalf("Error creating request: %s", err)
	}

	req.Header.Set("X-TrackerToken", piv.apiToken)

	resp, err := client.Do(req)

	if err != nil {
		piv.logger.Fatalf("Error sending request: %s", err)
	}

	if resp.StatusCode != 200 {
		piv.logger.Fatalf("Error getting Pivotal data: %s", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&piv.Intervals)
}
