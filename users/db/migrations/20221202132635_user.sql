-- migrate:up
CREATE TABLE `users` (`id` string PRIMARY KEY,`name` string NULL);

-- migrate:down
drop table users;

