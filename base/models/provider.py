from django.db import models

from .base import BaseModel


class Provider(BaseModel):

    name = models.CharField(max_length=50, primary_key=True)

    class Meta:
        db_table = "providers"
        verbose_name = "Provider"
        verbose_name_plural = "Providers"

    def __str__(self):
        return self.name
