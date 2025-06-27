CREATE TABLE  IF NOT EXISTS image_messages(
                                              message_id int PRIMARY KEY ,
                                              image_id  TEXT,
                                              FOREIGN KEY (image_id) REFERENCES images(id)
);
