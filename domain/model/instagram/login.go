package instagram

const (
	RequiredStepNone      RequiredCase = "none"
	RequiredStep2F        RequiredCase = "2f"
	RequiredStepChallenge RequiredCase = "challenge"
)

type RequiredCase string

type Login struct {
	ExternalID  string
	Credentials Credentials
	Required    Required
}

type Credentials struct {
	Username string
	Password string
	Proxy    string
}

type Required struct {
	Case    RequiredCase
	Options RequiredOptions
}

type RequiredOptions struct {
	Identifier    string
	Step          string
	CheckpointUrl string
	Method        string
	Code          string
}
