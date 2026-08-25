# Go ORM — Developer-Friendly Database Toolkit

> A type-safe, developer-friendly relational database toolkit for Go, designed around explicit APIs, composable abstractions, and a direct relationship with SQL.

## Philosophy

This project aims to answer a simple question:

> **What should a modern database toolkit feel like when it is designed from the ground up around Go's philosophy?**

The goal is not to hide the database behind a magical abstraction. The goal is to make common database work **simple and expressive**, while keeping behavior **explicit, inspectable, composable, and controllable**.

### Core principles

#### 1. Simple by default

The common path should require very little ceremony.

```go
users, err := db.User.
    Query(ctx).
    Where(User.Age.Gt(18)).
    OrderBy(User.CreatedAt.Desc()).
    Limit(20).
    All()

if err != nil {
    return err
}
```

A developer should be able to express intent without manually assembling SQL for ordinary operations.

#### 2. Explicit when it matters

The abstraction should never obscure important behavior.

Context, errors, transactions, connection lifetime, and database-specific behavior remain visible in the API.

```go
err := db.Transaction(ctx, func(tx *Tx) error {
    if err := tx.User.Create(ctx, user); err != nil {
        return err
    }

    return tx.AuditLog.Create(ctx, event)
})
```

There is no hidden global error state and no implicit transaction magic.

#### 3. Compile-time safety where possible

Prefer generated, typed Go APIs over strings and runtime discovery.

Instead of:

```go
Where("emial = ?", value)
```

the generated API should make valid fields discoverable:

```go
Where(User.Email.Eq(value))
```

A typo should ideally become a compiler error rather than a production bug.

#### 4. Composable abstractions

The library should be built from small pieces that can be combined.

```go
query := db.User.
    Query(ctx).
    Where(User.Active.Eq(true)).
    OrderBy(User.CreatedAt.Desc())

users, err := query.Limit(50).All()
```

The query builder, SQL compiler, executor, transaction layer, and connection pool should remain independently understandable.

#### 5. Never trap the developer

The high-level API is the default—not a prison.

When the abstraction stops being useful, developers should be able to move down the stack.

```go
rows, err := db.Raw(ctx, `
    SELECT id, email
    FROM users
    WHERE last_login > $1
`, cutoff)
```

The toolkit should make SQL easier to use, not make SQL inaccessible.

#### 6. Database knowledge should remain visible

The ORM should not encourage developers to forget how relational databases work.

Relationships, indexes, constraints, transactions, isolation, joins, and query plans should have clear representations in the system.

The abstraction exists to reduce repetitive work—not to remove understanding.

---

# Developer Experience

The intended experience is:

```text
        Developer intent
               │
               ▼
       Typed Go API
               │
               ▼
        Query builder
               │
               ▼
           SQL AST
               │
               ▼
        SQL compiler
               │
               ▼
       Database executor
               │
               ▼
          PostgreSQL
```

A typical application should feel concise:

```go
posts, err := db.Post.
    Query(ctx).
    Where(Post.Published.Eq(true)).
    Where(Post.AuthorID.Eq(userID)).
    Include(Post.Author).
    OrderBy(Post.CreatedAt.Desc()).
    Limit(20).
    All()
```

But the generated API should ultimately compile into ordinary SQL and parameters.

The developer should always have a path to inspect what is happening.

```go
sql, args := query.SQL()

fmt.Println(sql)
fmt.Println(args)
```

This makes the abstraction understandable rather than magical.

---

# Architecture

The system is intentionally divided into layers.

```text
┌──────────────────────────────────────────────────────┐
│                    Application                       │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                 Generated Go Client                  │
│                                                      │
│  User.Query()   Post.Query()   User.Email.Eq(...)   │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                    Query API                         │
│                                                      │
│  Where / Select / Include / OrderBy / Limit / ...   │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                     Query AST                        │
│                                                      │
│  Select / Insert / Update / Join / Predicate / ...  │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                   SQL Compiler                       │
│                                                      │
│          AST → SQL + bound parameters               │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                  Execution Layer                     │
│                                                      │
│     Context / Transactions / Pool / Rows / Errors   │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────┐
│                    PostgreSQL                        │
└──────────────────────────────────────────────────────┘
```

