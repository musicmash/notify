# Generated by Django 3.0.8 on 2020-07-17 07:52

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ("base", "0004_dailyschedule_weeklyschedule"),
    ]

    operations = [
        migrations.AlterUniqueTogether(
            name="connection",
            unique_together={("user_name", "provider", "settings")},
        ),
    ]
