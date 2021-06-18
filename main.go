package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to a database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=root password=root")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}
	defer conn.Close()

	log.Println("Connected to database!")

	// test my connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping database!")
	}

	log.Println("Pinged database!")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert a row
	query := `INSERT INTO users (first_name, last_name) VALUES ($1, $2)`
	_, err = conn.Exec(query, "Jack", "Brown")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a row!")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// update a row
	stmt := `UPDATE users SET first_name = $1 WHERE id = $2`
	_, err = conn.Exec(stmt, "Jackie", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Updated one or more rows!")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// get one row by id
	query = `SELECT id, first_name, last_name FROM users WHERE id = $1`

	var firstName, lastName string
	var id int

	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("QueryRow returns", id, firstName, lastName)

	// delete a row
	query = `DELETE FROM users WHERE id = $1`
	_, err = conn.Exec(query, 2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted a row!")

	// get rows again
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("SELECT id, first_name, last_name FROM users;")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is", id, firstName, lastName)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}

	fmt.Println("-------------------")

	return nil
}
