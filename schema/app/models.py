from django.db import models


class User(models.Model):
    name = models.CharField(max_length=300)
    email = models.EmailField(null=False, db_index=True)
    password = models.CharField(max_length=300)
