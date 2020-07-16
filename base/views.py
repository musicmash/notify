from rest_framework import status
from rest_framework.decorators import api_view
from rest_framework.parsers import JSONParser
from rest_framework.response import Response

from base.tasks import add_releases

from .serializers import NotificationSerializer


@api_view(["POST"])
def new_releases(request):
    data = JSONParser().parse(request)

    serializer = NotificationSerializer(data=data, many=True)
    serializer.is_valid(raise_exception=True)

    add_releases.delay(data)

    return Response(status=status.HTTP_200_OK)
