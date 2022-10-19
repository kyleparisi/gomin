package app

import (
    "time"
)

type AppUser struct {
  CreatedAt time.Time
  Email string
  Id int64
  Name string
  Password string
  UpdatedAt time.Time
}

type DjangoMigrations struct {
  App string
  Applied time.Time
  Id int
  Name string
}
