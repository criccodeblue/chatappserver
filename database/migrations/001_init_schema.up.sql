CREATE TYPE ChatStatusType AS ENUM('pending', 'sent', 'received', 'read');

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    handle TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts(
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL ,
    contact_user_id SERIAL NOT NULL ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (contact_user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS chats(
    chat_id SERIAL PRIMARY KEY,
    sender_id SERIAL NOT NULL ,
    receiver_id SERIAL NOT NULL ,
    message_data TEXT,
    status ChatStatusType NOT NULL ,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);
