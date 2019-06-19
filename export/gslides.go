package export

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/slides/v1"
	"log"
	"os"
	"pivex/credentials"
	"pivex/pivotal"
	"strconv"
)

type GSlides struct {
	creds        *credentials.GoogleSlides
	gDriveSvc    *drive.Service
	gSlidesSvc   *slides.Service
	logger       *log.Logger
	forceCreate  bool
	pivIteration pivotal.Iteration
}

const (
	softwareTeamDriveId = "0AEiahtpqanW1Uk9PVA"
	sprintFolder        = "1ScvGIRhj780z_yWzn1ptHovjL-0V9vKK"
)

func New(creds *credentials.GoogleSlides, forceCreate bool, logger *log.Logger, pivIteration pivotal.Iteration) *GSlides {
	gs := GSlides{
		creds:        creds,
		logger:       logger,
		forceCreate:  forceCreate,
		pivIteration: pivIteration,
	}

	gs.initClients()
	gs.delExisting("tmp")

	return &gs
}

func (gs *GSlides) initClients() {
	client := gs.creds.Oauth2Config.Client(context.Background(), gs.creds.Oauth2Token)

	gDriveSvc, err := drive.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	gs.gDriveSvc = gDriveSvc

	gSlidesSvc, err := slides.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	gs.gSlidesSvc = gSlidesSvc
}

func (gs *GSlides) delExisting(slideName string) {
	teamDriveFiles, err := gs.gDriveSvc.Files.
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
					gs.gDriveSvc.Files.
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
		if story.CurrentState != "accepted" {
			continue
		}

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
					"Total\nFeatures: %d\tChores: %d\tBugs: %d\nAccepted\nFeatures: %d\tChores: %d\tBugs: %d\n",
					totalFeatures, totalChores, totalBugs, acceptedFeatures, acceptedChores, acceptedBugs),
			},
		},
	)

	return requests
}

func (gs *GSlides) createPres() {
	slideName := "Sprint Demo " + strconv.Itoa(gs.pivIteration.Number)

	gs.delExisting(slideName)

	driveFile := &drive.File{
		MimeType:    "application/vnd.google-apps.presentation",
		Name:        slideName,
		TeamDriveId: softwareTeamDriveId,
		Parents:     []string{sprintFolder},
	}

	presentation, err := gs.gDriveSvc.Files.Create(driveFile).SupportsTeamDrives(true).Do()

	if err != nil {
		gs.logger.Fatalf("Unable to create presentation. %v", err)
	}

	requests := gs.genSprintAccomplishments()

	gs.logger.Printf("Created presentation with ID: %s", presentation.Id)

	requests = append(requests, genSlides(&gs.pivIteration.Stories)...)

	body := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}
	response, err := gs.gSlidesSvc.Presentations.BatchUpdate(presentation.Id, body).Do()
	if err != nil {
		gs.logger.Fatalf("Unable to create slide. %v", err)
	}

	gs.logger.Printf("Created slide with ID: %s", response.Replies[0].CreateSlide.ObjectId)
}

func (gs *GSlides) Export() {
	gs.createPres()
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
