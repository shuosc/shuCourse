create table Token
(
    id         bigserial    not null,
    token_hash varchar(256) not null
);

create unique index Token_id_uindex
    on Token (id);

create unique index Token_tokenHash_uindex
    on Token (token_hash);

alter table Token
    add constraint Token_pk
        primary key (id);

create table CourseSelectionUrl
(
    semester_id bigint       not null,
    url         varchar(128) not null
);

create table Campus
(
    id   smallserial  not null,
    name varchar(128) not null
);

create unique index Campus_id_uindex
    on Campus (id);

create unique index Campus_name_uindex
    on Campus (name);

alter table Campus
    add constraint Campus_pk
        primary key (id);

INSERT INTO Campus(id, name)
VALUES (1, '本部'),
       (2, '延长'),
       (3, '嘉定');

create table Course
(
    id     varchar(16)  not null,
    name   varchar(128) not null,
    credit smallint     not null
);

create unique index Course_id_uindex
    on Course (id);

alter table Course
    add constraint Course_pk
        primary key (id);

create table CourseByTeacher
(
    id                   bigserial   not null,
    course_id            varchar(16) not null
        constraint CourseByTeacher_Course_id_fk
            references Course (id)
            on update cascade on delete cascade,
    semester_id          bigint      not null,
    in_course_teacher_id varchar(8)  not null,
    teacher_id           varchar(16) not null,
    constraint sameCourseByTeacher EXCLUDE (course_id WITH =, in_course_teacher_id WITH =)
);

create unique index CourseByTeacher_id_uindex
    on CourseByTeacher (id);

alter table CourseByTeacher
    add constraint CourseByTeacher_pk
        primary key (id);

create table Class
(
    course_by_teacher_id bigint       not null
        constraint Class_CourseByTeacher_id_fk
            references CourseByTeacher (id)
            on update cascade on delete cascade,
    campus_id            smallint     not null
        constraint Class_Campus_id_fk
            references Campus (id)
            on update cascade on delete cascade,
    place                varchar(128) not null,
    weeks                smallint[]   not null,
    weekday              smallint     not null,
    begin_sector         smallint     not null,
    end_sector           smallint     not null
);
