import requests
from rest_framework import serializers, viewsets, status
from rest_framework.response import Response
from .models import File
from users.models import MCUser
from .serializers import FileSerializer

class FileViewSet(viewsets.ModelViewSet):
    queryset = File.objects.all()
    serializer_class = FileSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)

        # Check if the file is None or missing
        if 'file' not in request.FILES:
            return Response({"file": ["This field is required."]}, status=status.HTTP_400_BAD_REQUEST)

        self.perform_create(serializer)
        headers = self.get_success_headers(serializer.data)
        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)

    def get_queryset(self):
        """Optionally filters the queryset based on query parameters."""
        queryset = File.objects.all()
        private = self.request.query_params.get('private')
        user = self.request.query_params.get('user')

        if private is not None:
            queryset = queryset.filter(private=private)

        if user is not None:
            queryset = queryset.filter(user=user)

        return queryset