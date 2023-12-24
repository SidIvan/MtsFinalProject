CREATE TABLE IF NOT EXISTS users (
   user_id     varchar(120) PRIMARY KEY,
   lat         varchar(72) NOT NULL,
   lng         varchar(300) UNIQUE NOT NULL
);