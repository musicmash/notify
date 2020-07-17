# Generated by Django 3.0.8 on 2020-07-17 08:01

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("base", "0005_auto_20200717_0752"),
    ]

    operations = [
        migrations.AlterUniqueTogether(name="connection", unique_together=set(),),
        migrations.AddConstraint(
            model_name="connection",
            constraint=models.UniqueConstraint(
                fields=("user_name", "provider", "settings"), name="Unique connection"
            ),
        ),
    ]
