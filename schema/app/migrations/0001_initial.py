# Generated by Django 3.2.16 on 2022-10-18 21:14

from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='User',
            fields=[
                ('id', models.BigAutoField(primary_key=True, serialize=False)),
                ('created_at', models.DateTimeField(auto_now_add=True)),
                ('updated_at', models.DateTimeField(auto_now=True)),
                ('name', models.CharField(max_length=300)),
                ('email', models.EmailField(max_length=254)),
                ('password', models.CharField(max_length=300)),
            ],
        ),
    ]
