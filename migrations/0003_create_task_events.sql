CREATE TABLE events (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id int NOT NULL,
    event_type text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    constraint fk_task_id
    foreign key (task_id)
    references tasks (id)
);