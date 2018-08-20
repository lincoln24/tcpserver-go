/****************************************************************************
Copyright (C) , Future Advanced Technology Research Institute.
File name     : main.c
Description   : XML文件管理客户端模块
Version       : v_1.0
Date          : 2018-05-07
****************************************************************************/
#include <sys/types.h>
#include <sys/socket.h>
#include <strings.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include <sys/wait.h>
#include <sys/time.h>

int main(int argc, char **argv)
{
    int ret = -1, iCount = 0;
    int sfp = -1, i;
    struct sockaddr_in s_add;
    fd_set wfds, rfds;
    struct timeval tv;
    char sendbuf[50000] = {0};
    char recvbuf[256];

    sendbuf[0] = 0x7E;
    sendbuf[1] = 1;
    sendbuf[2] = 1;//设备类型
    sendbuf[3] = 0;//设备编号
    sendbuf[4] = 2;
    sendbuf[5] = 1;//数据类型
    sendbuf[6] = 0x9C;//长度
    sendbuf[7] = 0x40;
    for (i=0;i<40000;i=i+4) {
        sendbuf[8+i]=0;
        sendbuf[9+i]=1;
        sendbuf[10+i]=2;
        sendbuf[11+i]=1;
    }

    for(;;)
    {        
        sfp = socket(AF_INET, SOCK_STREAM, 0);
        if (sfp == -1) {
            printf("socket error : %s", (char *)strerror(errno));

            continue;
        }
        // USER_PRINTF_MSG(PRINTF_LEVEL_INFO,"pid=%d",gettid());
        s_add.sin_family = AF_INET;
        s_add.sin_port = htons(50000);
        s_add.sin_addr.s_addr = inet_addr("39.108.4.211");
        memset(s_add.sin_zero, 0, 8);

        // printf("connect start");
        ret = connect(sfp, (const struct sockaddr *)&s_add, sizeof(struct sockaddr));
        if (ret == -1) {
            printf("connect error : %s", (char *)strerror(errno));
            close(sfp);

            continue;
        }
        // printf( "connect sucess,fd=%d", sfp);
        tv.tv_sec = 10;
        tv.tv_usec = 0;
        FD_ZERO(&wfds);
        FD_ZERO(&rfds);
        FD_SET(sfp, &wfds);
        FD_SET(sfp, &rfds);

        for(;;)
        {
            // printf( "begin send,ret=%d", ret);
            // printf( "sendbuf : %s", sendbuf);
            // if (FD_ISSET(sfp, &wfds))
            // {
            ret = send(sfp, sendbuf, 40008, 0);
            if (ret <= 0)
            {
                printf( "send fail,ret=%d", ret);
                break;
            }
            // }
            // else
            // {
            // printf( "begin send,ret=%d",ret);
            //     close(sfp);
            //     return -1;
            // }

            ret = select(sfp + 1, NULL, &rfds, NULL, &tv);
            if (ret < 0)
            {
                printf( "select fail,fd=%d", sfp);
                break;
            }

            if (FD_ISSET(sfp, &rfds))
            {
                memset(recvbuf, 0, sizeof(recvbuf));
                for (;;)
                {
                    ret = recv(sfp, recvbuf, sizeof(recvbuf), MSG_DONTWAIT);
                    // printf( "read finish,ret=%d",ret);
                    if (ret <= 0)
                    {
                        usleep(50000);//休息50ms再接收
                    }
                    else
                    {
                        break;//数据接收完毕
                    }
                }
            }
            else
            {
                break;
            }
        }
        close(sfp);
    }

    // USER_PRINTF_MSG(PRINTF_LEVEL_DEBUG, "recvbuf=%s", data);
    return 0;
}
