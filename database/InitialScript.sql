CREATE DATABASE leaving_employees;

USE leaving_employees;

CREATE TABLE IF NOT EXISTS workspace (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    isActive BOOL NOT NULL DEFAULT TRUE,
    botAccessToken VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS user
(
    userID varchar(255) NOT NULL,
    workspaceID varchar(255) NOT NULL,
    isDeleted bool NOT NULL,
    name varchar(255) NOT NULL,
    PRIMARY KEY (userID, workspaceID)
#     FOREIGN KEY (workspaceID) REFERENCES workspace (id) ON DELETE CASCADE
);
