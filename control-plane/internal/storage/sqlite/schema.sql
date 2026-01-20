create table if not exists jobs (
  id text primary key,
  status text not null
);

create table if not exists sessions (
  id text primary key,
  status text not null,
  runtime_id text
);

create table if not exists session_steps (
  id text primary key,
  session_id text not null,
  command text not null,
  status text not null
);

create table if not exists policies (
  id text primary key,
  version integer not null
);

create table if not exists audit_events (
  id text primary key,
  action text not null,
  outcome text not null
);

create table if not exists idempotency_keys (
  key text primary key,
  value text not null
);
