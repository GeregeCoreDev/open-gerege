# 0000. Use Architecture Decision Records

Date: 2026-01-11

## Status

Accepted

## Context

We need to record architectural decisions made in this project.
It is important to not only record "what" we decided, but "why" we decided it.
This allows new team members to understand the history of the project and prevents re-litigating past decisions without new context.

## Decision

We will use Architecture Decision Records (ADRs) to document significant decisions that affect the structure, non-functional characteristics, dependencies, interfaces, or construction techniques of the project.

We will follow a lightweight format similar to [MADR](https://adr.github.io/madr/):

- **Title**: Short summary of the decision
- **Status**: Proposed, Accepted, Deprecated, etc.
- **Context**: The problem we are solving and the forces at play.
- **Decision**: The solution we chose.
- **Consequences**: The positive and negative implications of the decision.

## Consequences

- **Positive**:
  - Clear history of decisions.
  - Easier onboarding for new developers.
  - Encourages thoughtful decision-making.

- **Negative**:
  - Slight overhead in documentation.
  - Requires discipline to keep up to date.
