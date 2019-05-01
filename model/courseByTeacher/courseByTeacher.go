package courseByTeacher

import "shuCourse/infrastructure"

type CourseByTeacher struct {
	Id                uint64 `json:"id"`
	CourseId          string `json:"course_id"`
	SemesterId        uint64 `json:"semester_id"`
	InCourseTeacherId string `json:"in_course_teacher_id"`
	TeacherId         string `json:"teacher_id"`
}

func Get(id uint64) (CourseByTeacher, error) {
	result := CourseByTeacher{Id: id}
	row := infrastructure.DB.QueryRow(`
	SELECT course_id,semester_id, in_course_teacher_id, teacher_id
	FROM CourseByTeacher
	WHERE id=$1;
	`, id)
	err := row.Scan(&result.CourseId, &result.SemesterId, &result.InCourseTeacherId, &result.TeacherId)
	return result, err
}

func Save(courseByTeacher CourseByTeacher) (uint64, error) {
	row := infrastructure.DB.QueryRow(`
	INSERT INTO CourseByTeacher(course_id,semester_id, in_course_teacher_id, teacher_id)
	VALUES ($1,$2,$3,$4)
	RETURNING id;
	`, courseByTeacher.CourseId, courseByTeacher.SemesterId, courseByTeacher.InCourseTeacherId, courseByTeacher.TeacherId)
	var result uint64
	err := row.Scan(&result)
	return result, err
}
