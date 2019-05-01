package course

import (
	"shuCourse/infrastructure"
)

type Course struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Credit uint8  `json:"credit"`
}

func Get(id string) (Course, error) {
	result := Course{Id: id}
	row := infrastructure.DB.QueryRow(`
	SELECT name, credit
	FROM Course
	WHERE id=$1;
	`, id)
	err := row.Scan(&result.Name, &result.Credit)
	return result, err
}

func Save(course Course) {
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO course(id, name, credit) 
	VALUES ($1,$2,$3)
	ON CONFLICT(id) DO
	UPDATE SET name=$2,
			   credit=$3;
	`, course.Id, course.Name, course.Credit)
}
