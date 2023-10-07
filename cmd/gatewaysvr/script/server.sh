#!/bin/sh

ulimit -c unlimited

SERVER_NAME="gatewaysvr"
SCRIPT_PATH=`pwd`
SERVER_PATH=`echo ${SCRIPT_PATH%/*}`
SERVER_BIN_PATH="${SCRIPT_PATH}/../bin"
echo "server_name=${SERVER_NAME}"


is_running()
{   
    proc_num=$(ps -ef | grep -w "${SERVER_NAME}" | grep -w "${SERVER_BIN_PATH}" | grep -v grep | wc -l)
    if [ ${proc_num} -gt 0 ];then
        echo "Server ${SERVER_NAME} has already running!"
        return 1
    else
        return 0
    fi
}

start()
{
    is_running
    if [ $? -eq 0 ]; then
        #remove_prestop_file
        export GOTRACEBACK=crash
        nohup ${SERVER_BIN_PATH}/${SERVER_NAME} > ${SERVER_NAME}.log 2>&1 &
        ret=$?
        if [ $ret -eq 0 ];then
            ps -C "$SERVER_NAME" -o "pid=" > ${SERVER_NAME}.pid
            # ps -ef | grep -w "$SERVER_NAME" | grep -v grep | awk '{print $2}' > ${SERVER_NAME}.pid
            echo "Start server ${SERVER_NAME} OK"
        else
            echo "Start server ${SERVER_NAME} FAILED code "$ret
            #如果不是守护，异常退出的时候脚本也算异常退出
            exit $ret
        fi
    else
        echo "Start server ${SERVER_NAME} FAILED"
    fi
}

stop()
{
    i=3
    stop_flag=0
    while [ $i -gt 0 ]
    do
        is_running
        if [ $? -eq 0 ]; then
            stop_flag=1
            break
        else
            killall ${SERVER_NAME}
            sleep 1
        fi

        ((i=$i-1))
    done

    if [ $stop_flag -eq 1 ];then
        echo "Stop server ${SERVER_NAME} OK"
        exit 0
    fi

    if [ ${stop_flag} -eq 0 ] ; then
        is_running
        if [ $? -eq 0 ]; then
            stop_flag=1
            break
        else
            ps -ef | grep -w "${SERVER_NAME}" | grep -w "${SERVER_BIN_PATH}" | grep -v grep | awk '{print $2}' | xargs kill -9
            sleep 1
        fi
        echo "stop_flag=$stop_flag"
        if [ $stop_flag -eq 1 ];then
            echo "Stop server ${SERVER_NAME} OK"
        else
            echo "Stop server ${SERVER_NAME} FAILED"
        fi
    fi
}

usage()
{
    echo "Usage: ./server/sh [start|stop|restart] [bot1|bot2|bot3|...] [camp1|camp2|camp3|...]"
}

if [ $# -lt 1 ];then
    usage
    exit
fi

if [ "$1" = "start" ];then
    start

elif [ "$1" = "stop" ];then
    stop
else
  usage
fi
