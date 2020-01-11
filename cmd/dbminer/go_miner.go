package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/goodlandsecurity/go_miner/dbminer"
)

var (
	db       *sql.DB
	server        = flag.String("server", "", "the database server")
	port     *int = flag.Int("port", 1433, "the database port")
	user          = flag.String("user", "", "the database user")
	password      = flag.String("password", "", "the database password")
	debug         = flag.Bool("debug", false, "enable debugging")
)

type MSSQLMiner struct {
	Host *string
	Db   sql.DB
}

func New(host *string) (*MSSQLMiner, error) {
	m := MSSQLMiner{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MSSQLMiner) connect() error {
	if *debug {
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" user:%s\n", *user)
		fmt.Printf(" password:%s\n", *password)
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%v;port=%d;",
		*server,
		*user,
		*password,
		*port)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Failed to connect:", err.Error())
	}
	fmt.Println("Connected to server!\n")
	m.Db = *db
	return nil
}

func (m *MSSQLMiner) GetSchema() (*dbminer.Schema, error) {
	var (
		name    string
		dbnames []string
	)
	var s = new(dbminer.Schema)

	dbsql := fmt.Sprintf("SELECT name FROM sys.databases WHERE name NOT IN ('master', 'tempdb', 'model', 'msdb')")

	fmt.Println("Searching for databases!\n")
	dbrows, err := m.Db.Query(dbsql)
	if err != nil {
		log.Fatal("Scan failed!")
	}
	defer dbrows.Close()

	for dbrows.Next() {
		err := dbrows.Scan(&name)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}
		dbnames = append(dbnames, name)

	}

	for _, dbname := range dbnames {
		tsql := fmt.Sprintf("USE %s; SELECT table_catalog as table_schema, table_name, column_name FROM information_schema.columns ORDER BY table_schema, table_name", dbname)

		schemarows, err := m.Db.Query(tsql)
		if err != nil {
			log.Fatal("Query failed:", err.Error())
		}
		defer schemarows.Close()

		var prevschema, prevtable string
		var db dbminer.Database
		var table dbminer.Table
		for schemarows.Next() {
			var table_schema, table_name, column_name string
			err := schemarows.Scan(&table_schema, &table_name, &column_name)
			if err != nil {
				log.Fatal(err)
			}

			if table_schema != prevschema {
				if prevschema != "" {
					db.Tables = append(db.Tables, table)
					s.Databases = append(s.Databases, db)
				}
				db = dbminer.Database{Name: table_schema, Tables: []dbminer.Table{}}
				prevschema = table_schema
				prevtable = ""
			}

			if table_name != prevtable {
				if prevtable != "" {
					db.Tables = append(db.Tables, table)
				}
				table = dbminer.Table{Name: table_name, Columns: []string{}}
				prevtable = table_name
			}
			table.Columns = append(table.Columns, column_name)
		}
		db.Tables = append(db.Tables, table)
		s.Databases = append(s.Databases, db)
		if err := schemarows.Err(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func main() {
	flag.Parse()
	mm, err := New(server)
	if err != nil {
		panic(err)
	}
	defer mm.Db.Close()

	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
