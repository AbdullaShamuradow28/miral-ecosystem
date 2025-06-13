from django.db import models
import random

class User(models.Model):
    mid = models.CharField(max_length=6, unique=True, blank=True, null=True)
    email = models.EmailField(unique=True)
    password = models.CharField(max_length=128)

    def save(self, *args, **kwargs):
        if not self.mid:
            self.mid = self.generate_mid()
        super().save(*args, **kwargs)

    def generate_mid(self):
        """Генерация уникального 6-значного идентификатора."""
        while True:
            mid = str(random.randint(100000, 999999))
            if not User.objects.filter(mid=mid).exists():
                return mid

    def __str__(self):
        return self.email

class UserProfile(models.Model):
    # user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='profile', null=True, blank=True)
    mid = models.CharField(max_length=6, unique=True, blank=True, null=True)
    first_name = models.CharField(max_length=30)
    last_name = models.CharField(max_length=30)
    about_me = models.CharField(max_length=255, blank=True, null=True)
    profile_picture_url = models.FileField(upload_to="pfps", default="pfps/avatar.png", blank=True, null=True)
    nickname = models.CharField(max_length=30, blank=True, null=True)

    def __str__(self):
        return f"{self.first_name} {self.last_name} ({self.user.mid})"

class MiralNewsAccountProfile(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='news_profile', null=True, blank=True)
    about_me = models.CharField(max_length=255)
    profile_picture_url = models.FileField(upload_to="pfps", default="pfps/avatar.png")


class ChatUser(models.Model):
    mid = models.CharField(max_length=6, unique=True, blank=True, null=True)
    name = models.CharField(max_length=255)
    email = models.CharField(max_length=255)
    password = models.CharField(max_length=255)

    def save(self, *args, **kwargs):
        if not self.mid:
            self.mid = self.generate_mid()
        super().save(*args, **kwargs)

    def generate_mid(self):
        """Генерация уникального 6-значного идентификатора."""
        while True:
            mid = str(random.randint(100000, 999999))
            if not User.objects.filter(mid=mid).exists():
                return mid
    


class MiralChatAccountProfile(models.Model):
    user = models.OneToOneField(ChatUser, on_delete=models.CASCADE, related_name='chat_profile', null=True, blank=True)
    about_me = models.CharField(max_length=255)
    profile_picture_url = models.FileField(upload_to="pfps", default="pfps/avatar.png")
    nickname = models.CharField(max_length=500)
    username = models.CharField(max_length=500)


# MAIN MIRAL CHAT USER MODEL (CURRENT<ACTIVE>)
class MCUser(models.Model):
    mid = models.CharField(max_length=6, unique=True, blank=True, null=True)
    name = models.CharField(max_length=255)
    about = models.CharField(max_length=255, default="")
    profile_picture = models.ImageField(upload_to="pictures", blank=True, null=True)
    is_premium = models.BooleanField(default=False, null=True, blank=True)
    last_seen_at = models.DateTimeField(null=True, blank=True)
    registered_at = models.DateTimeField(null=True, blank=True)
    unregistered_at = models.DateTimeField(null=True, blank=True)