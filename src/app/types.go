package app

type AppUser struct {
  Id int
  Name string
  Email string
  Password string
}

type DjangoMigrations struct {
  Id int
  App string
  Name string
}
