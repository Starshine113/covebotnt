package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Starshine113/covebotnt/cbdb"
	"github.com/jackc/pgx/v4/pgxpool"
)

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

func initDB() (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), config.Auth.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	if err := initDBIfNotInitialised(db); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] There was an error while initialising the database: %v\n", time.Now().Format(time.RFC3339), err)
		os.Exit(1)
	}
	sugar.Infof("Connected to database.")
	// update DB if it's not updated
	sugar.Infof("Target database version: %v", cbdb.DBVersion)
	err = updateDB(db)
	if err != nil {
		sugar.Panicf("Error updating database: %v", err)
	}
	return db, nil
}

func initDBIfNotInitialised(db *pgxpool.Pool) error {
	var exists bool
	err := db.QueryRow(context.Background(), "select exists (select from information_schema.tables where table_schema = 'public' and table_name = 'info')").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // the database has been initialised so we're done
	}

	// ...it's not initialised and we have to do that
	_, err = db.Exec(context.Background(), initDBSql)
	if err != nil {
		return err
	}
	sugar.Infof("Successfully initialised the database.")
	return nil
}

func updateDB(db *pgxpool.Pool) (err error) {
	var dbVersion int
	err = db.QueryRow(context.Background(), "select schema_version from info").Scan(&dbVersion)
	if err != nil {
		return err
	}
	initialDBVersion := dbVersion
	for dbVersion < cbdb.DBVersion {
		_, err = db.Exec(context.Background(), cbdb.DBVersions[dbVersion-1])
		if err != nil {
			return err
		}
		dbVersion++
		sugar.Infof("Updated database to version %v", dbVersion)
	}
	if initialDBVersion < cbdb.DBVersion {
		sugar.Infof("Successfully updated database to target version %v", cbdb.DBVersion)
	}
	return nil
}
