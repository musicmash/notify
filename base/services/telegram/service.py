import os
from typing import Optional

import telegram
from django.template.loader import render_to_string


class TelegramService:

    def __init__(self, telegram_token: Optional[str] = None):
        token = telegram_token or os.getenv("TELEGRAM_BOT_TOKEN")

        self.bot = telegram.Bot(token)

    def send_notication(self, chat_id: int, user_name: str, release: dict):
        release_url = f"https://itunes.apple.com/us/{release['type']}/{release['itunes_id']}"

        keyboard = [
            [telegram.InlineKeyboardButton("Listen on Apple Music", url=release_url)]
        ]

        explicit = "🅴 " if release["explicit"] else ""

        text = render_to_string("telegram/notification.djt", {
            "invisible_text": "\u200c\u200c",
            "poster": release["poster"],
            "title": release["title"],
            "artist_name": release["artist_name"],
            "release_type": release["type"],
            "explicit": explicit,
        })

        self.bot.send_message(
            chat_id=chat_id,
            text=text,
            parse_mode=telegram.ParseMode.MARKDOWN_V2,
            reply_markup=telegram.InlineKeyboardMarkup(keyboard),
        )
