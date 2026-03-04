package valueobjects

import (
	"strings"

	domerrs "github.com/JorgeGorrito/PT-News-API/internal/domain/errors"
)

type ArticleStatus string

const (
	Draft     ArticleStatus = "BORRADOR"
	Published ArticleStatus = "PUBLICADO"
)

func (s ArticleStatus) IsValid() bool {
	switch s {
	case Draft, Published:
		return true
	default:
		return false
	}
}

func (s ArticleStatus) String() string {
	return string(s)
}

func (s ArticleStatus) Equals(other ArticleStatus) bool {
	return s == other
}

func NewArticleStatus(status string) (ArticleStatus, error) {
	s := ArticleStatus(strings.ToUpper(strings.TrimSpace(status)))
	if !s.IsValid() {
		return "", domerrs.NewDomainError("estado de artículo inválido: "+status, domerrs.InvalidArticleStatusError)
	}
	return s, nil
}
