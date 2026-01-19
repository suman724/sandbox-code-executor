create table if not exists jobs (
  id text primary key,
  status text not null
);

create table if not exists sessions (
  id text primary key,
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
