from .models import MCUser
from .serializers import McUserSerializer
from rest_framework import viewsets, status
from rest_framework.response import Response
from .models import UserProfile
from .serializers import UserProfileSerializer
from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.decorators import action
from rest_framework import serializers
from .models import ChatUser 
from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.decorators import action
from .models import User, MiralNewsAccountProfile, UserProfile
from .serializers import UserSerializer, UserProfileSerializer, ChatUserSerializer
from .models import ChatUser, MiralChatAccountProfile
from .serializers import MiralChatAccountProfileSerializer

class UserViewSet(viewsets.ModelViewSet):
    queryset = User.objects.all()
    serializer_class = UserSerializer

    @action(detail=False, methods=['get'], url_path='(?P<mid>\w+)')
    def retrieve_by_mid(self, request, mid=None):
        try:
            user = User.objects.get(mid=mid)
            serializer = self.get_serializer(user)
            return Response(serializer.data)
        except User.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)


class MiralProfileViewSet(viewsets.ModelViewSet):
    queryset = UserProfile.objects.all()
    serializer_class = UserProfileSerializer

    @action(detail=False, methods=['get', 'post', 'put', 'delete'], url_path='get_by_mid/(?P<mid>[^/.]+)')
    def retrieve_by_mid(self, request, mid=None):
        try:
            profile = UserProfile.objects.get(mid=mid)
            serializer = self.get_serializer(profile)
            return Response(serializer.data)
        except UserProfile.DoesNotExist:
            return Response({'detail': 'Profile not found.'}, status=status.HTTP_404_NOT_FOUND)

    def create(self, request, *args, **kwargs):
        mid = request.data.get('mid')
        try:
            profile = UserProfile.objects.get(mid=mid)
            serializer = self.get_serializer(profile, data=request.data, partial=True)
            if serializer.is_valid():
                serializer.save()
                return Response(serializer.data, status=status.HTTP_200_OK)
            else:
                return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
        except UserProfile.DoesNotExist:
            serializer = self.get_serializer(data=request.data)
            if serializer.is_valid():
                serializer.save()
                return Response(serializer.data, status=status.HTTP_201_CREATED)
            else:
                return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


class ChatUserSerializer(serializers.ModelSerializer):
    class Meta:
        model = ChatUser
        fields = ['mid', 'name', 'email', 'password']
class ChatUserViewSet(viewsets.ModelViewSet):
    queryset = ChatUser.objects.all()
    serializer_class = ChatUserSerializer

    @action(detail=False, methods=['post'], url_path='login')
    def login(self, request):
        mid = request.data.get('email')
        password = request.data.get('password')

        if not mid or not password:
            return Response({'detail': 'email and password are required.'}, status=status.HTTP_400_BAD_REQUEST)

        try:
            user = ChatUser.objects.get(email=mid)
            if user.password == password:  # Basic password check (replace with hashing)
                serializer = self.get_serializer(user)
                return Response(serializer.data)
            else:
                return Response({'detail': 'Invalid password.'}, status=status.HTTP_401_UNAUTHORIZED)
        except ChatUser.DoesNotExist:
            return Response({'detail': 'User not found. Please create an account.'}, status=status.HTTP_404_NOT_FOUND)

    @action(detail=False, methods=['get'], url_path='get_by_mid/(?P<mid>[^/.]+)')
    def retrieve_by_mid(self, request, mid=None):
        try:
            user = ChatUser.objects.get(mid=mid)
            serializer = self.get_serializer(user)
            return Response(serializer.data)
        except ChatUser.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)

    @action(detail=False, methods=['get'], url_path='get_by_email/(?P<email>[^/.]+)/?')
    def retrieve_by_email(self, request, email=None):
        try:
            user = ChatUser.objects.get(email=email)
            serializer = self.get_serializer(user)
            return Response(serializer.data)
        except ChatUser.DoesNotExist:
            all_users = ChatUser.objects.all()
            serializer = self.get_serializer(all_users, many=True)
            return Response(serializer.data)

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
            user = serializer.save()  # Save the user and get the instance
            return Response({'mid': user.mid}, status=status.HTTP_201_CREATED) # Return mid
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

class MiralChatAccountProfileViewSet(viewsets.ModelViewSet):
    queryset = MiralChatAccountProfile.objects.all()
    serializer_class = MiralChatAccountProfileSerializer

    @action(detail=False, methods=['get'], url_path='get_by_mid/(?P<mid>[^/.]+)')
    def retrieve_by_mid(self, request, mid=None):
        try:
            user = ChatUser.objects.get(mid=mid)
            profile = MiralChatAccountProfile.objects.get(user=user)
            serializer = self.get_serializer(profile)
            return Response(serializer.data)
        except ChatUser.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)
        except MiralChatAccountProfile.DoesNotExist:
            return Response({'detail': 'Profile not found.'}, status=status.HTTP_404_NOT_FOUND)

    def create(self, request, *args, **kwargs):
        mid = request.data.get('mid')
        try:
            user = ChatUser.objects.get(mid=mid)
            try:
                profile = MiralChatAccountProfile.objects.get(user=user)
                serializer = self.get_serializer(profile, data=request.data, partial=True)
                if serializer.is_valid():
                    serializer.save()
                    return Response(serializer.data, status=status.HTTP_200_OK)
                else:
                    return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
            except MiralChatAccountProfile.DoesNotExist:
                serializer = self.get_serializer(data=request.data)
                if serializer.is_valid():
                    serializer.save(user=user)
                    return Response(serializer.data, status=status.HTTP_201_CREATED)
                else:
                    return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

        except ChatUser.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)


class ChatAccountProfilesViewSet(viewsets.ModelViewSet):
    queryset = MCUser.objects.all()
    serializer_class = McUserSerializer

    @action(detail=False, methods=['get'], url_path='get_by_mid/(?P<mid>[^/.]+)')
    def retrieve_by_mid(self, request, mid=None):
        try:
            user = MCUser.objects.get(mid=mid)
            serializer = self.get_serializer(user)
            return Response(serializer.data)
        except MCUser.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)

    @action(detail=False, methods=['put'], url_path='update') #correct url
    def update_by_mid(self, request): #remove pk from arguments
        try:
            mid = request.data.get('mid')
            user = MCUser.objects.get(mid=mid)
            serializer = self.get_serializer(user, data=request.data, partial=True)
            if serializer.is_valid():
                serializer.save()
                return Response(serializer.data)
            return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
        except MCUser.DoesNotExist:
            return Response({'detail': 'User not found.'}, status=status.HTTP_404_NOT_FOUND)

    @action(detail=False, methods=['post'], url_path='new')
    def new(self, request):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
            try:
                serializer.save()
                return Response(serializer.data, status=status.HTTP_201_CREATED)
            except Exception as e:
                logger.error(f"Error saving new MCUser: {e}")
                return Response({'detail': 'Error saving user.'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)