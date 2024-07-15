#!/bin/bash

# 检查目录是否存在
if [ -d "./publish" ]; then
    # 使用 find 命令列出目录中的所有文件并整理成列表
    file_list=$(find ./publish -maxdepth 1 -type f)
    
    # 打印文件列表
    echo "文件列表："
    for file in $file_list
    do
        echo "$file"
    done
else
    echo "./publish 目录不存在"
fi
