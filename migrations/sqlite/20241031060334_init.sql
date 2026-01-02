-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects (    
    name TEXT PRIMARY KEY NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,   
    html_content TEXT NOT NULL,
    note_date DATETIME NOT NULL UNIQUE
);

CREATE TABLE project_excerpts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_name TEXT NOT NULL REFERENCES projects(name) ON DELETE CASCADE ON UPDATE CASCADE,
    note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    excerpt TEXT NOT NULL,
    note_date DATETIME NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE project_excerpts;
DROP TABLE projects;
DROP TABLE notes;

-- +goose StatementEnd
