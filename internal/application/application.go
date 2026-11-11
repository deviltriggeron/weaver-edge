package application

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type Application struct {
	source         ports.EventSource
	policy         ports.PolicyEngine
	core           ports.CoreClient
	contextService ports.ContextService
	auditService   ports.AuditService
}

func NewApplication(
	source ports.EventSource,
	policy ports.PolicyEngine,
	core ports.CoreClient,
	contextService ports.ContextService,
	auditService ports.AuditService,
) Application {
	return Application{
		source:         source,
		policy:         policy,
		core:           core,
		contextService: contextService,
		auditService:   auditService,
	}
}

// TODO: Replace EventSource.Receive polling with webhook-driven message handling.
// Refactor Application.Handle to accept a single domain.Message.
// Ensure idempotent processing by message.ID to prevent duplicate handling
// when webhook events are delivered more than once.
// Endpoint: Events
func (a *Application) Handle(ctx context.Context, message domain.Message) error {
	err := a.contextService.SaveMessage(ctx, message)
	if err != nil {
		return domain.ErrSaveContext
	}
	traceID := uuid.New()
	resource := createResource(message)

	if !a.policy.Allow(ctx, resource) {
		event := domain.AuditEvent{
			TraceID:   traceID,
			Action:    domain.AuditActionPolicyCheck,
			Result:    domain.AuditResultDenied,
			MessageID: &message.ID,
		}

		err := a.auditService.Save(ctx, event)
		if err != nil {
			return domain.ErrSaveAudit
		}

		return nil
	}

	log.Printf("sending message to core: %+v\n", message)

	decision, err := a.core.SendMessage(ctx, message, nil)
	if err != nil {
		event := domain.AuditEvent{
			TraceID:   traceID,
			Action:    domain.AuditActionMessageReceived,
			Result:    domain.AuditResultFailure,
			MessageID: &message.ID,
		}

		auditErr := a.auditService.Save(ctx, event)
		if auditErr != nil {
			log.Printf("error save error sending %s to audit: %v\n", message.ID, err)
			return domain.ErrSaveAudit
		}

		return fmt.Errorf("send message %s to core: %w", message.ID, err)
	}

	event := domain.AuditEvent{
		TraceID:   traceID,
		Action:    domain.AuditActionMessageReceived,
		Result:    domain.AuditResultSuccess,
		MessageID: &message.ID,
		Decision:  decision,
	}

	err = a.auditService.Save(ctx, event)
	if err != nil {
		log.Printf("error save %s to audit: %v\n", message.ID, err)
		return domain.ErrSaveAudit
	}

	return nil
}

// TODO: Record policy denial in Audit.
func (a *Application) GetContext(ctx context.Context, request domain.ContextRequest) (*domain.ContextResponse, error) {
	resource := domain.Resource{
		ChannelID: request.ChannelID,
		ThreadID:  request.ThreadID,
	}

	if !a.policy.Allow(ctx, resource) {
		auditEvent := domain.AuditEvent{
			Action: domain.AuditActionContextRequest,
			Result: domain.AuditResultDenied,
		}
		err := a.auditService.Save(ctx, auditEvent)
		if err != nil {
			log.Printf("error save request to audit: %v\n", err)
			return nil, domain.ErrSaveAudit
		}

		return &domain.ContextResponse{
			Status: domain.ContextStatusDenied,
			Code:   domain.ContextDenialNotAllowed,
		}, nil
	}

	message, err := a.contextService.GetContext(ctx, request)
	if err != nil {
		// TODO:
		// 1. Distinguish between context not found and technical errors.
		// 2. Record the context request failure in Audit.
		// 3. Return ContextStatusUnavailable for technical errors.
		// 4. Return an appropriate denial code when context is unavailable.
		// 5. Ensure Core receives the context retrieval result.
		return &domain.ContextResponse{
			Status: domain.ContextStatusDenied,
			Code:   domain.ContextDenialNotFound,
		}, nil
	}

	auditEvent := domain.AuditEvent{
		Action: domain.AuditActionContextRequest,
		Result: domain.AuditResultSuccess,
	}
	err = a.auditService.Save(ctx, auditEvent)
	if err != nil {
		log.Printf("error save request to audit: %v\n", err)
		return nil, domain.ErrSaveAudit
	}

	return &domain.ContextResponse{
		Status:   domain.ContextStatusAllowed,
		Messages: message,
	}, nil
}

func createResource(message domain.Message) domain.Resource {
	return domain.Resource{
		ChannelID: message.ChannelID,
		ThreadID:  message.ThreadID,
	}
}
