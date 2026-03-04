package entities

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	constants "github.com/JorgeGorrito/PT-News-API/internal/domain/constants"
	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
)

type Author struct {
	BaseEntity[int64]
	name      string
	email     string
	biography string
	createdAt time.Time
}

func NewAuthor(name, email string) (*Author, error) {
	author := &Author{
		BaseEntity: BaseEntity[int64]{id: 0},
		biography:  "",
		createdAt:  time.Now().UTC(),
	}

	if err := author.SetName(name); err != nil {
		return nil, err
	}

	if err := author.SetEmail(email); err != nil {
		return nil, err
	}

	return author, nil
}

func (a *Author) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return domerrs.NewDomainError("El nombre del autor no puede estar vacío", domerrs.EmptyAuthorNameError)
	}

	if len(name) < constants.MinAuthorNameLength {
		return domerrs.NewDomainError(fmt.Sprintf("El nombre del autor debe tener al menos %d caracteres", constants.MinAuthorNameLength), domerrs.InvalidAuthorNameError)
	}

	if len(name) > constants.MaxAuthorNameLength {
		return domerrs.NewDomainError(fmt.Sprintf("El nombre del autor no puede exceder %d caracteres", constants.MaxAuthorNameLength), domerrs.InvalidAuthorNameError)
	}

	a.name = name
	return nil
}

func (a *Author) SetEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return domerrs.NewDomainError("El email del autor no puede estar vacío", domerrs.EmptyAuthorEmailError)
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return domerrs.NewDomainError("El formato del email no es válido", domerrs.InvalidAuthorEmailError)
	}

	a.email = email
	return nil
}

func (a *Author) SetBiography(biography string) error {
	biography = strings.TrimSpace(biography)

	if len(biography) > constants.MaxAuthorBiographyLength {
		return domerrs.NewDomainError(fmt.Sprintf("La biografía no puede exceder %d caracteres", constants.MaxAuthorBiographyLength), domerrs.InvalidAuthorBiographyError)
	}

	a.biography = biography
	return nil
}

func (a *Author) Name() string {
	return a.name
}

func (a *Author) Email() string {
	return a.email
}

func (a *Author) Biography() string {
	return a.biography
}

func (a *Author) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Author) UpdateProfile(name, email, biography string) error {
	if err := a.SetName(name); err != nil {
		return err
	}

	if err := a.SetEmail(email); err != nil {
		return err
	}

	if biography != "" {
		if err := a.SetBiography(biography); err != nil {
			return err
		}
	}

	return nil
}

func (a *Author) String() string {
	return fmt.Sprintf("Author{ID: %d, Name: %s, Email: %s}", a.ID(), a.name, a.email)
}
