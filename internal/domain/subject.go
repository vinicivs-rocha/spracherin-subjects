package domain

import (
	"errors"
	"maps"

	"github.com/google/uuid"
)

var (
	ErrEmptyTitle       = errors.New("title should not be empty")
	ErrEmptyDescription = errors.New("description should not be empty")
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

func (si SubjectID) String() string {
	return si.uuid.String()
}

type SubjectTitle string

func (st SubjectTitle) IsZero() bool {
	return st == ""
}

func (st SubjectTitle) String() string {
	return string(st)
}

type SubjectDescription struct {
	text  string
	links map[string]string
}

func NewSubjectDescription(text string, links map[string]string) (SubjectDescription, error) {
	if text == "" {
		return SubjectDescription{}, ErrEmptyDescription
	}

	newLinks := make(map[string]string, len(links))

	maps.Copy(newLinks, links)

	return SubjectDescription{text: text, links: newLinks}, nil
}

func (sd SubjectDescription) IsZero() bool {
	return sd.text == ""
}

func (sd SubjectDescription) Text() string {
	return sd.text
}

func (sd SubjectDescription) Links() map[string]string {
	newLinks := make(map[string]string, len(sd.links))
	maps.Copy(newLinks, sd.links)
	return newLinks
}

type Concept struct {
	name string
}

func NewConcept(name string) (Concept, error) {
	if name == "" {
		return Concept{}, errors.New("name should not be empty")
	}

	return Concept{
		name: name,
	}, nil
}

func (c Concept) IsZero() bool {
	return c.name == ""
}

func (c Concept) Name() string {
	return c.name
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

func (s *Subject) ExtractConcepts(tokens []byte) error {
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

func (s *Subject) GetTitle() SubjectTitle {
	return s.title
}

func (s *Subject) GetDescription() SubjectDescription {
	return s.description
}

func (s *Subject) GetConcepts() []Concept {
	newConcepts := make([]Concept, len(s.concepts))
	copy(newConcepts, s.concepts)
	return newConcepts
}
