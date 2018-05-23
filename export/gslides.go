package gslides

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
	"os/user"
)

type GSlides struct {
	apiCreds string
	apiToken string
}

var (
	gDriveSrv *drive.Service
	gsSrv     *slides.Service

	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	googleApiRoot = func() string {
		usr, err := user.Current()
		if err != nil {
			logger.Fatal(err)
		}

		return fmt.Sprintf("%s/.google-api", usr.HomeDir)
	}()
	apiCreds = googleApiRoot + "/pivex-creds.json"
	apiToken = googleApiRoot + "/pivex-token.json"
)

func New() *GSlides {
	gDriveSrv, gsSrv = getClients()

	gs := GSlides{}

	return &gs
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	if _, err := os.Stat(googleApiRoot); os.IsNotExist(err) {
		os.Mkdir(googleApiRoot, 0600)
	}

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
		"Go to the following link in your browser and copy the authorization code, then paste it in the line " +
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
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func getClients() (driveSrv *drive.Service, slidesSrv *slides.Service) {
	b, err := ioutil.ReadFile(apiCreds)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(
		b,
		"https://www.googleapis.com/auth/presentations.readonly",
		"https://www.googleapis.com/auth/drive",
	)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	driveSrv, err = drive.New(getClient(config))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	client := getClient(config)

	slidesSrv, err = slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	return
}

func delExisting(slideName string) {
	r, err := gDriveSrv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()

	if err != nil {
		logger.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(r.Files) == 0 {
		logger.Printf("No files named %s exist", slideName)
	} else {
		for _, i := range r.Files {
			if i.Name == slideName {
				gDriveSrv.Files.Delete(i.Id).Do()

				logger.Printf("Deleted filename %s (%s)", i.Name, i.Id)
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

func createPres(pivInterval *pivotal.Interval) {
	slideName := "sprint-" + strconv.Itoa(pivInterval.Number)

	delExisting(slideName)

	p := &slides.Presentation{
		Title: slideName,
	}
	presentation, err := gsSrv.Presentations.Create(p).Fields(
		"presentationId",
	).Do()
	if err != nil {
		log.Fatalf("Unable to create presentation. %v", err)
	}
	fmt.Printf("Created presentation with ID: %s", presentation.PresentationId)

	requests := genSlides(&pivInterval.Stories)

	body := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}
	response, err := gsSrv.Presentations.BatchUpdate(presentation.PresentationId, body).Do()
	if err != nil {
		log.Fatalf("Unable to create slide. %v", err)
	}
	fmt.Printf("Created slide with ID: %s", response.Replies[0].CreateSlide.ObjectId)
}

func (gs *GSlides) Export(pivInterval *pivotal.Interval) {
	createPres(pivInterval)
}

func (gs *GSlides) DelAuth() {
	os.Remove(apiCreds)
	os.Remove(apiToken)

	logger.Printf("Deleted authentication files\n%s\n%s", apiCreds, apiToken)
}
