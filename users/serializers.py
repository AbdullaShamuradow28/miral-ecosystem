from rest_framework import serializers
from .models import User, UserProfile, ChatUser

class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ['mid', 'email', 'password']  # Оставляем только email и password

class UserProfileSerializer(serializers.ModelSerializer):
    class Meta:
        model = UserProfile  # Используем UserProfile
        fields = ['mid', 'first_name', 'last_name', 'about_me', 'profile_picture_url', 'nickname']

from rest_framework import serializers
from .models import ChatUser, MiralChatAccountProfile


class ChatUserSerializer(serializers.ModelSerializer):
    class Meta:
        model = ChatUser
        fields = ['mid', 'name', 'email', 'password']


class MiralChatAccountProfileSerializer(serializers.ModelSerializer):
    class Meta:
        model = MiralChatAccountProfile
        fields = ['user', 'about_me', 'profile_picture_url', 'nickname', 'username']

from .models import MCUser

import uuid
import base64
from django.core.files.base import ContentFile
from rest_framework import viewsets, status, serializers
from rest_framework.decorators import action
from rest_framework.response import Response
from .models import MCUser
import logging

logger = logging.getLogger(__name__)

class McUserSerializer(serializers.ModelSerializer):
    profile_picture = serializers.ImageField(max_length=None, use_url=True, allow_null=True, required=False) #added

    class Meta:
        model = MCUser
        fields = ['mid', 'name', 'about', 'profile_picture', 'is_premium']