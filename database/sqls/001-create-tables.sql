
create table IF not exists `students`
(id int not null auto_increment primary key, password varchar(128) not null, student_id varchar(100) not null unique, name varchar(100) , email varchar(100) not null, is_superuser bool not null, last_login datetime null)DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

create table IF not exists `classifications`
(id int not null auto_increment primary key, name varchar(50) not null) DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

create table IF not exists `equipments`
(id int not null auto_increment primary key, name varchar(50) not null, classifications_id int, user varchar(20) null, state varchar(20) not null DEFAULT '貸出可能', remarks text null, FOREIGN KEY fk_classifications_id(classifications_id) REFERENCES classifications (id) ON DELETE CASCADE ON UPDATE CASCADE) DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

create table IF not exists `sessions`
(session_key varchar(100) not null primary key, session_date datetime not null, session_data text not null);