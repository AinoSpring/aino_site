# Generated by Django 5.1 on 2024-08-11 19:11

import markdownx.models
from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ("pages", "0005_router_name"),
    ]

    operations = [
        migrations.AlterField(
            model_name="post",
            name="body",
            field=markdownx.models.MarkdownxField(),
        ),
    ]
