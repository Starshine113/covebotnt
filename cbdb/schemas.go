package cbdb

// DBVersion is the current database version
const DBVersion = 3

// DBVersions is a slice of schemas for every database version
var DBVersions []string = []string{
	`create table if not exists notes (
		id 			serial		primary key,
		guild_id 	text		not null,
		user_id 	text		not null,
		mod_id 		text		not null,
		note 		text		not null,
		created 	timestamp	not null default (current_timestamp at time zone 'utc')
	);
	
	update public.info set schema_version = 2;`,

	`create table if not exists mod_log (
		id			serial		primary key,
		guild_id	text		not null,
		user_id		text		not null,
		mod_id		text		not null,
		type		modaction	not null,
		reason		text		not null,
		created		timestamp	not null default (current_timestamp at time zone 'utc')
	);
	
	update public.info set schema_version = 3;`,
}
