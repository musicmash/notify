from django.db import models

from .base import BaseModel


class DailySchedule(BaseModel):

    connection = models.ForeignKey("base.Connection", on_delete=models.CASCADE)
    notify_at = models.TimeField(auto_now=False, auto_now_add=False)

    class Meta:
        db_table = "daily_schedules"
        verbose_name = "Daily schedule"
        verbose_name_plural = "Daily schedules"

    def __str__(self):
        return f"{self.connection.user_name} {self.notify_at}"
