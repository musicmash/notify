from typing import List

from celery import group

from base.models import AsapSchedule
from base.models.connection import Connection
from base.serializers import NotificationSerializer
from base.services import TELEGRAM
from base.services.telegram.service import TelegramService
from notify.celery import app


@app.task(autoretry_for=(Exception,), max_retries=None, retry_backoff=True)
def add_releases(data: List[dict]):
    serializer = NotificationSerializer(data=data, many=True)
    serializer.is_valid(raise_exception=True)

    jobs = []

    for user_data in serializer.validated_data:

        user_name = user_data["user_name"]
        releases = user_data["releases"]

        connection = (
            AsapSchedule.objects.select_related("connection")
            .filter(connection__user_name=user_name)
            .count()
        )

        if not connection:
            continue

        job = group(
            send_user_release.s(user_name=user_name, release=release) for release in releases
        )

        jobs.append(job)

    group(jobs).delay()


@app.task(max_retries=None, retry_backoff=True)
def send_user_release(user_name: str, release: dict):
    user_connections = Connection.objects.filter(user_name=user_name).all()

    group_task = group(
        send_release.s(connection_id=connection.id, release=release)
        for connection in user_connections
    )

    group_task()


@app.task(max_retries=None, retry_backoff=True)
def send_release(connection_id: int, release: dict):
    connection: Connection = Connection.objects.get(pk=connection_id)

    services = {
        TELEGRAM: TelegramService,
    }

    service = services[connection.provider_id]()

    service.send_release(connection=connection, release=release)
