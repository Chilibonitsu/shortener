create table if not exists urls(
	id integer primary key,
	alias text not null unique,
	url text not null
);