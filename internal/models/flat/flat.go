package flat

type Model struct {
	FlatID   uint32 // may not be unique
	HouseID  uint32 // is unique
	Price    uint32
	RoomsNum uint32
	Status   Status
}

type Status int32

// Defines values for Status.
const (
	Undefined Status = iota
	Approved
	Created
	Declined
	OnModeration

	ApprovedString     = "approved"
	CreatedString      = "created"
	DeclinedString     = "declined"
	OnModerationString = "on moderation"
)

// Enum value maps for FlatStatus.
var (
	StatusName = map[Status]string{
		Approved:     ApprovedString,
		Created:      CreatedString,
		Declined:     DeclinedString,
		OnModeration: OnModerationString,
	}

	StatusValue = map[string]Status{
		ApprovedString:     Approved,
		CreatedString:      Created,
		DeclinedString:     Declined,
		OnModerationString: OnModeration,
	}
)

func (status Status) String() string {
	s, ok := StatusName[status]
	if !ok {
		return StatusName[Undefined]
	}

	return s
}
func GetStatusFromString(statusString string) Status {
	val, ok := StatusValue[statusString]
	if !ok {
		return Undefined
	}

	return val
}

type NewParams struct {
	HouseID  uint32
	Price    uint32
	RoomsNum uint32
}

func New(params NewParams) Model {

	flat := Model{
		HouseID:  params.HouseID,
		Price:    params.Price,
		RoomsNum: params.RoomsNum,
		Status:   Created,
	}

	return flat
}
