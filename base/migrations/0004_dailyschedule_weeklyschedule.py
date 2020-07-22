# Generated by Django 3.0.8 on 2020-07-16 21:23

import django.db.models.deletion
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("base", "0003_asapschedule_connection"),
    ]

    operations = [
        migrations.CreateModel(
            name="WeeklySchedule",
            fields=[
                (
                    "id",
                    models.AutoField(
                        auto_created=True, primary_key=True, serialize=False, verbose_name="ID"
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
                ("day", models.PositiveSmallIntegerField(verbose_name="Day of week")),
                ("notify_at", models.TimeField()),
                (
                    "connection",
                    models.ForeignKey(
                        on_delete=django.db.models.deletion.CASCADE, to="base.Connection"
                    ),
                ),
            ],
            options={
                "verbose_name": "Weekly schedule",
                "verbose_name_plural": "Weekly schedules",
                "db_table": "weekly_schedules",
            },
        ),
        migrations.CreateModel(
            name="DailySchedule",
            fields=[
                (
                    "id",
                    models.AutoField(
                        auto_created=True, primary_key=True, serialize=False, verbose_name="ID"
                    ),
                ),
                ("created_at", models.DateTimeField(auto_now_add=True)),
                ("updated_at", models.DateTimeField(auto_now=True)),
                ("notify_at", models.TimeField()),
                (
                    "connection",
                    models.ForeignKey(
                        on_delete=django.db.models.deletion.CASCADE, to="base.Connection"
                    ),
                ),
            ],
            options={
                "verbose_name": "Daily schedule",
                "verbose_name_plural": "Daily schedules",
                "db_table": "daily_schedules",
            },
        ),
    ]