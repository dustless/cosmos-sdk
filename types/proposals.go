package types

const (
	MaxDescriptionLength int = 5000
	MaxTitleLength       int = 140
)

// ProposalContent is an interface that has title, description, and proposaltype for the Proposal
type ProposalContent interface {
	GetTitle() string
	GetDescription() string
	ProposalRoute() string
	ProposalType() string
}

// Text Proposals
type ProposalAbstract struct {
	Title       string `json:"title"`       //  Title of the proposal
	Description string `json:"description"` //  Description of the proposal
}

func NewProposalAbstract(title, description string) ProposalAbstract {
	return ProposalAbstract{
		Title:       title,
		Description: description,
	}
}

func (abs ProposalAbstract) ValidateBasic() Error {
	// XXX
	if len(abs.Title) == 0 {

	}
	if len(abs.Title) > MaxTitleLength {

	}
	if len(abs.Description) == 0 {

	}
	if len(abs.Description) > MaxDescriptionLength {

	}
	return nil
}

// nolint
func (abs ProposalAbstract) GetTitle() string       { return abs.Title }
func (abs ProposalAbstract) GetDescription() string { return abs.Description }
