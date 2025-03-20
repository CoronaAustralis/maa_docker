#!/bin/bash

task_list_0=('start' 'infrast' 'recruit' 'fight' 'award')
task_list_1=('start' 'infrast' 'recruit' 'mall' 'fight' 'award')

if [[ $1 -eq 0 ]]; then
    for i in "${task_list_0[@]}"; do
        maa run $i
    done
fi

if [[ $1 -eq 1 ]]; then
    for i in "${task_list_1[@]}"; do
        maa run $i
    done
fi

adb shell am force-stop com.hypergryph.arknights.bilibili
