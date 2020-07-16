from typing import List, Optional

from celery import group

from base.models import Connection
from base.serializers import NotificationSerializer
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

        connection: Optional[Connection] = Connection.objects.filter(
            user_name=user_name, provider_name="telegram"
        ).first()

        if not connection:
            continue

        job = group(
            telegram_release.s(chat_id=connection.settings, user_name=user_name, release=release)
            for release in releases
        )

        jobs.append(job)

    group(jobs).delay()


@app.task(max_retries=None, retry_backoff=True)
def telegram_release(chat_id: int, user_name: str, release: dict):
    TelegramService().send_notication(chat_id=chat_id, user_name=user_name, release=release)
