# Generated by Django 3.0.8 on 2020-07-16 19:13

from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = []

    operations = [
        migrations.CreateModel(
            name="Provider",
            fields=[
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
                ("name", models.CharField(max_length=50, primary_key=True, serialize=False)),
            ],
            options={
                "verbose_name": "Provider",
                "verbose_name_plural": "Providers",
                "db_table": "providers",
            },
        ),
    ]
