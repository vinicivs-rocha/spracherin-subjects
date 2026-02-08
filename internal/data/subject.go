package data

import (
	"context"

	"github.com/vinicivs-rocha/spracherin-subjects/domain"
)

type SubjectReader interface {
	All(ctx context.Context) ([]*domain.Subject, error)
	One(ctx context.Context, id domain.SubjectID) (*domain.Subject, error)
}

type SubjectWriter interface {
	Save(ctx context.Context, subject *domain.Subject) error
	Remove(ctx context.Context, subject *domain.Subject) error
}

type SubjectRepository interface {
	SubjectReader
	SubjectWriter
}
