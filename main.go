package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int
	Username string
	Email    string
}

type Role struct {
	ID   int
	Name string
}

type UserRole struct {
	UserID int
	RoleID int
}

func main() {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Query users with their associated roles via the pivot table
	rows, err := db.Query(`
		SELECT u.id, u.username, u.email, r.id, r.name
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON ur.role_id = r.id
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var role Role
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &role.ID, &role.Name)
		if err != nil {
			panic(err)
		}
		user.Roles = append(user.Roles, role)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	// Print users with their associated roles
	for _, user := range users {
		fmt.Printf("User ID: %d\n", user.ID)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Println("Roles:")
		for _, role := range user.Roles {
			fmt.Printf("- ID: %d, Name: %s\n", role.ID, role.Name)
		}
		fmt.Println("-----------------------")
	}
}
