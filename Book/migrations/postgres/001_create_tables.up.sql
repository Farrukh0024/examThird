create table books (
    id uuid primary key not null,
    name varchar(30),
    author_name varchar(30),
    page_number int not null
);