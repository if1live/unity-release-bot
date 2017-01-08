CREATE TABLE `versions` (
    `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
    `version` VARCHAR(15) NULL,
    `category` VARCHAR(15) NULL,
    `link` VARCHAR(255) NULL,
    `date` DATE,
    `created` DATE DEFAULT (datetime('now','localtime'))
);
CREATE UNIQUE INDEX ux01 ON versions(version);
