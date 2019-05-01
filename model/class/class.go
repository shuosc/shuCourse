package class

import (
	"github.com/lib/pq"
	"shuCourse/infrastructure"
)

type Time struct {
	Weeks       []int64 `json:"weeks"`
	Weekday     uint8   `json:"weekday"`
	BeginSector uint8   `json:"begin_sector"`
	EndSector   uint8   `json:"end_sector"`
}

type Class struct {
	Time
	CourseByTeacherId uint64 `json:"course_by_teacher_id"`
	CampusId          uint8  `json:"campus_id"`
	Place             string `json:"place"`
}

func Get(courseByTeacherId uint64) ([]Class, error) {
	rows, err := infrastructure.DB.Query(`
	SELECT campus_id,place,weeks,weekday,begin_sector,end_sector
	FROM Class
	WHERE course_by_teacher_id=$1;
	`, courseByTeacherId)
	if err != nil {
		return nil, err
	}
	var result []Class
	for rows.Next() {
		currentResult := Class{CourseByTeacherId: courseByTeacherId}
		err := rows.Scan(&currentResult.CampusId, &currentResult.Place, pq.Array(&currentResult.Weeks), &currentResult.Weekday, &currentResult.BeginSector, &currentResult.EndSector)
		if err != nil {
			return nil, err
		}
		result = append(result, currentResult)
	}
	return result, nil
}

func Save(class Class) {
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO Class(course_by_teacher_id, campus_id, place, weeks, weekday, begin_sector, end_sector)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, class.CourseByTeacherId, class.CampusId, class.Place, pq.Array(class.Weeks), class.Weekday, class.BeginSector, class.EndSector)
}
