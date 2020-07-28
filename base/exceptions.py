from rest_framework import status
from rest_framework.exceptions import APIException


class AlreadyExistError(APIException):
    status_code = status.HTTP_400_BAD_REQUEST
    default_detail = "object already exist"
    default_code = "already_exist"
