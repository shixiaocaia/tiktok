#!/bin/sh

ulimit -c unlimited

SERVER_NAME=$1

echo "SERVER_NAME=${SERVER_NAME}"

array=("gatewaysvr" "commentsvr" "favoritesvr" "relationsvr" "usersvr" "videosvr" "all")

run()
{
    if [[ ! "${array[@]}"  =~ ${SERVER_NAME} ]]; then
        echo "server name is not correct"
        exit 0
    elif [ ${SERVER_NAME} = "gatewaysvr" ];then
        cd ../cmd/gatewaysvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "commentsvr" ];then
        cd ../cmd/commentsvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "favoritesvr" ];then
        cd ../cmd/favoritesvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "relationsvr" ];then
        cd ../cmd/relationsvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "usersvr" ];then
        cd ../cmd/usersvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "videosvr" ];then
        cd ../cmd/videosvr/script
        ./build.sh
    elif [ ${SERVER_NAME} = "all" ];then
        cd ../cmd/usersvr/script
        ./build.sh
        cd ../../videosvr/script
        ./build.sh
        cd ../../favoritesvr/script
        ./build.sh
        cd ../../relationsvr/script
        ./build.sh
        cd ../../commentsvr/script
        ./build.sh
        cd ../../gatewaysvr/script
        ./build.sh
    fi
}


usage()
{
    echo "Usage: ./server/sh [start|stop|restart] [gatewaysvr|commentsvr|favoritesvr|relationsvr|usersvr|videosvr|all|...]"
}

if [ $# -lt 1 ];then
    usage
    exit
fi

run