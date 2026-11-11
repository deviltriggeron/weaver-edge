[English](architecture.md) | [Русский](architecture.ru.md)

# Weaver Architecture

## 1. System Purpose

Weaver is a system for managing agreements, commitments, and work-related events emerging from organizational communications.

It connects external workplace systems, organizational data, and intelligent information processing.

The Weaver architecture is built around three distinct areas of responsibility:

- **Weaver Edge** — access to workplace data and integration with external systems.
- **Weaver Core** — business logic, event processing, and workflow orchestration.
- **AI Service** — data analysis and interpretation.

The main architectural principle is:

> **Edge is responsible for data access. Core is responsible for orchestrating data processing. AI is responsible for intelligent analysis.**

These components must not take over one another's responsibilities.

## 2. Overall Architecture

```text
External Systems
       |
       v
  Weaver Edge
       |
       | Authorized Data / Events
       v
   Weaver Core
       |
       v
    AI Service
```

Components communicate through explicitly defined interfaces and contracts.

### Weaver Edge

Operates at the boundary between external workplace systems and Weaver's internal components.

It is responsible for connecting to data sources, receiving and storing original messages, enforcing access policies, and providing authorized data.

### Weaver Core

Orchestrates work-event processing, interacts with AI, handles analysis results, and coordinates subsequent system actions.

### AI Service

Performs intelligent analysis, including message classification, event and commitment extraction, and determining whether additional context is required.

AI has no direct access to Edge or its database.

## 3. Weaver Edge

Weaver Edge is the integration and data-access layer of the system.

It operates alongside an organization's workplace systems and provides Weaver with controlled access to their data.

Edge responsibilities include:

- connecting to external platforms;
- receiving incoming messages and events;
- normalizing platform-specific data;
- storing original messages;
- providing access to stored data;
- enforcing access policies through `PolicyEngine`;
- auditing access decisions;
- communicating with Weaver Core through a defined API contract.

Edge answers the question:

> "Can this data be accessed within the scope of the current request?"

### Responsibility Boundaries

Edge must not perform semantic analysis of workplace messages or determine the meaning of agreements.

It does not decide whether a message represents a promise, proposal, assumption, or confirmed agreement.

These responsibilities belong to Core and AI.

### Source of Truth

Original messages received from external systems are the source data.

AI analysis results must not replace or modify original messages.

The system must preserve the ability to trace every processed event back to its source data.

## 4. Weaver Core

Weaver Core is the business logic and orchestration layer of Weaver.

It manages the processing lifecycle of work-related events and coordinates interactions between Edge and AI Service.

Core responsibilities include:

- receiving work events from Edge;
- orchestrating event processing;
- interacting with AI Service;
- providing data for analysis;
- processing structured AI results;
- determining whether additional context is required;
- creating requests to Edge;
- managing multi-step workflows;
- coordinating subsequent processing.

Core answers the question:

> "What do these data mean, and what should the system do next?"

Core does not own original workplace data and must not bypass Edge to retrieve data protected by access policies.

Core's internal implementation is a separate part of the system. Its interaction with Edge is defined by an external contract, not by Core's internal domain models.

## 5. AI Service

AI Service is the intelligent analysis component within the Weaver architecture.

It receives data for analysis and returns structured results.

AI may:

- classify workplace messages;
- extract events and commitments;
- identify participants and context;
- distinguish proposals, assumptions, and confirmed agreements;
- determine whether additional context is required;
- generate clarification requests.

AI must not:

- access Edge directly;
- connect to Edge's database;
- independently determine data-access permissions;
- perform actions on behalf of users without Core orchestration;
- act as the source of truth for original messages.

### Requesting Additional Context

When AI detects that information is insufficient, it returns a structured result to Core indicating that additional context is required.

Core then orchestrates data retrieval through Edge.

The responsibilities are therefore separated:

- AI determines what information is missing.
- Core creates and sends the request.
- Edge checks access permissions and returns authorized data.
- Core provides the retrieved context to AI.

