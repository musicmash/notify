import os

import sentry_sdk
from sentry_sdk.integrations.celery import CeleryIntegration
from sentry_sdk.integrations.django import DjangoIntegration
from sentry_sdk.integrations.redis import RedisIntegration

sentry_sdk.init(
    dsn=os.getenv("SENTRY_DSN"),
    integrations=[DjangoIntegration(), CeleryIntegration(), RedisIntegration()],
)

# Build paths inside the project like this: os.path.join(BASE_DIR, ...)
BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

APP_NAME = os.getenv("APP_NAME")

# Celery settings
REDIS_HOST = os.getenv("REDIS_HOST")
REDIS_PORT = os.getenv("REDIS_PORT")
REDIS_PASS = os.getenv("REDIS_PASS")

CELERY_BROKER_URL = f"redis://:{REDIS_PASS}@{REDIS_HOST}:{REDIS_PORT}/0"
CELERY_RESULT_BACKEND = f"redis://:{REDIS_PASS}@{REDIS_HOST}:{REDIS_PORT}/0"

#: Only add pickle to this list if your broker is secured
#: from unwanted access (see userguide/security.html)
CELERY_ACCEPT_CONTENT = ["json"]
CELERY_TASK_SERIALIZER = "json"

# Quick-start development settings - unsuitable for production
# See https://docs.djangoproject.com/en/3.0/howto/deployment/checklist/

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = os.getenv("SECRET_KEY")

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = False

ALLOWED_HOSTS = ["*"]


# Application definition

INSTALLED_APPS = [
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "rest_framework",
    "base",
]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "django.middleware.common.CommonMiddleware",
]

ROOT_URLCONF = f"{APP_NAME}.urls"


LOGGING = {
    "version": 1,
    "disable_existing_loggers": False,
    "handlers": {"console": {"class": "logging.StreamHandler"}},
    "loggers": {
        "django": {"handlers": ["console"], "level": os.getenv("DJANGO_LOG_LEVEL", "INFO")}
    },
}


TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ]
        },
    }
]

