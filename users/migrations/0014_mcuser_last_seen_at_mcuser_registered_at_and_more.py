# Generated by Django 5.1.1 on 2025-03-29 09:48

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('users', '0013_mcuser_is_premium'),
    ]

    operations = [
        migrations.AddField(
            model_name='mcuser',
            name='last_seen_at',
            field=models.DateTimeField(blank=True, null=True),
        ),
        migrations.AddField(
            model_name='mcuser',
            name='registered_at',
            field=models.DateTimeField(blank=True, null=True),
        ),
        migrations.AddField(
            model_name='mcuser',
            name='unregistered_at',
            field=models.DateTimeField(blank=True, null=True),
        ),
    ]
