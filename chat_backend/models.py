from django.db import models
from users.models import MCUser

# Create your models here.
class Message(models.Model):
    pass

class History(models.Model):
    user = models.ForeignKey(MCUser, on_delete=models.CASCADE, related_name='histories')
    content_text = models.TextField(blank=True, null=True)
    content_image = models.ImageField(upload_to='histories/images/', blank=True, null=True)
    content_video = models.FileField(upload_to='histories/videos/', blank=True, null=True)
    content_audio = models.FileField(upload_to='histories/audio/', blank=True, null=True)
    forwarded_post = models.TextField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)
    expires_at = models.DateTimeField()
    viewed_by = models.ManyToManyField(MCUser, related_name='viewed_histories', blank=True)

    def __str__(self):
        return f"История {self.user.name} от {self.created_at}"

    class Meta:
        ordering = ['-created_at']