AI never receives direct access to workplace resources.

## 6. Data Access and PolicyEngine

Access control is the responsibility of Weaver Edge.

All requests for protected workplace data must pass through Edge's access-control mechanisms.

The general flow is:

```text
Core
 |
 | ContextRequest
 v
Edge
 |
 v
PolicyEngine
 |
 +---- DENIED ----> Audit
 |
 +---- ALLOWED ---> Repository
                         |
                         v
                         DB
                         |
                         v
                  ContextResponse
                         |
                         v
                        Core
```

`PolicyEngine` determines whether the requested data may be provided.

Core must not independently decide whether a user is authorized to access a particular message or other workplace resource.

### Access Decision Results

Access checks must distinguish three fundamentally different outcomes:

```text
ALLOWED
    -> requested data may be returned

DENIED
    -> access is denied and a safe response is returned

SYSTEM ERROR
    -> a technical error occurred
```

Access denial is a normal protocol outcome, not a system error.

A technical error must never be interpreted as permission to access data.

## 7. Core ↔ Edge

Core and Edge communicate through an explicitly defined external API contract.

The contract defines:

- supported operations;
- request structures;
- response structures;
- identifier propagation rules;
- error handling;
- access-control outcomes;
- authentication requirements.

Internal domain types must not be used directly as transport models.

For example:

```text
Edge Domain
    |
    | Mapping
    v
Transport DTO
    |
    | HTTP API
    v
Core
```

This allows Edge and Core to evolve their internal implementations independently without coupling their internal models.

### Authentication

Communication between Edge and Core must be protected.

Separate authentication mechanisms must be used for different communication directions:

- Edge → Core: Edge authenticates with Core.
- Core → Edge: Core authenticates with Edge when requesting context.

Secrets used for these directions must not be conflated.

Communication between separate servers must use a secure connection.

## 8. Domain and Transport

Domain models describe the business domain and internal business logic of a component.

Transport DTOs describe the external communication contract.

Example structure:

```text
transport/dto
    ContextRequest
    ContextResponse
    ContextMessage

domain
    ContextQuery
    ContextResult
    Message
```

The conversion flow is:

```text
ContextRequest
      |
      v
ContextQuery
      |
      v
Application
      |
      v
ContextResult
      |
      v
ContextResponse
```

The Transport layer converts external requests into models understood by the Application layer and converts results back into transport responses.

It must not contain business logic.

In particular, `domain.Message` must not be used directly as a DTO for Core ↔ Edge communication.

## 9. Historical Context Retrieval

One of the primary interaction scenarios is retrieving additional context for message analysis.

The target flow is:

```text
AI Service
    |
    | NEED_CONTEXT
    v
Weaver Core
    |
    | ContextRequest
    v
Weaver Edge
    |
    v
PolicyEngine
    |
    +---- DENIED ----------------> Audit
    |
    +---- ALLOWED
             |
             v
         Repository
             |
             v
             DB
             |
             v
       ContextResponse
             |
             v
        Weaver Core
             |
             v
         AI Service
```

AI does not request arbitrary resources directly.

Context requests must be limited to supported operations and data scopes.

For example:

```text
kind:
    messages

scope:
    current_thread
    current_channel
```

The list of supported operations and scopes is defined by the Edge API contract.

Edge independently checks whether the requested context can be accessed and applies the relevant access policies.

## 10. Data Storage and Repository

Weaver Edge owns the original workplace data received from external systems.

Edge's database is an internal infrastructure detail and is not directly accessible to Core or AI Service.

Data access is organized through Repository interfaces defined on the Edge side.

The current interfaces are:

```go
type (
    ContextRepository interface {
        SaveMessage(ctx context.Context, msg domain.Message) error
        FetchContext(ctx context.Context, req domain.ContextRequest) ([]domain.Message, error)
    }

    PolicyRepository interface {
        Load(ctx context.Context) (*domain.PolicySnapshot, error)
    }

    AuditRepository interface {
        Record(ctx context.Context, event domain.AuditEvent) error
    }
)
```

