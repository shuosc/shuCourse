alter table coursebyteacher
    drop constraint sameCourseByTeacher;

alter table coursebyteacher
    add constraint sameCourseByTeacher EXCLUDE (semester_id WITH =, course_id WITH =, in_course_teacher_id WITH =);
