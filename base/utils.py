from django.db import IntegrityError
from django.http import Http404
from rest_framework import exceptions
from rest_framework.exceptions import PermissionDenied
from rest_framework.response import Response
from rest_framework.views import set_rollback

from base.exceptions import AlreadyExistError


def custom_exception_handler(exc, context):
    """
    Returns the response that should be used for any given exception.

    By default we handle the REST framework `APIException`, and also
    Django's built-in `Http404` and `PermissionDenied` exceptions.

    Any unhandled exceptions may return `None`, which will cause a 500 error
    to be raised.
    """
    if isinstance(exc, Http404):
        exc = exceptions.NotFound()
    elif isinstance(exc, PermissionDenied):
        exc = exceptions.PermissionDenied()

    if isinstance(exc, IntegrityError):
        exc = AlreadyExistError()

    if isinstance(exc, exceptions.APIException):
        headers = {}
        if getattr(exc, "auth_header", None):
            headers["WWW-Authenticate"] = exc.auth_header
        if getattr(exc, "wait", None):
            headers["Retry-After"] = "%d" % exc.wait

        if isinstance(exc.detail, (list, dict)):
            data = exc.detail
        else:
            data = {"error": exc.detail.lower().rstrip(".")}

        set_rollback()
        return Response(data, status=exc.status_code, headers=headers)

    return None
