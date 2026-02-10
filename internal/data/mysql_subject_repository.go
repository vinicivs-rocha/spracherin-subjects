package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/domain"
)

type MySQLSubjectRepository struct {
	db *sql.DB
}

func NewMySQLSubjectRepository(db *sql.DB) *MySQLSubjectRepository {
	return &MySQLSubjectRepository{db: db}
}

type subjectConceptDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

const (
	subjectSelectAllQuery = `
SELECT id, title, description_text, description_links, concepts
FROM subjects`

	subjectSelectOneQuery = `
SELECT id, title, description_text, description_links, concepts
FROM subjects
WHERE id = ?`

	subjectUpsertQuery = `
INSERT INTO subjects (id, title, description_text, description_links, concepts)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	title = VALUES(title),
	description_text = VALUES(description_text),
	description_links = VALUES(description_links),
	concepts = VALUES(concepts)`

	subjectDeleteQuery = `
DELETE FROM subjects
WHERE id = ?`
)

func (r *MySQLSubjectRepository) All(ctx context.Context) ([]*domain.Subject, error) {
	rows, err := r.db.QueryContext(ctx, subjectSelectAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := make([]*domain.Subject, 0)

	for rows.Next() {
		var id string
		var title string
		var descriptionText string
		var descriptionLinksJSON []byte
		var conceptsJSON []byte

		err = rows.Scan(&id, &title, &descriptionText, &descriptionLinksJSON, &conceptsJSON)
		if err != nil {
			return nil, err
		}

		subject, err := subjectFromValues(id, title, descriptionText, descriptionLinksJSON, conceptsJSON)
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *MySQLSubjectRepository) One(ctx context.Context, id domain.SubjectID) (*domain.Subject, error) {
	if id.IsZero() {
		return nil, errors.New("id should not be empty")
	}

	var idValue string
	var title string
	var descriptionText string
	var descriptionLinksJSON []byte
	var conceptsJSON []byte

	err := r.db.QueryRowContext(ctx, subjectSelectOneQuery, id.String()).
		Scan(&idValue, &title, &descriptionText, &descriptionLinksJSON, &conceptsJSON)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return subjectFromValues(idValue, title, descriptionText, descriptionLinksJSON, conceptsJSON)
}

func (r *MySQLSubjectRepository) Save(ctx context.Context, subject *domain.Subject) error {
	if subject == nil {
		return errors.New("subject should not be nil")
	}

	id := subject.GetID()
	if id.IsZero() {
		return errors.New("id should not be empty")
	}

	title := subject.GetTitle()
	if title.IsZero() {
		return errors.New("title should not be empty")
	}

	description := subject.GetDescription()
	if description.IsZero() {
		return errors.New("description should not be empty")
	}

	descriptionLinks := description.Links()
	if descriptionLinks == nil {
		descriptionLinks = map[string]string{}
	}

	descriptionLinksJSON, err := json.Marshal(descriptionLinks)
	if err != nil {
		return err
	}

	concepts := subject.GetConcepts()
	conceptDTOs := make([]subjectConceptDTO, 0, len(concepts))
	for _, concept := range concepts {
		conceptDTOs = append(conceptDTOs, subjectConceptDTO{
			Name:        concept.Name(),
			Description: concept.Description(),
		})
	}

	conceptsJSON, err := json.Marshal(conceptDTOs)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(
		ctx,
		subjectUpsertQuery,
		id.String(),
		title.String(),
		description.Text(),
		descriptionLinksJSON,
		conceptsJSON,
	)
	return err
}

func (r *MySQLSubjectRepository) Remove(ctx context.Context, subject *domain.Subject) error {
	if subject == nil {
		return errors.New("subject should not be nil")
	}

	id := subject.GetID()
	if id.IsZero() {
		return errors.New("id should not be empty")
	}

	res, err := r.db.ExecContext(ctx, subjectDeleteQuery, id.String())
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrSubjectNotFound
	}

	return nil
}

func subjectFromValues(
	id string,
	title string,
	descriptionText string,
	descriptionLinksJSON []byte,
	conceptsJSON []byte,
) (*domain.Subject, error) {
	subjectID, err := domain.SubjectIDFromString(id)
	if err != nil {
		return nil, err
	}

	subjectTitle := domain.SubjectTitle(title)
	if subjectTitle.IsZero() {
		return nil, errors.New("title should not be empty")
	}

	descriptionLinks := make(map[string]string)
	if len(descriptionLinksJSON) > 0 {
		if err := json.Unmarshal(descriptionLinksJSON, &descriptionLinks); err != nil {
			return nil, err
		}
	}

	description, err := domain.NewSubjectDescription(descriptionText, descriptionLinks)
	if err != nil {
		return nil, err
	}

	conceptDTOs := make([]subjectConceptDTO, 0)
	if len(conceptsJSON) > 0 {
		if err := json.Unmarshal(conceptsJSON, &conceptDTOs); err != nil {
			return nil, err
		}
	}

	concepts := make([]domain.Concept, 0, len(conceptDTOs))
	for _, dto := range conceptDTOs {
		concept, err := domain.NewConcept(dto.Name, dto.Description)
		if err != nil {
			return nil, err
		}
		concepts = append(concepts, concept)
	}

	return domain.NewSubject(subjectID, subjectTitle, description, concepts), nil
}

var _ SubjectRepository = (*MySQLSubjectRepository)(nil)
