package story

type State string

const (
	ACCEPTED = State("accepted")
	DELIVERED = State("delivered")
	FINISHED = State("finished")
	STARTED = State("started")
	REJECTED = State("rejected")
	PLANNED = State("planned")
	UNSTARTED = State("unstarted")
	UNSCHEDULED = State("unscheduled")
)
