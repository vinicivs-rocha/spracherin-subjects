package application

import (
	"context"
	"errors"

	"github.com/vinicivs-rocha/spracherin-subjects/data"
	"github.com/vinicivs-rocha/spracherin-subjects/domain"
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
}

func createSubject(command DescribeSubjectCommand) (*domain.Subject, error) {
	var id domain.SubjectID
	var err error

	if command.ID == "" {
		id = domain.NewSubjectID()
	} else {
		id, err = domain.SubjectIDFromString(command.ID)
	}

	if err != nil {
		return &domain.Subject{}, err
	}

	title := domain.SubjectTitle(command.Title)

	if title.IsZero() {
		return &domain.Subject{}, errors.New("title should not be empty")
	}

	description, err := domain.NewSubjectDescription(command.Description.Text, command.Description.Links)

	if err != nil {
		return &domain.Subject{}, err
	}

	concepts := make([]domain.Concept, 0)

	newSubject := domain.NewSubject(id, title, description, concepts)

	err = newSubject.ExtractConcepts()

	if err != nil {
		return &domain.Subject{}, err
	}

	return newSubject, nil
}

func (ds *DescribeSubject) detectChanges(current *domain.Subject, next *domain.Subject) error {
	diff, err := current.EvaluateConceptsDiff(next)

	if err != nil {
		return err
	}

	changes := data.NewSubjectChanges(next.GetID(), diff)
	return ds.messager.MessageDetectedChanges(changes)
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
		return domain.SubjectID{}, errors.New("title should not be empty")
	}

	description, err := domain.NewSubjectDescription(command.Description.Text, command.Description.Links)

	if err != nil {
		return domain.SubjectID{}, err
	}

	concepts := make([]domain.Concept, 0)

	new := domain.NewSubject(id, title, description, concepts)

	err = new.ExtractConcepts()

	if err != nil {
		return domain.SubjectID{}, err
	}

	current, err := ds.repository.One(ctx, new.GetID())

	if err != nil {
		return domain.SubjectID{}, err
	}

	if current != nil {
		ds.detectChanges(current, new)
	}

	err = ds.repository.Save(ctx, new)

	if err != nil {
		return domain.SubjectID{}, err
	}

	return new.GetID(), nil
}
