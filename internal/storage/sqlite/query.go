package sqlite

const createTableUrlQuery = `
create table if not exists url(
	id integer primary key,
	alias text not null unique,
	url text not null
);
create index if not exists idx_alias on url(alias);
`

const insertUrlQuery = `
insert into url (url, alias) values (?, ?);
`

const getUrlQuery = `
select url from url where alias = ?;
`
