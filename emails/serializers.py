from rest_framework import serializers
from .models import Email

class EmailSerializer(serializers.ModelSerializer):
    sender = serializers.EmailField()  # Directly use EmailField
    recipient = serializers.EmailField()  # Directly use EmailField

    class Meta:
        model = Email
        fields = ['id', 'sender', 'recipient', 'subject', 'body']

    def create(self, validated_data):
        return Email.objects.create(**validated_data)  # Create Email without sender/recipient extraction