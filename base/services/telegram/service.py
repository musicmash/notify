import os
from html import unescape
from typing import Optional

import telegram
from dateutil.parser import parse
from django.template.loader import render_to_string
from django.utils import timezone

from base.models.connection import Connection

ANNOUNCE = "announce"
RELEASE = "release"


class TelegramService:
    def __init__(self, telegram_token: Optional[str] = None):
        token = telegram_token or os.getenv("TELEGRAM_BOT_TOKEN")

        self.bot = telegram.Bot(token)

    def send_release(self, connection: Connection, release: dict):
        release_url = f"https://itunes.apple.com/us/{release['type']}/{release['itunes_id']}"

        keyboard = [[telegram.InlineKeyboardButton("Listen on Apple Music", url=release_url)]]

        explicit = "🅴 " if release["explicit"] else ""
        released = parse(release["released"])

        release_types = {RELEASE: "telegram/release.djt", ANNOUNCE: "telegram/announce.djt"}

        release_type = ANNOUNCE if released > timezone.now() else RELEASE

        text = render_to_string(
            release_types[release_type],
            {
                "invisible_text": "\u200c\u200c",
                "poster": release["poster"],
                "title": unescape(release["title"]),
                "artist_name": unescape(release["artist_name"]),
                "release_type": release["type"],
                "explicit": explicit,
                "release_date": released,
            },
        )

        self.bot.send_message(
            chat_id=connection.settings,
            text=text,
            parse_mode=telegram.ParseMode.MARKDOWN,
            reply_markup=telegram.InlineKeyboardMarkup(keyboard),
        )
