CREATE TABLE IF NOT EXISTS authors (
    id INTEGER PRIMARY KEY,
    name text NOT NULL,
    bio text
);
CREATE TABLE IF NOT EXISTS counters (
    "name" text PRIMARY key,
    "counter" integer NOT NULL
);
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY,
    "name" text not null,
    parent_id integer,
    FOREIGN KEY(parent_id) REFERENCES projects(id)
);
CREATE TABLE IF NOT EXISTS tags (
    tag TEXT PRIMARY KEY,
    tag_group_id INTEGER NOT NULL,
    FOREIGN KEY(tag_group_id) REFERENCES tag_groups(id)
);
CREATE TABLE IF NOT EXISTS tag_groups (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER NOT NULL,
    root_id INTEGER NOT NULL,
    color_code integer not null,
    emoji text,
    FOREIGN KEY(parent_id) REFERENCES tag_groups(id),
    FOREIGN KEY(parent_id) REFERENCES tag_groups(id)
);
CREATE TABLE IF NOT EXISTS statuss (
    id INTEGER PRIMARY KEY,
    "name" text not null,
    project_id integer not null,
    color_code integer not null,
    emoji text,
    done boolean not null,
    progressing boolean not null,
    archived boolean not null,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);
CREATE TABLE IF NOT EXISTS status_changes (
    id INTEGER PRIMARY KEY,
    "name" text not null,
    project_id integer not null,
    source_status_id integer not null,
    target_status_id integer not null,
    emoji text,
    FOREIGN KEY(source_status_id) REFERENCES statuss(id),
    FOREIGN KEY(target_status_id) REFERENCES statuss(id),
    FOREIGN KEY(project_id) REFERENCES projects(id)
);
CREATE TABLE IF NOT EXISTS status_change_fields (
    id INTEGER PRIMARY KEY,
    status_change_id integer not null,
    field_id integer not null,
    FOREIGN KEY(status_change_id) REFERENCES status_changes(id),
    FOREIGN KEY(field_id) REFERENCES fields(id)
);
CREATE TABLE IF NOT EXISTS fields (
    id INTEGER PRIMARY KEY,
    label text NOT NULL,
    input_type text NOT NULL,
    project_id integer not null,
    done boolean not null,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);
CREATE TABLE IF NOT EXISTS ticket_templates (
    id INTEGER PRIMARY KEY,
    label text NOT NULL,
    input_type text NOT NULL,
    project_id integer not null,
    done boolean not null,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);
CREATE TABLE IF NOT EXISTS tickets (
    id INTEGER PRIMARY KEY,
    label text NOT NULL,
    input_type text NOT NULL,
    project_id integer not null,
    done boolean not null,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);