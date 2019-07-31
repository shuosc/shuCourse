package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"shuCourse/infrastructure"
	"shuCourse/model/class"
	courseModel "shuCourse/model/course"
	"shuCourse/model/courseByTeacher"
	"shuCourse/model/courseSelectionUrl"
	"strconv"
)

func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	course, err := courseModel.Get(id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	data, _ := json.Marshal(course)
	_, _ = w.Write(data)
}

func GetCourseByTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	course, err := courseByTeacher.Get(uint64(id))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	data, _ := json.Marshal(course)
	_, _ = w.Write(data)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Query().Get("course_by_teacher_id"))
	classObject, err := class.Get(uint64(id))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	data, _ := json.Marshal(classObject)
	_, _ = w.Write(data)
}

type DefaultFormat struct {
	Id                uint64        `json:"id"`
	CourseId          string        `json:"course_id"`
	InCourseTeacherId string        `json:"in_course_teacher_id"`
	TeacherId         string        `json:"teacher_id"`
	OnSemesterId      uint8         `json:"on_semester_id"`
	Name              string        `json:"name"`
	Credit            uint8         `json:"credit"`
	Classes           []class.Class `json:"classes"`
}

func ThreeIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	inCourseTeacherId := r.URL.Query().Get("in_course_teacher_id")
	courseId := r.URL.Query().Get("course_id")
	semesterId, err := strconv.Atoi(r.URL.Query().Get("semester_id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}
	result := DefaultFormat{
		CourseId:          courseId,
		InCourseTeacherId: inCourseTeacherId,
		OnSemesterId:      uint8(semesterId),
	}
	row := infrastructure.DB.QueryRow(`
	SELECT CourseByTeacher.id,
       CourseByTeacher.teacher_id,
       Course.name,
       Course.credit
	FROM Course, CourseByTeacher
	WHERE Course.id = CourseByTeacher.course_id
	  AND CourseByTeacher.in_course_teacher_id = $1
	  AND CourseByTeacher.course_id = $2
	  AND CourseByTeacher.semester_id = $3;
	`, inCourseTeacherId, courseId, semesterId)
	err = row.Scan(&result.Id, &result.TeacherId, &result.Name, &result.Credit)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	result.Classes, err = class.Get(result.Id)
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

func IdOnlyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	result := DefaultFormat{
		Id: uint64(id),
	}
	row := infrastructure.DB.QueryRow(`
	SELECT CourseByTeacher.teacher_id,
	       Course.id,
	       CourseByTeacher.in_course_teacher_id,
	       CourseByTeacher.semester_id,
	       Course.name,
	       Course.credit
	FROM Course,
	     CourseByTeacher
	WHERE CourseByTeacher.id = $1
	  AND Course.id = CourseByTeacher.course_id;
	`, id)
	err := row.Scan(&result.TeacherId, &result.CourseId, &result.InCourseTeacherId, &result.OnSemesterId, &result.Name, &result.Credit)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	result.Classes, err = class.Get(result.Id)
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

func DefaultFormatHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("id") != "" {
		IdOnlyHandler(w, r)
	} else {
		ThreeIdHandler(w, r)
	}
}

func CourseSelectionURLHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	url, err := courseSelectionUrl.GetBySemesterId(id)
	if err != nil {
		log.Println("user tried to get CourseSelectionURL for semester", id, "which does not exist")
		w.WriteHeader(404)
		return
	}
	response := struct {
		Url string `json:"url"`
	}{url}
	data, _ := json.Marshal(response)
	_, _ = w.Write(data)
}

func PingPongHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
