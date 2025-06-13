import os
import requests
from django.db import models
from users.models import MCUser

class Group(models.Model):
    name = models.CharField(max_length=255)
    creator = models.ForeignKey(MCUser, on_delete=models.CASCADE)
    
    def __str__(self):
        return f'{self.name} ({self.creator.mid})'

class File(models.Model):
    name = models.CharField(max_length=255)
    comment = models.TextField(max_length=1025)
    file = models.FileField(upload_to='files/')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    user = models.ForeignKey(MCUser, on_delete=models.CASCADE)
    private = models.BooleanField(default=False)
    group = models.ForeignKey(Group, on_delete=models.CASCADE, related_name='files', null=True)

    def __str__(self):
        return f'{self.name} ({self.user.mid}) --> C: {self.created_at} U: {self.updated_at} <--'