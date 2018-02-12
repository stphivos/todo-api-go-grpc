import grpc
import time
import todos_pb2
import todos_pb2_grpc

from locust import Locust, TaskSet, task, events


class GrpcClient:
    def __init__(self, host):
        channel = grpc.insecure_channel(host)
        self.stub = todos_pb2_grpc.TodosStub(channel)

    def __getattr__(self, name):
        func = self.stub.__getattribute__(name)

        def wrapper(*args, **kwargs):
            start_time = time.time()
            try:
                response = func(*args, **kwargs).SerializeToString()
            except Exception as e:
                total_time = int((time.time() - start_time) * 1000)
                events.request_failure.fire(request_type="grpc", name=name, response_time=total_time, exception=e)
                print e
            else:
                total_time = int((time.time() - start_time) * 1000)
                events.request_success.fire(request_type="grpc", name=name, response_time=total_time, response_length=len(response))

        return wrapper


class GrpcLocust(Locust):
    def __init__(self, *args, **kwargs):
        super(GrpcLocust, self).__init__(*args, **kwargs)
        self.client = GrpcClient(self.host)


class ApiUser(GrpcLocust):
    min_wait = 100
    max_wait = 1000

    class task_set(TaskSet):
        @task()
        def get_todos(self):
            self.client.GetTodos(todos_pb2.Request(token='xyz'))
