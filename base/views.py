from django.db import IntegrityError
from rest_framework import status
from rest_framework.decorators import api_view
from rest_framework.exceptions import ValidationError
from rest_framework.parsers import JSONParser
from rest_framework.response import Response
from rest_framework.viewsets import ModelViewSet

from base.authentication import UserNameAuthentication
from base.models.connection import Connection
from base.tasks import add_releases

from .serializers import ConnectionSerializer, NotificationSerializer


@api_view(["POST"])
def new_releases(request):
    data = JSONParser().parse(request)

    serializer = NotificationSerializer(data=data, many=True)
    serializer.is_valid(raise_exception=True)

    add_releases.delay(data)

    return Response(status=status.HTTP_200_OK)


class ConnectionsViewSet(ModelViewSet):
    authentication_classes = [UserNameAuthentication]

    queryset = Connection.objects.all()
    serializer_class = ConnectionSerializer

    def get_queryset(self):
        return Connection.objects.filter(user_name=self.request.user.name).all()

    def create(self, request, *args, **kwargs):
        request.data["user_name"] = request.user.name

        try:
            response = super().create(request, *args, **kwargs)
        except IntegrityError:
            raise ValidationError("Connection already exists")

        return response

    def update(self, request, *args, **kwargs):
        request.data["user_name"] = request.user.name

        return super().update(request, *args, **kwargs)

    def partial_update(self, request, *args, **kwargs):
        request.data["user_name"] = request.user.name

        return super().partial_update(request, *args, **kwargs)