`ContextRepository` is responsible for storing original messages and retrieving historical context.

`PolicyRepository` provides the Application layer with the current access policy.

`AuditRepository` is responsible for persisting audit events.

The Application layer depends on interfaces rather than concrete storage implementations.

The current Edge implementation uses PostgreSQL.

This allows storage infrastructure to evolve without moving business logic into concrete Repository implementations.

## 11. Audit and Data Provenance

Auditing is necessary to explain system behavior and reconstruct the sequence of actions.

The architecture distinguishes two areas of responsibility.

### Edge Audit

Edge records decisions related to access to workplace data.

The minimum information includes:

- event timestamp;
- `request_id`;
- resource;
- operation;
- policy decision;
- actor;
- connector;
- denial code, if access is denied.

Internal policy details must not be shared with Core or AI unless necessary.

### System Audit

System auditing must make it possible to reconstruct the actual sequence of system events and actions.

Related events must preserve causal relationships, for example through `ParentEventID`.

Audit records must describe actual actions and outcomes, not merely AI intentions or assumptions.

### Provenance

Analysis results must retain their connection to the original messages.

This is necessary for verification, explanation, and subsequent correction of AI results.

AI is not the source of truth. The original message from the external system remains the source of the underlying information.

## 12. External Integrations

Weaver Edge provides connections to external workplace platforms through dedicated integration components.

Telegram is the first integration.

The architecture may later be extended with adapters for:

- Slack;
- Discord;
- Microsoft Teams;
- other enterprise communication platforms.

These integrations represent possible development directions, not a claim that they are currently implemented.

Integration components must convert platform-specific events into a consistent internal Edge representation.

This prevents Core from becoming coupled to the specifics of individual external platforms.

## 13. Target Vertical Slice

The first target end-to-end Weaver scenario is:

```text
Telegram
    |
    v
Weaver Edge
    |
    v
Receive and Normalize
    |
    v
Store Original Message
    |
    v
Weaver Core
    |
    v
AI Service
    |
    | NEED_CONTEXT
    v
Weaver Core
    |
    | ContextRequest
    v
Weaver Edge
    |
    v
PolicyEngine
    |
    +---- DENIED
    |
    +---- ALLOWED
             |
             v
         Repository
             |
             v
             DB
             |
             v
       ContextResponse
             |
             v
        Weaver Core
             |
             v
         AI Service
```

This is a target scenario, not a statement that the entire flow has already been implemented.

The goal of the first vertical slice is to establish a working permission-aware cycle:

```text
receive a workplace message
          |
          v
store the original message
          |
          v
send the event to Core
          |
          v
analyze through AI
          |
          v
determine whether context is required
          |
          v
check access through Edge
          |
          v
return authorized data
          |
          v
continue analysis
```

## 14. Reliability and Operations

The architecture must account not only for the normal processing flow but also for failures of individual components.

Key requirements include:

- idempotent processing of incoming events;
- controlled retries;
- timeouts for inter-service requests;
- graceful handling of unavailable dependencies;
- graceful process shutdown;
- proper cleanup of connections and resources;
- preservation of original data integrity;
- no bypassing of access checks during failures.

Infrastructure errors, business-logic errors, and access denials must be handled as distinct outcomes.

## 15. System Evolution Principles

New architectural abstractions should be introduced in response to concrete use cases.

The system should not prematurely introduce:

- a universal retrieval engine;
- a complex Work Graph;
- a message bus;
- a distributed cache;
- separate microservices without a clear need;
- a complex policy control plane.

First, implement a concrete vertical scenario. Then extract reusable components from solutions that prove their value.

The primary goal of the current development stage is:

> **Establish a reliable permission-aware path from an original workplace message to its analysis and understanding by Core/AI, while preserving access control, data provenance, and auditability.**

---

*Document status: Architecture description for the Weaver Alpha project. Some scenarios describe target behavior and may not yet be implemented.*