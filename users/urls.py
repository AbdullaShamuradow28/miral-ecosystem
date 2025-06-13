from django.urls import path, include
from rest_framework.routers import DefaultRouter
from .views import UserViewSet, MiralProfileViewSet, ChatUserViewSet, ChatAccountProfilesViewSet
from miral_cloud.views import FileViewSet
from emails.views import EmailViewSet

router = DefaultRouter()
router.register(r'users', UserViewSet, basename='user')
router.register(r'profiles', MiralProfileViewSet)
router.register(r'chatusers', ChatUserViewSet)
router.register(r'chat/account/profile', ChatAccountProfilesViewSet)
router.register(r'storage', FileViewSet)
router.register(r'emails', EmailViewSet)

urlpatterns = [
    path('', include(router.urls)),
]