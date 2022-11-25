--liquibase formatted sql

--changeset Ali Khattab:Initial Script

create sequence aquarius.hibernate_sequence start with 1 increment by 1;

create table aquarius.task (task_id bigint not null, task_name varchar(255), task_description varchar(255), owner varchar(255), due_date timestamp, creation_date timestamp);

create table aquarius.user (user_id varchar(255) not null, dark_mode boolean, first_name varchar(255), surname varchar(255));

--rollback drop sequence aquarius.hibernate_sequence;

--rollback drop table aquarius.task;

--rollback drop table aquarius.user;