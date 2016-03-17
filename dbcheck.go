// Copyright 2016 Gabriel Guzman <gabe@lifewaza.com>.
// All rights reserved.  Use of this source code is
// governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/ziutek/mymysql/godrv"
)

var host, port, database, user, password string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Connect to host.")
	flag.StringVar(&port, "port", "3306", "Port number to use for connection.")
	flag.StringVar(&database, "database", "", "Database to use. If empty, will check all databases")
	flag.StringVar(&user, "user", "", "User for login.")
	flag.StringVar(&password, "password", "", "Password to use when connecting to server.")
	flag.Parse()
}

func main() {
	// Connect to the mysql server
	con, err := sql.Open("mymysql", "tcp:"+host+":"+port+"*"+database+"/"+user+"/"+password)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer con.Close()

	checkServerVersion(con)
	checkServerSettings(con)
	checkTables(con)
	checkColumns(con)
}

func checkServerVersion(con *sql.DB) {
	// Minimum supported version for utf8mb4 is 5.5.3
	row := con.QueryRow("SELECT version()")

	var fullVersion string
	var version []string
	err := row.Scan(&fullVersion)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	version = strings.Split(fullVersion, "-")
	fmt.Printf("Version: %s\n", version[0])
	var major, minor, release int
	ver := strings.Split(version[0], ".")
	major, err = strconv.Atoi(ver[0])
	minor, err = strconv.Atoi(ver[1])
	release, err = strconv.Atoi(ver[2])
	if !(major >= 5 && minor >= 5 && release >= 3) {
		log.Fatal("MySQL server version must be 5.5.3 or higher to support utf8mb4.")
	}
}

func checkServerSettings(con *sql.DB) {
	rows, err := con.Query("SHOW VARIABLES WHERE Variable_name LIKE 'character_set_%' OR Variable_name LIKE 'collation%'")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	var vname, value string
	for rows.Next() {
		err = rows.Scan(&vname, &value)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		if value != "utf8mb4" {
			fmt.Printf("Variable %s must be set to utf8mb4 (currently %s).\n", vname, value)
		}
	}
}

func checkTables(con *sql.DB) {
	// Based on: http://stackoverflow.com/a/6103613
	dbName := "ecom-middleware"
	tableName := "ecom_account"
	rows, err := con.Query("select c.character_set_name from information_schema.tables as t, information_schema.collation_character_set_applicability as c where c.collation_name = t.table_collation and t.table_schema = '" + dbName + "' and t.table_name = '" + tableName + "'")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	var charSetName string
	for rows.Next() {
		err = rows.Scan(&charSetName)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("%s.%s is not utf8mb4 (currently %s).\n", dbName, tableName, charSetName)
	}
}

func checkColumns(con *sql.DB) {
	rows, err := con.Query("SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM information_schema.columns WHERE CHARACTER_SET_NAME IS NOT NULL AND (CHARACTER_SET_NAME != 'utf8mb4' OR COLLATION_NAME != 'utf8mb4_unicode_ci')")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	var tableSchema, tableName, columnName string
	for rows.Next() {
		err = rows.Scan(&tableSchema, &tableName, &columnName)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("%s.%s.%s is not utf8mb4.\n", tableSchema, tableName, columnName)
	}
}

func createArtisanMigration() {
	const artisanMigration = `
<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class {{.migrationName}} extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
`
}
