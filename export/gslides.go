package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/slides/v1"
	"pivex/pivotal"
	"google.golang.org/api/drive/v3"
	"strconv"
)

type GSlides struct {
	credsPath    string
	apiCreds     string
	apiTok       string
	gDriveSrv    *drive.Service
	gsSrv        *slides.Service
	logger       *log.Logger
	forceCreate  bool
	pivIteration pivotal.Iteration
}

const (
	softwareTeamDriveId = "0AEiahtpqanW1Uk9PVA"
	sprintFolder        = "1ScvGIRhj780z_yWzn1ptHovjL-0V9vKK"
)

func New(credsPath string, forceCreate bool, logger *log.Logger, pivIteration pivotal.Iteration) *GSlides {
	apiCreds := fmt.Sprintf("%s/api-creds.json", credsPath)
	apiTok := fmt.Sprintf("%s/api-token.json", credsPath)
	gDriveSrv, gsSrv := getClients(apiCreds, apiTok)

	gs := GSlides{
		credsPath:    credsPath,
		apiCreds:     apiCreds,
		apiTok:       apiTok,
		gDriveSrv:    gDriveSrv,
		gsSrv:        gsSrv,
		logger:       logger,
		forceCreate:  forceCreate,
		pivIteration: pivIteration,
	}

	return &gs
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, apiToken string) *http.Client {
	tok, err := tokenFromFile(apiToken)

	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(apiToken, tok)
	}

	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf(
		"Go to the following link in your browser and copy the authorization code, then paste it in the line "+
			"below\n%v\nAuthorization code: ",
		authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	json.NewEncoder(f).Encode(token)
}

func getClients(apiCreds string, apiTok string) (driveSrv *drive.Service, slidesSrv *slides.Service) {
	b, err := ioutil.ReadFile(apiCreds)
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json
	config, err := google.ConfigFromJSON(
		b,
		"https://www.googleapis.com/auth/presentations.readonly",
		"https://www.googleapis.com/auth/drive",
	)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	driveSrv, err = drive.New(getClient(config, apiTok))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	client := getClient(config, apiTok)

	slidesSrv, err = slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	return
}

