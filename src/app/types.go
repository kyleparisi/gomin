package app

type AppUser struct {
  Email string
  Id int64
  Name string
  Password string
}

type DjangoMigrations struct {
  App string
  Id int
  Name string
}
