# 0001. Use Zustand for State Management

Date: 2026-01-11

## Status

Accepted

## Context

We needed a global state management solution for our Next.js application to handle:
- User session data (Profile, Roles)
- UI state (Sidebar, Modals)
- Organization selection

The solution needed to be compatible with Next.js App Router (Server/Client components boundaries) and TypeScript.

## Decision

We chose **Zustand**.

## Consequences

**Positive:**
- **Simplicity**: Minimal boilerplate compared to Redux.
- **Performance**: Selectors allow components to subscribe to specific slices of state, preventing unnecessary re-renders.
- **DX**: Simple hook-based API (`useStore(selector)`) feels natural in React.
- **Bundle Size**: Very small (<2kb) compared to other libraries.
- **Middleware**: Built-in support for persisting state to `localStorage` (`persist` middleware).

**Negative:**
- **Structure**: Being unopinionated, it requires discipline to organize stores logically (e.g., separate files for separate domains like `user.ts`, `org.ts`).
