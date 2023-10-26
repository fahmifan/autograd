PRAGMA foreign_keys = ON;

CREATE TABLE "users" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" INT NOT NULL,
    "created_at" DATETIMET NOT NULL,
    "updated_at" DATETIMET NOT NULL,
    "deleted_at" DATETIMET
);