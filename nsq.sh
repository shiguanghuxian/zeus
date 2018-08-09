#!/bin/bash

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE

mkdir -p var

function check_pid() {
    if [ -f var/$1.pid ];then
        pid=`cat var/$1.pid`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function start() {
    nsqlookupd
    nsqd
    nsqadmin
}

function stop() {
    stopnsqadmin
    sleep 1
    stopnsqlookupd
    sleep 1
    stopnsqd
}

function restart(){
    stop
    sleep 1
    start
}

function nsqlookupd(){
    $(check_pid nsqlookupd);
    running=$?
    if [ $running -gt 0 ];then
        echo -n "nsqlookupd now is running already, pid="
        cat var/nsqlookupd.pid
        return 1
    fi

    nohup ./nsqlookupd &> var/nsqlookupd.log &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > var/nsqlookupd.pid
        echo "nsqlookupd started..., pid=$!"
    else
        echo "nsqlookupd failed to start."
        return 1
    fi
}

function stopnsqlookupd() {
    pid=`cat var/nsqlookupd.pid`
    if [ "$pid" != "" ];then
        kill $pid
        rm -f var/nsqlookupd.pid
        echo "nsqlookupd stoped..."
    else
        echo "stop nsqlookupd error"
    fi
}

function nsqd(){
    $(check_pid nsqd);
    running=$?
    if [ $running -gt 0 ];then
        echo -n "nsqd now is running already, pid="
        cat var/nsqd.pid
        return 1
    fi

    nohup ./nsqd --lookupd-tcp-address=127.0.0.1:4160 &> var/nsqd.log &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > var/nsqd.pid
        echo "nsqd started..., pid=$!"
    else
        echo "nsqd failed to start."
        return 1
    fi
}

function stopnsqd() {
    pid=`cat var/nsqd.pid`
    if [ "$pid" != "" ];then
        kill $pid
        rm -f var/nsqd.pid
        echo "nsqd stoped..."
    else
        echo "stop nsqd error"
    fi
}

function nsqadmin(){
    $(check_pid nsqadmin);
    running=$?
    if [ $running -gt 0 ];then
        echo -n "nsqadmin now is running already, pid="
        cat var/nsqadmin.pid
        return 1
    fi

    nohup ./nsqadmin --lookupd-http-address=127.0.0.1:4161 &> var/nsqadmin.log &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > var/nsqadmin.pid
        echo "nsqadmin started..., pid=$!"
    else
        echo "nsqadmin failed to start."
        return 1
    fi
}

function stopnsqadmin() {
    pid=`cat var/nsqadmin.pid`
    if [ "$pid" != "" ];then
        kill $pid
        rm -f var/nsqadmin.pid
        echo "nsqadmin stoped..."
    else
        echo "stop nsqadmin error"
    fi
}

function help() {
    echo "$0 pack|start|stop|restart"
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "pack" ];then
    pack
else
    help
fi
