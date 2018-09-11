package main

// User describes the data structure for firestore user
type User struct {
	Name  string `firestore:"name"`
	Email string `firestore:"email"`
}
