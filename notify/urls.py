"""URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.0/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.urls import path

from base.views import ConnectionsViewSet, new_releases

connection_list = ConnectionsViewSet.as_view({"get": "list", "post": "create"})
connection_detail = ConnectionsViewSet.as_view(
    {"get": "retrieve", "put": "update", "patch": "partial_update", "delete": "destroy"}
)

urlpatterns = [
    path("releases", new_releases, name="releases"),
    path("connections", connection_list, name="connection-list"),
    path("connections/<int:pk>", connection_detail, name="connection-detail"),
]
