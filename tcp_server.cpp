// Server side C/C++ program to demonstrate Socket programming 
#include <unistd.h> 
#include <stdio.h> 
#include <sys/socket.h> 
#include <stdlib.h> 
#include <netinet/in.h> 
#include <string.h> 

int main(int argc, char const *argv[]) { 
	int server_fd, nsocket, valread; 
    int port = atoi(argv[1]);
	struct sockaddr_in address; 
	int opt = 1; 
	int addrlen = sizeof(address); 
	char buffer[256] = {0}; 
	char msg[256] = {0}; 
	
	//CREATE SOCKET
	if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) { 
		printf("error_socket"); 
		exit(EXIT_FAILURE); 
	} 
	
	//CONFIGURATION
	if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt)))  { 
		printf("error_setsockopt"); 
		exit(EXIT_FAILURE); 
	}

	address.sin_family = AF_INET; 
	address.sin_addr.s_addr = INADDR_ANY; 
	address.sin_port = htons(port); 
	
	//ATTACHING SOCKET
	if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0)  {
		printf("error_bind"); 
		exit(EXIT_FAILURE); 
	} 

	if (listen(server_fd, 3) < 0) { 
		printf("error_listen"); 
		exit(EXIT_FAILURE); 
	}

	if ((nsocket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen)) < 0) {
		printf("error_accept"); 
		exit(EXIT_FAILURE); 
	} 

	while((valread = read(nsocket, buffer, 256)) > 0) {
		printf("%s", buffer);

		fgets(msg, 256, stdin);
		send(nsocket, msg, 256, 0); 
	} 
	
	printf("Hello message sent\n"); 
	return 0; 
} 
