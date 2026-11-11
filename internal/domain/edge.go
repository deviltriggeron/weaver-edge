package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	ID            uuid.UUID
	TraceID       uuid.UUID
	ParentEventID *uuid.UUID
	Action        AuditAction
	Result        AuditResult
	MessageID     *uuid.UUID
	Decision      *Decision
	Error         string
	CreatedAt     time.Time
}

type AuditAction string

const (
	AuditActionMessageReceived AuditAction = "message_received"
	AuditActionPolicyCheck     AuditAction = "policy_check"
	AuditActionCoreRequest     AuditAction = "core_request"
	AuditActionContextRequest  AuditAction = "context_request"
)

type AuditResult string

const (
	AuditResultAllowed AuditResult = "allowed"
	AuditResultDenied  AuditResult = "denied"
	AuditResultSuccess AuditResult = "success"
	AuditResultFailure AuditResult = "failure"
)

type Platform string

const (
	Telegram Platform = "telegram"
)

type PendingRequestStatus string

const (
	Pending   PendingRequestStatus = "pending"
	Resolved  PendingRequestStatus = "resolved"
	Cancelled PendingRequestStatus = "cancelled"
)

type (
	Resource struct {
		ChannelID string
		ThreadID  string
	}

	Policy struct {
		ID        uuid.UUID
		Resource  Resource
		CreatedAt time.Time
	}

	PolicySnapshot struct {
		Allowed []Policy
	}

	PendingRequest struct {
		ID                uuid.UUID
		TraceID           uuid.UUID
		OriginalMessageID uuid.UUID
		ResponseMessageID *uuid.UUID
		Reason            ClarificationReason
		Question          string
		Status            PendingRequestStatus
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}
)
