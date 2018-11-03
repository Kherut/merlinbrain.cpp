#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>
#include <cstdlib>

bool port_available(int portno) {
    int sockfd;
    socklen_t clilen;
    struct sockaddr_in serv_addr, cli_addr;
    bool available = true;

    sockfd = socket(AF_INET, SOCK_STREAM, 0);

    bzero((char *) &serv_addr, sizeof(serv_addr));

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(portno);

    if (bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0)
        available = false;

    close(sockfd);

    return available;
}

int main(int argc, char **argv) {
    int i = atoi(argv[1]);
    int max = atoi(argv[2]);

    while(i <= max && !port_available(i))
        i++;

    printf("%d", i);

    return 0;
}