Supporting the main pipeline are two major subsystems:

```text
                    Schema
                      │
             ┌────────┴────────┐
             ▼                 ▼
      Code Generator      Migration Engine
             │                 │
             ▼                 ▼
       Typed Go API       Schema → SQL
```

---

# Schema and Code Generation

The project should favor generated code for the application-facing API.

A schema describes the database model:

```text
User
 ├── ID
 ├── Email
 ├── Name
 └── CreatedAt

Post
 ├── ID
 ├── Title
 ├── Body
 ├── AuthorID
 └── CreatedAt
```

The generator turns that information into typed Go constructs:

```go
User.ID
User.Email
User.Name
User.CreatedAt

Post.ID
Post.Title
Post.AuthorID
```

This provides:

- compiler-checked field references
- IDE autocomplete
- discoverable relationships
- fewer string-based queries
- predictable generated code
- a clear boundary between schema and application code

Generated code should remain readable. A developer should be able to open it and understand what the toolkit generated.

---

# Query System

Queries are represented as structured data rather than strings.

For example:

```go
db.User.
    Query(ctx).
    Where(
        User.Age.Gt(18),
        User.Active.Eq(true),
    ).
    OrderBy(
        User.CreatedAt.Desc(),
    ).
    Limit(20)
```

Internally this becomes something conceptually similar to:

```text
SELECT
    user.id,
    user.email,
    user.name
FROM users AS user
WHERE
    user.age > $1
    AND user.active = $2
ORDER BY user.created_at DESC
LIMIT $3
```

The important separation is:

```text
Intent
  ↓
AST
  ↓
Dialect-specific SQL
```

The query API should not need to know how PostgreSQL syntax is emitted.

---

# SQL and Database Dialects

The database layer should use explicit dialect boundaries.

```go
type Dialect interface {
    QuoteIdentifier(name string) string
    Placeholder(position int) string
    Compile(node Node) (string, []any, error)
}
```

This allows database-specific behavior to live in a well-defined layer.

The initial implementation may focus exclusively on PostgreSQL.

Supporting additional databases should happen only when the abstractions genuinely justify it.

> **Do not design for every database before understanding one database deeply.**

PostgreSQL is the reference database because it provides enough relational and database-engineering complexity to make the architecture meaningful.

---

# Transactions

Transactions are first-class.

```go
err := db.Transaction(ctx, func(tx *Tx) error {
    account, err := tx.Account.
        Query(ctx).
        Where(Account.ID.Eq(accountID)).
        First()

    if err != nil {
        return err
    }

    if err := tx.Payment.Create(ctx, payment); err != nil {
        return err
    }

    return tx.Account.UpdateBalance(ctx, account.ID, -payment.Amount)
})
```

The transaction system must understand:

- connection pinning
- commit
- rollback
- context cancellation
- nested operations
- savepoints
- isolation levels
- concurrent access

The abstraction should make transactions convenient without pretending that transaction semantics are simple.

---

# Connection Management

Database connections are shared resources and must be managed explicitly and safely.

The execution layer should eventually cover:

- connection pooling
- acquisition and release
- maximum open connections
- idle connections
- connection lifetime
- context-aware acquisition
- graceful shutdown
- health checking

Concurrency correctness is part of the design, not an afterthought.

---

# Errors

Errors should remain ordinary Go errors.

```go
user, err := db.User.
    Query(ctx).
    Where(User.ID.Eq(id)).
    First()

if err != nil {
    if errors.Is(err, ErrNotFound) {
        // handle missing user
    }

    return err
}
```

Errors should carry useful context without becoming opaque framework-specific objects.

Where appropriate, errors should preserve their underlying cause:

```go
fmt.Errorf("create user: %w", err)
```

The developer should be able to use standard Go error handling tools.

---

# Context and Cancellation

Every potentially blocking database operation should accept a context.

```go
user, err := db.User.
    Query(ctx).
    Where(User.ID.Eq(id)).
    First()
```

The context should propagate through:

```text
HTTP request
      ↓
Application
      ↓
ORM
      ↓
Connection pool
      ↓
Database driver
      ↓
PostgreSQL
```

If the request is cancelled or times out, the database operation should be able to stop and release its resources correctly.

---

