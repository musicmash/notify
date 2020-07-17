from django.db import models

from .base import BaseModel


class Connection(BaseModel):

    user_name = models.CharField(max_length=255)
    provider = models.ForeignKey(
        "base.Provider",
        db_column="provider_name",
        verbose_name="Provider name",
        on_delete=models.CASCADE,
    )
    settings = models.CharField(max_length=255)

    class Meta:
        unique_together = [["user_name", "provider", "settings"]]
        db_table = "connections"
        verbose_name = "Connection"
        verbose_name_plural = "Connections"

    def __str__(self):
        return f"{self.user_name} - {self.provider.name} ({self.settings})"
