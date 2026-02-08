package data

import (
	"context"

	"github.com/vinicivs-rocha/spracherin-subjects/internal/domain"
)

type SubjectChanges struct {
	subjectID domain.SubjectID
	diff      *domain.ConceptsDiff
}

func NewSubjectChanges(subjectID domain.SubjectID, diff *domain.ConceptsDiff) SubjectChanges {
	return SubjectChanges{
		subjectID: subjectID,
		diff:      diff,
	}
}

type Messager interface {
	MessageDetectedChanges(ctx context.Context, changes SubjectChanges) error
}
