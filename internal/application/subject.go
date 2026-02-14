package application

import (
	"context"

	"github.com/vinicivs-rocha/spracherin-subjects/internal/data"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/domain"
)

type DescribeSubjectCommandDescription struct {
	Text  string            `json:"text"`
	Links map[string]string `json:"links"`
}

type DescribeSubjectCommand struct {
	ID          string                            `json:"id"`
	Title       string                            `json:"title"`
	Description DescribeSubjectCommandDescription `json:"description"`
}

type DescribeSubject struct {
	repository data.SubjectRepository
	messager   data.Messager
	tokenizer  data.Tokenizer
}

func NewDescribeSubject(repository data.SubjectRepository, messager data.Messager) *DescribeSubject {
	return &DescribeSubject{
		repository: repository,
		messager:   messager,
	}
}

func (ds *DescribeSubject) detectChanges(ctx context.Context, current *domain.Subject, next *domain.Subject) error {
	diff, err := current.EvaluateConceptsDiff(next)

	if err != nil {
		return err
	}

	changes := data.NewSubjectChanges(next.GetID(), diff)
	return ds.messager.MessageDetectedChanges(ctx, changes)
}

func (ds *DescribeSubject) Execute(ctx context.Context, command DescribeSubjectCommand) (domain.SubjectID, error) {
	var id domain.SubjectID
	var err error

	if command.ID == "" {
		id = domain.NewSubjectID()
	} else {
		id, err = domain.SubjectIDFromString(command.ID)
	}

	if err != nil {
		return domain.SubjectID{}, err
	}

	title := domain.SubjectTitle(command.Title)

	if title.IsZero() {
		return domain.SubjectID{}, domain.ErrEmptyTitle
	}

	description, err := domain.NewSubjectDescription(command.Description.Text, command.Description.Links)

	if err != nil {
		return domain.SubjectID{}, err
	}

	concepts := make([]domain.Concept, 0)

	new := domain.NewSubject(id, title, description, concepts)

	tokens, err := ds.tokenizer.Encode(new.GetDescription().Text())

	if err != nil {
		return domain.SubjectID{}, err
	}

	err = new.ExtractConcepts(tokens)

	if err != nil {
		return domain.SubjectID{}, err
	}

	current, err := ds.repository.One(ctx, new.GetID())

	if err != nil {
		return domain.SubjectID{}, err
	}

	if current != nil {
		ds.detectChanges(ctx, current, new)
	}

	err = ds.repository.Save(ctx, new)

	if err != nil {
		return domain.SubjectID{}, err
	}

	return new.GetID(), nil
}
