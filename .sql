CREATE TABLE users (
                      id     uuid  DEFAULT gen_random_uuid() PRIMARY KEY,
                      email    text  NOT NULL
);

CREATE TABLE tokens (
    user_id uuid references users(id),
    refresh text  NOT NULL ,
    ip text NOT NULL
);

insert into users (email) values ( 'ritenshe@gmail.ru');