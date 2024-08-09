package user

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/Ropho/avito-bootcamp-assignment/internal/models/hash"
)

type Model struct {
	ID    uuid.UUID
	Email string

	EncryptedPassword string
	Salt              string
	Type              Type
}

type Type int32

// Defines values for Status.
const (
	Undefined Type = iota
	Client
	Moderator

	ClientString    = "client"
	ModeratorString = "moderator"
)

// Enum value maps for UserType.
var (
	TypeName = map[Type]string{
		Client:    ClientString,
		Moderator: ModeratorString,
	}

	TypeValue = map[string]Type{
		ClientString:    Client,
		ModeratorString: Moderator,
	}
)

func (userType Type) String() string {
	s, ok := TypeName[userType]
	if !ok {
		return TypeName[Undefined]
	}

	return s
}
func GetTypeFromString(typeString string) Type {
	val, ok := TypeValue[typeString]
	if !ok {
		return Undefined
	}

	return val
}

type NewParams struct {
	Email    string
	Password string
	UserType string
}

func New(params NewParams) (Model, error) {

	salt, err := hash.GenerateRandomString(hash.SaltLen)
	if err != nil {
		return Model{}, fmt.Errorf("unable to generate salt: [%w]", err)
	}
	encrPass := hash.ComputeHashWithSalt(params.Password, salt)

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return Model{}, fmt.Errorf("unable to generate user uuuid: [%w]", err)
	}

	user := Model{
		ID:                userUUID,
		Email:             params.Email,
		EncryptedPassword: encrPass,
		Salt:              salt,
		Type:              GetTypeFromString(params.UserType),
	}

	return user, nil
}