DJANGO_TELEGRAMBOT = {
    "MODE": "WEBHOOK",  # (Optional [str]) # The default value is WEBHOOK,
    # otherwise you may use 'POLLING'
    # NB: if use polling you must provide to run
    # a management command that starts a worker
    "WEBHOOK_SITE": "https://musicmash.me",
    "WEBHOOK_PREFIX": "/v1/c8f2a745-48f1-4529-85e0-e4329bbfd679/telegram/web-hook",
    # (Optional[str]) # If this value is specified,
    # a prefix is added to webhook url
    # 'WEBHOOK_CERTIFICATE' : 'cert.pem', # If your site use self-signed
    # certificate, must be set with location of your public key
    # certificate.(More info at https://core.telegram.org/bots/self-signed )
    "STRICT_INIT": True,  # If set to True, the server will fail to start if some of the
    # apps contain telegrambot.py files that cannot be successfully
    # imported.
    "DISABLE_SETUP": False,  # If set to True, there will be no tries to set webhook or read
    # updates from the telegram server on app's start
    # (useful when developing on local machine; makes django's startup faster)
    # (Optional [str])  # The default name for file name containing telegram handlers
    # which has to be placed inside your local app(s).
    # Default is 'telegrambot'. Example is to put "telegrambot_handlers.py" file to local app's folder.
    "BOT_MODULE_NAME": "telegrambot",
    "BOTS": [
        {
            # "ID": "MainBot",  # Unique identifier for your bot (used in your code only)
            "TOKEN": os.getenv("TELEGRAM_BOT_TOKEN"),  # Your bots token (provided by botfather)
            "CONTEXT": True,  # Use context based handler functions (should be true for future versions)
            # 'ALLOWED_UPDATES':(Optional[list[str]]), # List the types of
            # updates you want your bot to receive. For example, specify
            # ``["message", "edited_channel_post", "callback_query"]`` to
            # only receive updates of these types. See ``telegram.Update``
            # for a complete list of available update types.
            # Specify an empty list to receive all updates regardless of type
            # (default). If not specified, the previous setting will be used.
            # Please note that this parameter doesn't affect updates created
            # before the call to the setWebhook, so unwanted updates may be
            # received for a short period of time.
            # 'TIMEOUT':(Optional[int|float]), # If this value is specified,
            # use it as the read timeout from the server
            # 'WEBHOOK_MAX_CONNECTIONS':(Optional[int]), # Maximum allowed number of
            # simultaneous HTTPS connections to the webhook for update
            # delivery, 1-100. Defaults to 40. Use lower values to limit the
            # load on your bot's server, and higher values to increase your
            # bot's throughput.
            # 'MESSAGEQUEUE_ENABLED':(Optinal[bool]), # Make this True if you want to use messagequeue
            # 'MESSAGEQUEUE_ALL_BURST_LIMIT':(Optional[int]), # If not provided 29 is the default value
            # 'MESSAGEQUEUE_ALL_TIME_LIMIT_MS':(Optional[int]), # If not provided 1024 is the default value
            # 'MESSAGEQUEUE_REQUEST_CON_POOL_SIZE':(Optional[int]), # If not provided 8 is the default value
            # 'POLL_INTERVAL' : (Optional[float]), # Time to wait between polling updates from Telegram in
            # seconds. Default is 0.0
            # 'POLL_CLEAN':(Optional[bool]), # Whether to clean any pending updates on Telegram servers before
            # actually starting to poll. Default is False.
            # 'POLL_BOOTSTRAP_RETRIES':(Optional[int]), # Whether the bootstrapping phase of the `Updater`
            # will retry on failures on the Telegram server.
            # |   < 0 - retry indefinitely
            # |     0 - no retries (default)
            # |   > 0 - retry up to X times
            # 'POLL_READ_LATENCY':(Optional[float|int]), # Grace time in seconds for receiving the reply from
            # server. Will be added to the `timeout` value and used as the read timeout from
            # server (Default: 2).
        },
        # Other bots here with same structure.
    ],
}

CACHES = {
    "default": {
        "BACKEND": "django_redis.cache.RedisCache",
        "LOCATION": f"redis://:{REDIS_PASS}@{REDIS_HOST}:{REDIS_PORT}/4",
        "OPTIONS": {
            "CLIENT_CLASS": "django_redis.client.DefaultClient",
        },
    }
}

WSGI_APPLICATION = f"{APP_NAME}.wsgi.application"


REST_FRAMEWORK = {
    "DEFAULT_PAGINATION_CLASS": "rest_framework.pagination.LimitOffsetPagination",
    "PAGE_SIZE": 100,
    "EXCEPTION_HANDLER": "base.utils.custom_exception_handler",
    "DEFAULT_AUTHENTICATION_CLASSES": [],
    "DEFAULT_PERMISSION_CLASSES": [],
    "UNAUTHENTICATED_USER": None,
}


# Database
# https://docs.djangoproject.com/en/3.0/ref/settings/#databases

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.postgresql",
        "NAME": os.getenv("POSTGRES_DB_NAME"),
        "USER": os.getenv("POSTGRES_DB_USER"),
        "PASSWORD": os.getenv("POSTGRES_DB_PASSWORD"),
        "HOST": os.getenv("POSTGRES_HOST"),
        "PORT": os.getenv("POSTGRES_PORT"),
        "CONN_MAX_AGE": 60,
    }
}


# Internationalization
# https://docs.djangoproject.com/en/3.0/topics/i18n/

LANGUAGE_CODE = "en-us"

TIME_ZONE = "UTC"

USE_I18N = False

USE_L10N = True

USE_TZ = True


# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/3.0/howto/static-files/

STATIC_URL = "/static/"
STATIC_ROOT = "static"
