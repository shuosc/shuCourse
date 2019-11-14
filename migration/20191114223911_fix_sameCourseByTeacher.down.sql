alter table coursebyteacher
    drop constraint sameCourseByTeacher;

alter table coursebyteacher
    add constraint sameCourseByTeacher EXCLUDE (course_id WITH =, in_course_teacher_id WITH =);
