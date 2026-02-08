package domain

import (
	"errors"
	"maps"

	"github.com/google/uuid"
)

type SubjectID struct {
	uuid uuid.UUID
}

func NewSubjectID() SubjectID {
	return SubjectID{uuid: uuid.New()}
}

func FromString(id string) (SubjectID, error) {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return SubjectID{}, err
	}

	return SubjectID{uuid: uuid}, nil
}

func (si SubjectID) IsZero() bool {
	return si.uuid == uuid.Nil
}

type Title string

func (t Title) IsZero() bool {
	return t == ""
}

type Description struct {
	text  string
	links map[string]string
}

func NewDescription(text string, links map[string]string) (Description, error) {
	if text == "" {
		return Description{}, errors.New("value should not be empty")
	}

	newLinks := make(map[string]string, len(links))

	maps.Copy(newLinks, links)

	return Description{text: text, links: newLinks}, nil
}

func (d Description) IsZero() bool {
	return d.text == ""
}

type Concept struct {
	name        string
	description string
}

func NewConcept(name string, description string) (Concept, error) {
	if name == "" {
		return Concept{}, errors.New("name should not be empty")
	}

	if description == "" {
		return Concept{}, errors.New("description should not be empty")
	}

	return Concept{
		name:        name,
		description: description,
	}, nil
}

type Subject struct {
	id          SubjectID
	title       Title
	description Description
	concepts    []Concept
}

func NewSubject(id SubjectID, title Title, description Description, concepts []Concept) *Subject {
	newConcepts := make([]Concept, len(concepts))
	copy(newConcepts, concepts)

	return &Subject{
		id:          id,
		title:       title,
		description: description,
		concepts:    newConcepts,
	}
}
