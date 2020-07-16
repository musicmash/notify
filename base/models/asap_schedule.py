from django.db import models

from .base import BaseModel


class AsapSchedule(BaseModel):

    connection = models.ForeignKey("base.Connection", on_delete=models.CASCADE)

    class Meta:
        db_table = "asap_schedules"
        verbose_name = "ASAP schedule"
        verbose_name_plural = "ASAP schedules"

    def __str__(self):
        return f"{self.connection.user_name}"
