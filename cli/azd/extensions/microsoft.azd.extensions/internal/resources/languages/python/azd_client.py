import grpc  
from grpc_interceptor import ClientCallDetails, ClientInterceptor  
  
class AuthHeaderInterceptor(ClientInterceptor):  
    def __init__(self, access_token):  
        self.access_token = access_token  
  
    def intercept(self, method, request_or_iterator, call_details):  
        metadata = call_details.metadata or []  
        metadata.append(('authorization', self.access_token))  
        new_call_details = ClientCallDetails(  
            call_details.method, call_details.timeout, metadata,  
            call_details.credentials, call_details.wait_for_ready)  
        return method(request_or_iterator, new_call_details)  
  
class AzdClient:  
    def __init__(self, server_address, access_token):  
        if not server_address.startswith("http"):  
            server_address = "http://" + server_address  
        self.channel = grpc.insecure_channel(server_address)  
        self.channel = grpc.intercept_channel(self.channel, AuthHeaderInterceptor(access_token))  
        # Initialize gRPC service clients here  
        # self.compose = ComposeServiceStub(self.channel)  
        # self.deployment = DeploymentServiceStub(self.channel)  
        # self.environment = EnvironmentServiceStub(self.channel)  
        # self.events = EventServiceStub(self.channel)  
        # self.project = ProjectServiceStub(self.channel)  
        # self.prompt = PromptServiceStub(self.channel)  
        # self.user_config = UserConfigServiceStub(self.channel)  
        # self.workflow = WorkflowServiceStub(self.channel)  
  
    def close(self):  
        self.channel.close()  