# Migrations

Schema changes should be represented as versioned migrations.

```bash
goorm migrate create add_user_email
goorm migrate status
goorm migrate apply
```

Conceptually:

```text
Current schema
      +
Desired schema
      │
      ▼
 Schema diff
      │
      ▼
Migration plan
      │
      ▼
 SQL migration
```

The migration engine should eventually understand:

- tables
- columns
- indexes
- constraints
- foreign keys
- unique constraints
- enum/type changes
- destructive changes
- migration ordering
- migration history

The system should be conservative around destructive operations.

A migration that can destroy data should never silently do so.

---

# Go Design Guidelines

The implementation should follow Go's strengths instead of fighting them.

### Prefer small interfaces

Bad:

```go
type Everything interface {
    Query()
    Insert()
    Update()
    Delete()
    Migrate()
    Generate()
    ...
}
```

Prefer focused interfaces:

```go
type Executor interface {
    Exec(ctx context.Context, query string, args ...any) error
}
```

Small interfaces make components easier to test and replace.

### Prefer composition

Favor:

```text
Query
 + Predicate
 + Ordering
 + Projection
 + Pagination
```

over a giant query object containing every possible behavior.

### Prefer explicit dependencies

Avoid hidden global database state.

Prefer:

```go
func NewClient(cfg Config) (*Client, error)
```

and:

```go
client := NewClient(...)
```

over:

```go
orm.Init(...)
```

with global state.

### Prefer standard library concepts

Where the standard library already provides a good abstraction, use it.

Examples:

- `context.Context`
- `error`
- `io`
- `sync`
- `time`
- `database/sql` concepts where appropriate

The project should feel like a Go library, not a foreign framework implemented in Go.

---

# Database Design Guidelines

The ORM should encourage sound relational design rather than hiding it.

### Relationships must be explicit

A relation should correspond to actual database structures:

```text
User 1 ─────── N Post
```

implemented through:

```text
posts.author_id → users.id
```

### Indexes are intentional

Indexes should be represented as schema decisions, not magical performance switches.

```text
users
 ├── PRIMARY KEY (id)
 ├── UNIQUE (email)
 └── INDEX (created_at)
```

### Constraints matter

The ORM should expose:

- primary keys
- foreign keys
- unique constraints
- check constraints
- not-null constraints

Application validation and database integrity are complementary, not interchangeable.

---

# Performance Philosophy

Performance should be measurable.

The project should include benchmarks for:

```text
Query construction
SQL compilation
Reflection / metadata lookup
Row scanning
Connection acquisition
Bulk operations
Concurrent queries
```

Avoid optimizing based on assumptions.

The development process should be:

```text
Measure
  ↓
Profile
  ↓
Understand
  ↓
Optimize
  ↓
Benchmark again
```

The project should use Go's race detector, benchmarks, fuzzing, and profiling tools as first-class development tools.

---

# Testing Strategy

Testing should happen at multiple levels.

```text
Unit tests
    ↓
Component tests
    ↓
SQL compiler tests
    ↓
Integration tests
    ↓
Concurrency tests
    ↓
Fuzz tests
    ↓
Benchmarks
```

### Unit tests

Test individual components without PostgreSQL.

### SQL compiler tests

Given an AST:

```text
Input AST
    ↓
Expected SQL
    ↓
Expected parameters
```

### Integration tests

Run against a real PostgreSQL instance.

Do not rely exclusively on mocks for database behavior.

### Concurrency tests

Test:

- concurrent queries
- pool exhaustion
- cancellation
- transaction contention
- connection reuse

### Fuzz tests

Fuzz:

- query expressions
- identifiers
- schema definitions
- SQL generation
- migration diffs

---

# Observability

A production database toolkit should make its behavior observable.

Potential hooks:

```go
type QueryObserver interface {
    BeforeQuery(ctx context.Context, q QueryInfo)
    AfterQuery(ctx context.Context, q QueryInfo, err error)
}
```

This can support:

- query timing
- query logging
- tracing
- metrics
- slow-query detection

Observability should be opt-in and should never require developers to modify their database code.

---

# Project Structure

The exact structure may evolve, but the intended separation is approximately:

