from rest_framework import viewsets, status
from rest_framework.response import Response
from .models import Email
from .serializers import EmailSerializer
from users.models import User  

class EmailViewSet(viewsets.ModelViewSet):
    queryset = Email.objects.all()
    serializer_class = EmailSerializer

    def create(self, request, *args, **kwargs):
        sender_email = request.data.get("sender")
        recipient_email = request.data.get("recipient")

        if sender_email is None or recipient_email is None:
            return Response({"detail": "Sender and recipient emails must be provided."}, status=status.HTTP_400_BAD_REQUEST)

        try:
            sender = User.objects.get(email=sender_email)  # Optional validation
        except User.DoesNotExist:
            return Response({"detail": f"Sender with email '{sender_email}' does not exist."}, status=status.HTTP_404_NOT_FOUND)

        try:
            recipient = User.objects.get(email=recipient_email)  # Optional validation
        except User.DoesNotExist:
            return Response({"detail": f"Recipient with email '{recipient_email}' does not exist."}, status=status.HTTP_404_NOT_FOUND)

        email_data = {
            'sender': sender_email,  # Now just passing the strings directly
            'recipient': recipient_email,
            'subject': request.data.get('subject'),
            'body': request.data.get('body')
        }

        serializer = self.get_serializer(data=email_data)
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status=status.HTTP_201_CREATED)

        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)