from django.db import models

from .base import BaseModel


class WeeklySchedule(BaseModel):

    connection = models.ForeignKey("base.Connection", on_delete=models.CASCADE)
    day = models.PositiveSmallIntegerField("Day of week")
    notify_at = models.TimeField(auto_now=False, auto_now_add=False)

    class Meta:
        db_table = "weekly_schedules"
        verbose_name = "Weekly schedule"
        verbose_name_plural = "Weekly schedules"

    def __str__(self):
        return f"{self.connection.user_name} - {self.day} - {self.notify_at}"
