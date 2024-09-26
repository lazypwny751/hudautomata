package lib

import (
    "database/sql"
    //"fmt"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func SetupDatabase(dbpath string) {
    db, err := sql.Open("sqlite3", dbpath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    createAppTable := `CREATE TABLE IF NOT EXISTS App (
        Version INTEGER
    );`
    _, err = db.Exec(createAppTable)
    if err != nil {
        log.Fatalf("App table couldn't created: %v", err)
    }

    createUserTable := `CREATE TABLE IF NOT EXISTS User (
        Name TEXT,
        PicturePath TEXT,
        Room TEXT,
        RoomID INTEGER,
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        Balance INTEGER
    );`
    _, err = db.Exec(createUserTable)
    if err != nil {
        log.Fatalf("User table couldn't created: %v", err)
    }

    createAdminTable := `CREATE TABLE IF NOT EXISTS Admin (
        UserName TEXT,
        Password TEXT,
        Price INTEGER
    );`
    _, err = db.Exec(createAdminTable)
    if err != nil {
        log.Fatalf("Admin table couldn't created: %v", err)
    }

	/*
		// It could be added with "verbose" feature in future.
	    fmt.Println("Database is ready for connection at: ", dbpath)
	*/
}
