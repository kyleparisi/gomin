from django.db import models


class User(models.Model):
    id = models.BigAutoField(primary_key=True)
    name = models.CharField(max_length=300)
    email = models.EmailField(null=False)
    password = models.CharField(max_length=300)
