from django.contrib.auth.models import AnonymousUser
from rest_framework import authentication
from rest_framework.exceptions import AuthenticationFailed


class UserName(AnonymousUser):
    def __init__(self, username):
        self.name = username

    @property
    def is_authenticated(self):
        return True


class UserNameAuthentication(authentication.BaseAuthentication):
    def authenticate(self, request):
        user_name = request.META.get("HTTP_X_USER_NAME", None)
        if not user_name:
            raise AuthenticationFailed("User name is required")

        return (UserName(user_name), None)
