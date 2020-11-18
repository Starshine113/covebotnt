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

// initDBSql is the initial SQL database schema
var initDBSql = `create type modaction as enum ('warn', 'mute', 'unmute', 'pause', 'hardmute', 'kick', 'tempban', 'ban');

create table if not exists guild_settings
(
	guild_id			text primary key,
	prefix				text default '',

    starboard_channel	text not null default '',
    react_limit			int not null default 100,
	emoji				text not null default '⭐',
    sender_can_react	boolean default false,
	react_to_starboard	boolean default true,
	
	mod_roles			text[] not null default array[]::text[],
	helper_roles		text[] not null default array[]::text[],
	mod_log				text not null default '',
	mute_role			text not null default '',
	pause_role			text not null default '',

	gatekeeper_roles	text[] not null default array[]::text[],
	member_roles		text[] not null default array[]::text[],
	gatekeeper_channel	text not null default '',
	gatekeeper_message	text not null default 'Please wait to be approved, {mention}.',
	welcome_channel		text not null default '',
	welcome_message		text not null default 'Welcome to {guild}, {mention}!'
);

create table if not exists starboard_messages
(
    message_id				text primary key,
    channel_id				text not null,
    server_id				text not null,
    starboard_message_id	text
);

create table if not exists starboard_blacklisted_channels
(
    channel_id	text primary key,
    server_id	text not null
);

create table if not exists info
(
    id						int primary key not null default 1, -- enforced only equal to 1
    schema_version			int,
    constraint singleton	check (id = 1) -- enforce singleton table/row
);

insert into info (schema_version) values (1);`
