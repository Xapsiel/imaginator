CREATE TABLE IF NOT EXISTS polls(
                                    id varchar(250) PRIMARY KEY ,
                                    message_id INTEGER NOT NULL ,
                                    poll_type STATUS_ENUM NOT NULL,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
