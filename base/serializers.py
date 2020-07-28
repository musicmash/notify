from rest_framework import serializers
from rest_framework.exceptions import ValidationError

from base.models.connection import Connection

TYPES = [
    ("album", "album"),
    ("single", "single"),
    ("music-video", "music-video"),
]


class ReleaseSerializer(serializers.Serializer):
    artist_id = serializers.IntegerField()
    artist_name = serializers.CharField()
    title = serializers.CharField()
    released = serializers.DateTimeField()
    itunes_id = serializers.CharField(allow_null=True)
    spotify_id = serializers.CharField(allow_null=True)
    deezer_id = serializers.CharField(allow_null=True)
    poster = serializers.URLField()
    type = serializers.ChoiceField(choices=TYPES)
    explicit = serializers.BooleanField()

    def validate(self, attrs):
        store_fields = (attrs["itunes_id"], attrs["deezer_id"], attrs["spotify_id"])

        if not any(store_fields):
            raise ValidationError(
                'Must be set only one of these fields: ["itunes_id", "deezer_id", "spotify_id"]'
            )
        return attrs


class NotificationSerializer(serializers.Serializer):
    user_name = serializers.CharField()
    releases = ReleaseSerializer(many=True)


class ConnectionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Connection
        fields = ["id", "user_name", "provider", "settings"]
        extra_kwargs = {
            "user_name": {"write_only": True},
        }
