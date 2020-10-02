import logging

import telegram
from django.core.cache import cache
from django.db import IntegrityError
from django.template.loader import render_to_string
from django_telegrambot.apps import DjangoTelegramBot
from telegram.ext import CallbackContext, CommandHandler
from telegram.update import Update

from base.models.connection import Connection
from base.services.telegram.service import TelegramService

logger = logging.getLogger(__name__)


# Define a few command handlers. These usually take the two arguments bot and
# update. Error handlers also receive the raised TelegramError object in error.
def start(update: Update, context: CallbackContext):
    telegram_id = update.message.from_user.id

    message_words = update.message.text.split()

    if len(message_words) == 1:
        # ToDo: ask user to attach telegram account through website
        return

    payload = message_words[1]
    user_name = cache.get(payload)

    if not user_name:
        text = render_to_string("telegram/expired.djt")
    else:
        try:
            Connection.objects.create(
                user_name=user_name, settings=telegram_id, provider_name="telegram"
            )
        except IntegrityError:
            text = render_to_string("telegram/already.djt")
        else:
            text = render_to_string("telegram/attach.djt")

    TelegramService().bot.send_message(
        chat_id=telegram_id,
        text=text,
        parse_mode=telegram.ParseMode.MARKDOWN,
    )


def main():
    logger.info("Loading handlers for telegram bot")

    # Default dispatcher (this is related to the first bot in settings.TELEGRAM_BOT_TOKENS)
    dp = DjangoTelegramBot.dispatcher
    # To get Dispatcher related to a specific bot
    # dp = DjangoTelegramBot.getDispatcher('BOT_n_token')     #get by bot token
    # dp = DjangoTelegramBot.getDispatcher('BOT_n_username')  #get by bot username

    # on different commands - answer in Telegram
    dp.add_handler(CommandHandler("start", start))
