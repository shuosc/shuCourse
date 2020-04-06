package campus

import "shuCourse/infrastructure"

var Campuses = map[string]uint8{}

func init() {
	rows, err := infrastructure.DB.Query(`
	SELECT id, unnest(name)
	FROM campus;
	`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id uint8
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		Campuses[name] = id
	}
}
