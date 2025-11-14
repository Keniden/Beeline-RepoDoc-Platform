CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE repos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT,
    url TEXT UNIQUE,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id) ON DELETE CASCADE,
    path TEXT,
    module TEXT,
    cohesion DOUBLE PRECISION,
    coupling DOUBLE PRECISION
);

CREATE TABLE functions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    file_id UUID REFERENCES files(id),
    name TEXT,
    signature TEXT
);

CREATE TABLE structs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    file_id UUID REFERENCES files(id),
    name TEXT
);

CREATE TABLE interfaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    file_id UUID REFERENCES files(id),
    name TEXT
);

CREATE TABLE edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    from_id UUID,
    to_id UUID,
    type TEXT
);

CREATE TABLE commits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    author TEXT,
    message TEXT,
    committed_at TIMESTAMPTZ
);

CREATE TABLE commits_files (
    commit_id UUID REFERENCES commits(id),
    repo_id UUID REFERENCES repos(id),
    file_id UUID REFERENCES files(id)
);

CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    name TEXT
);

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id UUID REFERENCES repos(id),
    name TEXT,
    path TEXT
);

CREATE MATERIALIZED VIEW module_activity AS
SELECT repo_id, module, count(*) AS file_count FROM files GROUP BY repo_id, module;

CREATE MATERIALIZED VIEW hot_files AS
SELECT repo_id, path, cohesion, coupling FROM files ORDER BY coupling DESC LIMIT 10;
