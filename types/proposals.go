package types

// ProposalContent is an interface that has title, description, and proposaltype for the Proposal
type ProposalContent interface {
	GetTitle() string
	GetDescription() string
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

// nolint
func (tp ProposalAbstract) GetTitle() string       { return tp.Title }
func (tp ProposalAbstract) GetDescription() string { return tp.Description }
