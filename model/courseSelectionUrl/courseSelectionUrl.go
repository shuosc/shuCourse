package courseSelectionUrl

import "shuCourse/infrastructure"

func GetBySemesterId(id uint64) (string, error) {
	var result string
	row := infrastructure.DB.QueryRow(`
	SELECT url
	FROM CourseSelectionUrl
	WHERE semester_id=$1;
	`, id)
	err := row.Scan(&result)
	return result, err
}

func Save(semesterId uint64, url string) {
	if url[:4] != "http" {
		Save(semesterId, "http://"+url)
	}
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO CourseSelectionUrl(semester_id, url)
	VALUES ($1,$2);
	`, semesterId, url)
}
