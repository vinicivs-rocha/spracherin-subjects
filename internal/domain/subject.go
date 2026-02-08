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

func SubjectIDFromString(id string) (SubjectID, error) {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return SubjectID{}, err
	}

	return SubjectID{uuid: uuid}, nil
}

func (si SubjectID) IsZero() bool {
	return si.uuid == uuid.Nil
}

type SubjectTitle string

func (st SubjectTitle) IsZero() bool {
	return st == ""
}

type SubjectDescription struct {
	text  string
	links map[string]string
}

func NewSubjectDescription(text string, links map[string]string) (SubjectDescription, error) {
	if text == "" {
		return SubjectDescription{}, errors.New("value should not be empty")
	}

	newLinks := make(map[string]string, len(links))

	maps.Copy(newLinks, links)

	return SubjectDescription{text: text, links: newLinks}, nil
}

func (sd SubjectDescription) IsZero() bool {
	return sd.text == ""
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

func (c Concept) IsZero() bool {
	return c.description == "" || c.name == ""
}

type Subject struct {
	id          SubjectID
	title       SubjectTitle
	description SubjectDescription
	concepts    []Concept
}

func NewSubject(id SubjectID, subjectTitle SubjectTitle, subjectDescription SubjectDescription, concepts []Concept) *Subject {
	newConcepts := make([]Concept, len(concepts))
	copy(newConcepts, concepts)

	return &Subject{
		id:          id,
		title:       subjectTitle,
		description: subjectDescription,
		concepts:    newConcepts,
	}
}

func (s *Subject) ExtractConcepts() error {
	// TODO: implement text concepts extraction
	return nil
}

type ConceptsDiff struct {
	Updated []Concept
	Removed []Concept
	Added   []Concept
}

func (s *Subject) EvaluateConceptsDiff(next *Subject) (*ConceptsDiff, error) {
	// TODO: implement diff
	return nil, nil
}

func (s *Subject) GetID() SubjectID {
	return s.id
}