func (gs *GSlides) delExisting(slideName string) {
	teamDriveFiles, err := gs.gDriveSrv.Files.
		List().
		SupportsTeamDrives(true).
		IncludeTeamDriveItems(true).
		Corpora("teamDrive").
		TeamDriveId(softwareTeamDriveId).
		Do()

	if err != nil {
		gs.logger.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(teamDriveFiles.Files) == 0 {
		gs.logger.Printf("No files named %s exist", slideName)
	} else {
		for _, teamFile := range teamDriveFiles.Files {
			if teamFile.Name == slideName {
				if gs.forceCreate {
					gs.gDriveSrv.Files.
						Delete(teamFile.Id).
						SupportsTeamDrives(true).
						Do()

					gs.logger.Printf("Deleted filename %s (%s)", teamFile.Name, teamFile.Id)
				} else {
					gs.logger.Printf("Slide %s already exists", teamFile.Name)
					os.Exit(1)
				}
			}
		}
	}
}

func genSlides(stories *[]pivotal.Story) ([]*slides.Request) {
	requests := make([]*slides.Request, 0)

	for _, story := range *stories {
		titleId := fmt.Sprintf("story-title-%d", story.Id)
		bodyId := fmt.Sprintf("story-body-%d", story.Id)

		requests = append(
			requests,
			&slides.Request{
				CreateSlide: &slides.CreateSlideRequest{
					SlideLayoutReference: &slides.LayoutReference{
						PredefinedLayout: "TITLE_AND_BODY",
					},
					PlaceholderIdMappings: []*slides.LayoutPlaceholderIdMapping{
						{
							LayoutPlaceholder: &slides.Placeholder{
								Type: "TITLE",
							},
							ObjectId: titleId,
						},
						{
							LayoutPlaceholder: &slides.Placeholder{
								Type: "BODY",
							},
							ObjectId: bodyId,
						},
					},
				},
			},
			&slides.Request{
				InsertText: &slides.InsertTextRequest{
					ObjectId: titleId,
					Text:     story.Name,
				},
			},
			&slides.Request{
				InsertText: &slides.InsertTextRequest{
					ObjectId: bodyId,
					Text:     story.Description,
				},
			},
		)
	}

	return requests
}

func (gs *GSlides) genSprintAccomplishments() ([]*slides.Request) {
	titleId := "sprint-accomplishments"
	bodyId := "sprint-accomplishments-body"

	totalFeatures, totalChores, totalBugs := gs.getStoryCounts(func(string) (bool) { return true })
	acceptedFeatures, acceptedChores, acceptedBugs := gs.getStoryCounts(
		func(state string) (bool) {
			return state == "accepted"
		})

	requests := make([]*slides.Request, 0)
	requests = append(
		requests,
		&slides.Request{
			CreateSlide: &slides.CreateSlideRequest{
				SlideLayoutReference: &slides.LayoutReference{
					PredefinedLayout: "TITLE_AND_BODY",
				},
				PlaceholderIdMappings: []*slides.LayoutPlaceholderIdMapping{
					{
						LayoutPlaceholder: &slides.Placeholder{
							Type: "TITLE",
						},
						ObjectId: titleId,
					},
					{
						LayoutPlaceholder: &slides.Placeholder{
							Type: "BODY",
						},
						ObjectId: bodyId,
					},
				},
			},
		},
		&slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: titleId,
				Text:     "Sprint accomplishments",
			},
		},
		&slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: bodyId,
				Text: fmt.Sprintf(
					"Total\nFeatures: %d\tChores: %d\tBugs: %d\nAccepted\nFeatures: %d\tChores: %d\tBugs: %d\n", totalFeatures, totalChores, totalBugs, acceptedFeatures, acceptedChores, acceptedBugs),
			},
		},
	)

	return requests
}

func (gs *GSlides) createPres() {
	slideName := "sprint-" + strconv.Itoa(gs.pivIteration.Number)

	gs.delExisting(slideName)

	driveFile := &drive.File{
		MimeType:    "application/vnd.google-apps.presentation",
		Name:        slideName,
		TeamDriveId: softwareTeamDriveId,
		Parents:     []string{sprintFolder},
	}

	presentation, err := gs.gDriveSrv.Files.Create(driveFile).SupportsTeamDrives(true).Do()

	if err != nil {
		gs.logger.Fatalf("Unable to create presentation. %v", err)
	}

	requests := gs.genSprintAccomplishments()

	gs.logger.Printf("Created presentation with ID: %s", driveFile.Id)

	requests = append(requests, genSlides(&gs.pivIteration.Stories)...)

	body := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}
	response, err := gs.gsSrv.Presentations.BatchUpdate(presentation.Id, body).Do()
	if err != nil {
		gs.logger.Fatalf("Unable to create slide. %v", err)
	}

	gs.logger.Printf("Created slide with ID: %s", response.Replies[0].CreateSlide.ObjectId)
}

func (gs *GSlides) Export() {
	gs.createPres()
}

func (gs *GSlides) DelTok() {
	os.Remove(gs.apiCreds)
	os.Remove(gs.apiTok)

	gs.logger.Printf("Deleted authentication files\n%s\n%s", gs.apiCreds, gs.apiTok)
}

func (gs *GSlides) getStoryCounts(constraint func(state string) bool) (featureCount int, choreCount int, bugCount int) {
	for _, story := range gs.pivIteration.Stories {
		storyType := story.StoryType
		storyState := story.CurrentState

		if storyType == "feature" && constraint(storyState) {
			featureCount++
		} else if storyType == "chore" && constraint(storyState) {
			choreCount++
		} else if storyType == "bug" && constraint(storyState) {
			bugCount++
		}
	}

	return
}
