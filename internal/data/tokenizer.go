package data

import "github.com/vinicivs-rocha/spracherin-subjects/internal/domain"

type Tokenizer interface {
	GetVocabulary() (domain.Vocabulary, error)
}
