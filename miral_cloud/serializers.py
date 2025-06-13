from rest_framework import serializers, viewsets, status
from rest_framework.response import Response
from .models import File
from users.models import MCUser
import requests

class FileSerializer(serializers.ModelSerializer):
    class Meta:
        model = File
        fields = ('id', 'name', 'comment', 'file', 'created_at', 'updated_at', 'user', 'private')
        read_only_fields = ('created_at', 'updated_at')