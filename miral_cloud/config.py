import json
import uuid

PROJECTS_LIMIT = 1
STORAGE_MAX_SIZE = 200


class Project:
    def __init__(self, name, owner):
        self.name = name
        self.owner = owner
        self.storage = 0
        self.collections = {}

    def to_dict(self):
        return {
            'name': self.name,
            'owner': self.owner,
            'storage': self.storage,
            'collections': list(self.collections.keys()),
        }

    def to_json(self):
        return json.dumps(self.to_dict())

    def __str__(self):
        return f'Project({self.name}, {self.owner}, {self.storage} GB)'

class StorageError(Exception):
    def __init__(self, message, code, reason):
        super().__init__(message)
        self.message = message
        self.code = code
        self.reason = reason

    def __str__(self):
        return f'StorageError({self.message}, {self.code}, {self.reason})'

    def to_dict(self):
        return {
            'message': self.message,
            'code': self.code,
            'reason': self.reason,
        }

    def to_json(self):
        return json.dumps(self.to_dict())

class MiralCloud:
    def __init__(self):
        self.projects = []
    
    def storage(self):
        def create_bucket(self, bucket_name, bucket_id, bucket_secret):
            # Implementation of bucket creation logic goes here
            return bucket_id, bucket_secret
        def buckets(self, bucket_name, bucket_id):
            # Implementation of bucket retrieval logic goes here
            def get_all(self):
                return "all"
            def get(self, bucket_id):
                if bucket_id == bucket_name:
                    return "bucket"
                else:
                    raise StorageError('Bucket not found', 404, 'Bucket not found')
            def delete(self, bucket_name):
                # Implementation of bucket deletion logic goes here
                return f'Bucket {bucket_name} deleted'
        def upload(self, bucket_name, file):
            # Implementation of file upload logic goes here
            return f'File {file} uploaded to bucket {bucket_name}'
        def download(self, bucket_name, file, output_path):
            # Implementation of file download logic goes here
            return f'File {file} downloaded from bucket {bucket_name} to {output_path}'
        def delete(self, bucket_name, file):
            # Implementation of file deletion logic goes here
            return f'File {file} deleted from bucket {bucket_name}'

    def create_project(self, name, owner):
        if len(self.projects) >= PROJECTS_LIMIT:
            raise StorageError('Storage limit reached', 413, 'Too many projects')

        project = Project(name, owner)
        self.projects.append(project)
        return project

    def get_project(self, name):
        for project in self.projects:
            if project.name == name:
                return project

        raise StorageError('Project not found', 404, 'Project not found')

    def delete_project(self, name):
        project = self.get_project(name)
        self.projects.remove(project)

    def collection(self, project_name, collection_name):
        project = self.get_project(project_name)

        if collection_name not in project.collections:
            project.collections[collection_name] = {}

        def get(doc_id=None):
            if doc_id:
                if doc_id in project.collections[collection_name]:
                    return project.collections[collection_name][doc_id]
                else:
                    raise StorageError('Document not found', 404, 'Document not found')
            else:
                return list(project.collections[collection_name].values())

        def add(data):
            doc_id = str(uuid.uuid4())
            project.collections[collection_name][doc_id] = data
            return doc_id

        def update(doc_id, data):
            if doc_id in project.collections[collection_name]:
                project.collections[collection_name][doc_id].update(data)
            else:
                raise StorageError('Document not found', 404, 'Document not found')

        def destroy(doc_id):
            if doc_id in project.collections[collection_name]:
                del project.collections[collection_name][doc_id]
            else:
                raise StorageError('Document not found', 404, 'Document not found')

        return {
            'get': get,
            'add': add,
            'update': update,
            'destroy': destroy,
        }

"""Example Usage"""
miral_cloud = MiralCloud()
project = miral_cloud.create_project('MyProject', 'User1')
collection = miral_cloud.collection('MyProject', 'users')

doc_id = collection['add']({'name': 'Alice', 'age': 25})
print(f"Document added with ID: {doc_id}")

doc = collection['get'](doc_id)
print(f"Retrieved document: {doc}")

collection['update'](doc_id, {'age': 31})
updated_doc = collection['get'](doc_id)
print(f"Updated document: {updated_doc}")

all_docs = collection['get']()
print(f"All documents: {all_docs}")

collection['destroy'](doc_id)
print("Document deleted")
try:
    collection['get'](doc_id)
except StorageError as e:
    print(e)

input("Enter smth to continue :> ")