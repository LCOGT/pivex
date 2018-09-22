package pivotal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pivex/pivotal/story"
	"strconv"
)

type Pivotal struct {
	credsPath  string
	pivUrl     string
	projUrl    string
	projectId  int
	apiToken   string
	Iterations []Iteration
	logger     *log.Logger
}

type Iteration struct {
	Number       int           `json:"number"`
	ProjectId    int           `json:"project_id"`
	Length       int           `json:"length"`
	TeamStrength int           `json:"team_strength"`
	Stories      []story.Story `json:"stories"`
	Start        string        `json:"start"`
	Finish       string        `json:"finish"`
	Kind         string        `json:"kind"`
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

func (piv *Pivotal) GetIterations() {
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

	json.NewDecoder(resp.Body).Decode(&piv.Iterations)
}

func (piv *Pivotal) GetStories() {

}
