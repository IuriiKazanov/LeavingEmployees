CREATE DATABASE leaving_employees;

USE leaving_employees;

CREATE TABLE IF NOT EXISTS user
(
    userID varchar(255) NOT NULL,
    workspaceID varchar(255) NOT NULL,
    isDeleted bool NOT NULL,
    name varchar(255) NOT NULL,
    PRIMARY KEY (userID, workspaceID)
);
