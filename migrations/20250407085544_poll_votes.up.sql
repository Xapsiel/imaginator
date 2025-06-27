CREATE TABLE IF NOT EXISTS poll_votes(
                                         user_id BIGINT PRIMARY KEY ,
                                         poll_id varchar(250)  UNIQUE ,
                                         answer_id int,
                                         FOREIGN KEY (poll_id) REFERENCES polls(id)
);