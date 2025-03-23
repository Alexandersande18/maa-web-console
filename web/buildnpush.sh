#!/bin/bash
# 方式1：使用 set -e 全局设置
set -eo pipefail

npm run build
docker build -t maasst .
docker tag maasst registry.cn-chengdu.aliyuncs.com/example/maa-web
docker push registry.cn-chengdu.aliyuncs.com/example/maa-web