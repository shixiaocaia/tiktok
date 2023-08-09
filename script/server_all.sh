#!/bin/sh

ulimit -c unlimited

ACTION=$1
SERVER_NAME=$2

echo "ACTION=${ACTION}"
echo "SERVER_NAME=${SERVER_NAME}"

array1=("start" "stop")

array2=("gatewaysvr" "commentsvr" "favoritesvr" "relationsvr" "usersvr" "videosvr" "all")

run()
{
    if [[ ! "${array1[@]}"  =~ ${ACTION} ]]; then
        echo "param1 should be start or stop"
        exit 0
    else
        if [[ ! "${array2[@]}"  =~ ${SERVER_NAME} ]]; then
            echo "server name is not correct"
            exit 0
        elif [ ${SERVER_NAME} = "gatewaysvr" ];then
            cd ../server/gatewaysvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "commentsvr" ];then
            cd ../server/commentsvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "favoritesvr" ];then
            cd ../server/favoritesvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "relationsvr" ];then
            cd ../server/relationsvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "usersvr" ];then
            cd ../server/usersvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "videosvr" ];then
            cd ../server/videosvr/script
            ./server.sh ${ACTION}
        elif [ ${SERVER_NAME} = "all" ];then
            SCRIPT_PATH=`pwd`
            echo "SCRIPT_PATH====${SCRIPT_PATH}"
            cd ../cmd/usersvr/script
            ./server.sh ${ACTION}
            cd ../../videosvr/script
            ./server.sh ${ACTION}
            cd ../../favoritesvr/script
            ./server.sh ${ACTION}
            cd ../../relationsvr/script
            ./server.sh ${ACTION}
            cd ../../commentsvr/script
            ./server.sh ${ACTION}
            cd ../../gatewaysvr/script
            ./server.sh ${ACTION}
        fi
    fi
}



usage()
{
    echo "Usage: ./server/sh [start|stop|restart] [gatewaysvr|commentsvr|favoritesvr|relationsvr|usersvr|videosvr|all|...]"
}

if [ $# -lt 2 ];then
    usage
    exit
fi

run


