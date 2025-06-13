import os
import requests
from django.db import models
from django.conf import settings
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from Crypto.Util.Padding import pad, unpad
from users.models import MCUser
from django.core.files.base import ContentFile

class EncryptedFile(models.Model):
    name = models.CharField(max_length=255)
    comment = models.TextField(max_length=1025)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    user = models.ForeignKey(MCUser, on_delete=models.CASCADE)
    private = models.BooleanField(default=False)
    encryption_key = models.BinaryField(max_length=32)
    iv = models.BinaryField(max_length=16)

    def encrypt_file(self, original_file):
        key = get_random_bytes(32)
        iv = get_random_bytes(16)
        cipher = AES.new(key, AES.MODE_CBC, iv)

        with open(original_file.path, 'rb') as f_in:
            plaintext = f_in.read()
            padded_plaintext = pad(plaintext, AES.block_size)
            ciphertext = cipher.encrypt(padded_plaintext)

        self.encryption_key = key
        self.iv = iv
        self.save()

        # Upload to remote server
        self.upload_to_remote(ciphertext, original_file.name + ".enc")

    def upload_to_remote(self, ciphertext, filename):
        """Uploads the encrypted file to the remote Apache server."""
        remote_url = f"http://127.0.1.1:8080/upload"  # Replace <REMOTE_IP>
        files = {'file': (filename, ciphertext, 'application/octet-stream')} # change content type if needed

        try:
            response = requests.post(remote_url, files=files)
            response.raise_for_status()

        except requests.exceptions.RequestException as e:
            print(f"Error uploading file: {e}")

    def decrypt_file(self, output_path):
        """Decrypts the encrypted file from the remote server and saves it locally."""
        remote_url = f"http://127.0.1.1:8080/{os.path.basename(self.name)}.enc"
        try:
            response = requests.get(remote_url)
            response.raise_for_status()
            ciphertext = response.content
            cipher = AES.new(self.encryption_key, AES.MODE_CBC, self.iv)
            padded_plaintext = cipher.decrypt(ciphertext)
            plaintext = unpad(padded_plaintext, AES.block_size)

            with open(output_path, 'wb') as f_out:
                f_out.write(plaintext)

        except requests.exceptions.RequestException as e:
            print(f"Error downloading or decrypting file: {e}")

    def delete(self, *args, **kwargs):
        """Deletes the encrypted file from the remote server."""
        remote_url = f"http://127.0.1.1:8080/{os.path.basename(self.name)}.enc" # replace <REMOTE_IP>

        try:
            response = requests.delete(remote_url) # Apache server must be configured to allow DELETE requests.
            response.raise_for_status()

        except requests.exceptions.RequestException as e:
            print(f"Error deleting file from remote server: {e}")

        super().delete(*args, **kwargs)