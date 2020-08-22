package test_utils

import (
	db2 "db"
	"github.com/jmoiron/sqlx"
)

const (
	dropAccountsQuery           = `DROP TABLE IF EXISTS accounts`
	dropAccountCredentialsQuery = `DROP TABLE IF EXISTS account_credentials;`
	dropObjectsQuery            = `DROP TABLE IF EXISTS objects;`
	dropCommandsQuery           = `DROP TABLE IF EXISTS commands;`
	dropCommandsSettingsQuery   = `DROP TABLE IF EXISTS commands_settings;`
	dropCommandsHeadersQuery    = `DROP TABLE IF EXISTS commands_headers;`
	dropCommandsCookiesQuery    = `DROP TABLE IF EXISTS commands_cookies;`

	addAccountsQuery = `
	CREATE TABLE IF NOT EXISTS accounts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash VARCHAR(32) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	addAccountQuery = `
	INSERT INTO accounts(hash) VALUES(:hash)`
	addAccountCredentialsQuery = `
	INSERT INTO account_credentials(login, password, account_hash)
	VALUES(:login, :password, :hash)`
	addObjectQuery = `
	INSERT INTO objects(name, hash)
	VALUES(:name, :hash)`
	addCommandQuery = `
	INSERT INTO commands(name, hash, object_name)
	VALUES(:name, :hash, :object_name)`
	addCommandSettingsQuery = `
	INSERT INTO commands_settings(method, base_url, endpoint, pass_arguments_in_url, command_hash)
	VALUES(:method, :base_url, :endpoint, :pass_arguments_in_url, :command_hash)`
	addCommandHeadersQuery = `
	INSERT INTO commands_headers(key, value, command_hash)
	VALUES(:key, :value, :command_hash)`
	addCommandCookiesQuery = `
	INSERT INTO commands_cookies(key, value, command_hash)
	VALUES(:key, :value, :command_hash)`

	AccountHash       = "some-hash"
	CredentialsLogin  = "some_login"
	CredentialsHash   = "some_hash"
	ObjectName        = "USER"
	ObjectHash        = "hash-1"
	CreateCommandName = "CREATE"
	GetCommandName    = "GET"
	PatchCommandName  = "PATCH"
	DeleteCommandName = "DELETE"
	CreateCommandHash = "hash-1"
	GetCommandHash    = "hash-2"
	PatchCommandHash  = "hash-3"
	DeleteCommandHash = "hash-4"

	NotExistsAccountHash = "not-exists-account-hash"
	NotExistsObjectHash  = "not-exists-object-hash"
)

var (
	accounts = []map[string]interface{}{
		{
			"hash": AccountHash,
		},
	}
	credentials = []map[string]interface{}{
		{
			"login":    CredentialsLogin,
			"password": "some_password",
			"hash":     CredentialsHash,
		},
	}
	objects = []map[string]interface{}{
		{
			"name": ObjectName,
			"hash": ObjectHash,
		},
	}
	commands = []map[string]interface{}{
		{
			"name":        CreateCommandName,
			"hash":        CreateCommandHash,
			"object_name": ObjectName,
		},
		{
			"name":        GetCommandName,
			"hash":        GetCommandHash,
			"object_name": ObjectName,
		},
		{
			"name":        PatchCommandName,
			"hash":        PatchCommandHash,
			"object_name": ObjectName,
		},
		{
			"name":        DeleteCommandName,
			"hash":        DeleteCommandHash,
			"object_name": ObjectName,
		},
	}
	Settings = []map[string]interface{}{
		{
			"method":                "POST",
			"base_url":              "http://link.com",
			"endpoint":              "user/",
			"pass_arguments_in_url": false,
			"command_hash":          CreateCommandHash,
		},
		{
			"method":                "GET",
			"base_url":              "http://link.com",
			"endpoint":              "user",
			"pass_arguments_in_url": true,
			"command_hash":          GetCommandHash,
		},
		{
			"method":                "PATCH",
			"base_url":              "http://link.com",
			"endpoint":              "user",
			"pass_arguments_in_url": true,
			"command_hash":          PatchCommandHash,
		},
		{
			"method":                "DELETE",
			"base_url":              "http://link.com",
			"endpoint":              "user",
			"pass_arguments_in_url": true,
			"command_hash":          DeleteCommandHash,
		},
	}
	Headers = []map[string]interface{}{
		{
			"key":          "X-Test-1",
			"value":        "test1",
			"command_hash": CreateCommandHash,
		},
		{
			"key":          "X-Test-2",
			"value":        "test2",
			"command_hash": CreateCommandHash,
		},
	}
	Cookies = []map[string]interface{}{
		{
			"key":          "Test-Value-1",
			"value":        "test1",
			"command_hash": CreateCommandHash,
		},
		{
			"key":          "Test-Value-2",
			"value":        "test2",
			"command_hash": CreateCommandHash,
		},
	}

	queryToData = map[string][]map[string]interface{}{
		addAccountCredentialsQuery: credentials,
		addObjectQuery:             objects,
		addCommandQuery:            commands,
		addCommandSettingsQuery:    Settings,
		addCommandHeadersQuery:     Headers,
		addCommandCookiesQuery:     Cookies,
	}
)

func DropTables(db *sqlx.DB) {
	for _, query := range []string{
		dropCommandsQuery, dropCommandsSettingsQuery,
		dropCommandsHeadersQuery, dropCommandsCookiesQuery,
		dropObjectsQuery, dropAccountCredentialsQuery,
		dropAccountsQuery,
	} {
		exec(db, query)
	}
}

func DropCommandsSettings(db *sqlx.DB) {
	exec(db, dropCommandsSettingsQuery)
}

func DropCommandsHeaders(db *sqlx.DB) {
	exec(db, dropCommandsHeadersQuery)
}

func DropCommandsCookies(db *sqlx.DB) {
	exec(db, dropCommandsCookiesQuery)
}

func ReplaceBaseURLAndInitTables(db *sqlx.DB, baseURL string) {
	for settingIndex := range Settings {
		Settings[settingIndex]["base_url"] = baseURL
	}

	InitTables(db)
}

func InitTablesWithAccounts(db *sqlx.DB) {
	DropTables(db)

	_, err := db.Exec(addAccountsQuery)
	if err != nil {
		panic(err)
	}

	err = db2.Install(db)
	if err != nil {
		panic(err)
	}

	tx := db.MustBegin()
	for query, data := range queryToData {
		applyData(tx, query, data)
	}
	applyData(tx, addAccountQuery, accounts)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func InitTables(db *sqlx.DB) {
	DropTables(db)
	err := db2.Install(db)
	if err != nil {
		panic(err)
	}

	tx := db.MustBegin()
	for query, data := range queryToData {
		applyData(tx, query, data)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func exec(db *sqlx.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func applyData(tx *sqlx.Tx, query string, data []map[string]interface{}) {
	for _, item := range data {
		_, err := tx.NamedExec(query, item)
		if err != nil {
			panic(err)
		}
	}
}
