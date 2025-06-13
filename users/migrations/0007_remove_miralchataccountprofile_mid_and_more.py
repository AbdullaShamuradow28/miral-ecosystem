# Generated by Django 5.1.1 on 2025-03-01 15:53

import django.db.models.deletion
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('users', '0006_miralchataccountprofile_remove_user_first_name_and_more'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='miralchataccountprofile',
            name='mid',
        ),
        migrations.RemoveField(
            model_name='miralnewsaccountprofile',
            name='mid',
        ),
        migrations.RemoveField(
            model_name='userprofile',
            name='mid',
        ),
        migrations.AddField(
            model_name='miralchataccountprofile',
            name='user',
            field=models.OneToOneField(blank=True, null=True, on_delete=django.db.models.deletion.CASCADE, related_name='chat_profile', to='users.user'),
        ),
        migrations.AddField(
            model_name='miralnewsaccountprofile',
            name='user',
            field=models.OneToOneField(blank=True, null=True, on_delete=django.db.models.deletion.CASCADE, related_name='news_profile', to='users.user'),
        ),
        migrations.AddField(
            model_name='user',
            name='mid',
            field=models.CharField(blank=True, max_length=6, null=True, unique=True),
        ),
        migrations.AlterField(
            model_name='userprofile',
            name='user',
            field=models.OneToOneField(blank=True, null=True, on_delete=django.db.models.deletion.CASCADE, related_name='profile', to='users.user'),
        ),
    ]
