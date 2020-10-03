CREATE TABLE IF NOT EXISTS "reports"
(
    "id"           integer primary key autoincrement,
    "user_id"      integer,
    "resource"     varchar(191) NOT NULL,
    "filename"     varchar(191) NOT NULL,
    `type`         varchar(191) NOT NULL,
    "status"       varchar(191) NOT NULL,
    `note`         varchar(191),
    "created_at"   datetime,
    "updated_at"   datetime
);
