# Generated by Django 5.1 on 2024-08-11 16:57

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('pages', '0004_router'),
    ]

    operations = [
        migrations.AddField(
            model_name='router',
            name='name',
            field=models.CharField(default='File', max_length=255),
            preserve_default=False,
        ),
    ]