from django.db import models

class Email(models.Model):
    sender = models.EmailField(max_length=255)  # Change to EmailField for better validation
    recipient = models.EmailField(max_length=255)  # Change to EmailField for better validation
    subject = models.CharField(max_length=255)
    body = models.TextField()

    def __str__(self):
        return f"Email from {self.sender} to {self.recipient}"