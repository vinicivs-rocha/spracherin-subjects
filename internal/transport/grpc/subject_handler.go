package grpc

import (
	"context"
	"errors"
	"fmt"

	subjectsv1 "github.com/joserocha/spracherin/spracherin-proto/gen/go/github.com/vinicivs-rocha/spracherin/spracherin-proto/gen/go/subjects/v1"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/application"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/data"
	"github.com/vinicivs-rocha/spracherin-subjects/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubjectHandler struct {
	subjectsv1.UnimplementedSubjectServiceServer
	repository data.SubjectRepository
	describe   *application.DescribeSubject
}

func NewSubjectHandler(repository data.SubjectRepository, messager data.Messager) *SubjectHandler {
	return &SubjectHandler{
		repository: repository,
		describe:   application.NewDescribeSubject(repository, messager),
	}
}

func (h *SubjectHandler) ListSubjects(ctx context.Context, _ *subjectsv1.ListSubjectsRequest) (*subjectsv1.ListSubjectsResponse, error) {
	subjects, err := h.repository.All(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list subjects")
	}

	response := &subjectsv1.ListSubjectsResponse{
		Subjects: make([]*subjectsv1.SubjectSummary, 0, len(subjects)),
	}

	for _, subject := range subjects {
		if subject == nil {
			continue
		}
		response.Subjects = append(response.Subjects, &subjectsv1.SubjectSummary{
			Id:    subject.GetID().String(),
			Title: subject.GetTitle().String(),
		})
	}

	return response, nil
}

func (h *SubjectHandler) GetSubject(ctx context.Context, req *subjectsv1.GetSubjectRequest) (*subjectsv1.GetSubjectResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := domain.SubjectIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	subject, err := h.repository.One(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch subject")
	}
	if subject == nil {
		return nil, status.Error(codes.NotFound, "subject not found")
	}

	return &subjectsv1.GetSubjectResponse{
		Subject: &subjectsv1.SubjectDetail{
			Id:          subject.GetID().String(),
			Title:       subject.GetTitle().String(),
			Description: subject.GetDescription().Text(),
		},
	}, nil
}

func (h *SubjectHandler) DescribeSubject(ctx context.Context, req *subjectsv1.DescribeSubjectRequest) (*subjectsv1.DescribeSubjectResponse, error) {
	if req.GetId() != "" {
		if _, err := domain.SubjectIDFromString(req.GetId()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}
	}

	command := application.DescribeSubjectCommand{
		ID:    req.GetId(),
		Title: req.GetTitle(),
		Description: application.DescribeSubjectCommandDescription{
			Text:  req.GetDescription(),
			Links: map[string]string{},
		},
	}

	id, err := h.describe.Execute(ctx, command)
	if err != nil {
		return nil, mapDescribeError(err)
	}

	return &subjectsv1.DescribeSubjectResponse{Id: id.String()}, nil
}

func (h *SubjectHandler) RemoveSubject(ctx context.Context, req *subjectsv1.RemoveSubjectRequest) (*subjectsv1.RemoveSubjectResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	id, err := domain.SubjectIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	subject, err := h.repository.One(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch subject")
	}
	if subject == nil {
		return nil, status.Error(codes.NotFound, "subject not found")
	}

	if err := h.repository.Remove(ctx, subject); err != nil {
		if errors.Is(err, data.ErrSubjectNotFound) {
			return nil, status.Error(codes.NotFound, "subject not found")
		}
		return nil, status.Error(codes.Internal, "failed to remove subject")
	}

	return &subjectsv1.RemoveSubjectResponse{}, nil
}

func mapDescribeError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrEmptyTitle),
		errors.Is(err, domain.ErrEmptyDescription):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, fmt.Sprintf("describe failed: %s", err.Error()))
	}
}
