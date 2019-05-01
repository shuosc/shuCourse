package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"shuCourse/infrastructure"
	"shuCourse/service/crawl"
	"shuCourse/service/token"
	"strconv"
)

func addAdmin(username string) {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(username), -1)
	_, err := infrastructure.DB.Exec(`
	INSERT INTO token(token_hash)
	VALUES ($1);
	`, string(encrypted))
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(token.GenerateJWT(username))
}

func setSemesterUrl(semesterId string, url string) {
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO courseSelectionUrl(semester_id, url)
	VALUES ($1,$2);
	`, semesterId, url)
}

func main() {
	whichFunction := os.Args[1]
	switch whichFunction {
	case "addAdmin":
		username := os.Args[2]
		addAdmin(username)
	case "setSemesterUrl":
		semesterId := os.Args[2]
		url := os.Args[3]
		setSemesterUrl(semesterId, url)
	case "fetchCourses":
		semesterId, _ := strconv.Atoi(os.Args[2])
		username := os.Args[3]
		password := os.Args[4]
		crawl.FetchAllCourses(uint64(semesterId), username, password)
	}
}