```text
/
├── cmd/
│   └── goorm/
│
├── internal/
│   ├── ast/
│   ├── compiler/
│   ├── dialect/
│   ├── executor/
│   ├── generator/
│   ├── migration/
│   ├── schema/
│   └── pool/
│
├── runtime/
│   ├── query/
│   ├── transaction/
│   └── errors/
│
├── generated/
│
├── examples/
│
├── integration/
│
├── benchmarks/
│
├── migrations/
│
├── docs/
│
└── go.mod
```

The package boundaries should follow responsibilities rather than implementation convenience.

---

# Development Roadmap

The project should grow incrementally.

## Phase 1 — Foundation

- [ ] PostgreSQL connection
- [ ] Basic configuration
- [ ] Context-aware execution
- [ ] Basic CRUD
- [ ] Error model
- [ ] Initial tests

## Phase 2 — Query Engine

- [ ] Query builder
- [ ] SQL AST
- [ ] PostgreSQL compiler
- [ ] Predicates
- [ ] Ordering
- [ ] Pagination
- [ ] Projections
- [ ] SQL inspection

## Phase 3 — Relations

- [ ] One-to-one
- [ ] One-to-many
- [ ] Many-to-many
- [ ] JOIN generation
- [ ] Relation loading
- [ ] N+1 detection/mitigation

## Phase 4 — Type-Safe API

- [ ] Schema representation
- [ ] Go code generation
- [ ] Typed fields
- [ ] Typed predicates
- [ ] Typed relations
- [ ] IDE-friendly API

## Phase 5 — Transactions and Pooling

- [ ] Transactions
- [ ] Rollback
- [ ] Savepoints
- [ ] Isolation levels
- [ ] Connection pooling
- [ ] Cancellation
- [ ] Concurrency safety

## Phase 6 — Migrations

- [ ] Schema model
- [ ] Schema diffing
- [ ] Migration generation
- [ ] Migration history
- [ ] Migration execution
- [ ] Destructive-change detection

## Phase 7 — Production Engineering

- [ ] Benchmarks
- [ ] Race testing
- [ ] Fuzz testing
- [ ] Profiling
- [ ] Query observability
- [ ] Metrics/tracing hooks
- [ ] Graceful shutdown

## Phase 8 — Advanced Database Features

- [ ] CTEs
- [ ] Window functions
- [ ] Upserts
- [ ] Bulk operations
- [ ] PostgreSQL-specific types
- [ ] JSON/JSONB operations
- [ ] Array operations
- [ ] Full-text search
- [ ] Advisory locks

---

# Design Boundaries

This project intentionally does **not** try to abstract everything.

The following should remain accessible:

```text
Go
 ↓
ORM
 ↓
Query Builder
 ↓
SQL
 ↓
Database
```

A developer should be able to move between these layers depending on the problem.

### High-level

```go
users, err := db.User.
    Query(ctx).
    Where(User.Active.Eq(true)).
    All()
```

### Mid-level

```go
query := db.User.
    Query(ctx).
    Where(User.Active.Eq(true))

sql, args := query.SQL()
```

### Low-level

```go
rows, err := db.Raw(ctx, `
    SELECT id, email
    FROM users
    WHERE active = $1
`, true)
```

The layers should cooperate rather than compete.

---

# What This Project Is Trying to Prove

This project is intentionally more than an ORM.

It is an exploration of:

### Go

- language design
- interfaces
- generics
- reflection
- concurrency
- synchronization
- contexts
- error handling
- code generation
- networking
- testing
- fuzzing
- benchmarking
- profiling
- package design

### Databases

- relational modeling
- SQL
- query planning
- indexes
- constraints
- joins
- transactions
- isolation
- locking
- connection management
- schema evolution
- migration systems
- PostgreSQL internals

### Software engineering

- API design
- abstraction boundaries
- developer experience
- backward compatibility
- observability
- performance engineering
- testing strategy
- documentation

The ultimate goal is not to build the largest ORM.

It is to build one where **the abstraction is pleasant enough for everyday development while the underlying system remains understandable to the developer.**

---

# Guiding Principle

> **Make the common path simple.  
> Make the uncommon path possible.  
> Make the underlying behavior understandable.  
> Never hide what Go or the database is doing.**

That is the standard this project should be designed against.
