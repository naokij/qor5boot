#!/bin/bash
set -e

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查参数
if [ "$#" -lt 1 ]; then
    echo -e "${RED}错误: 请指定环境 (例如: dev)${NC}"
    echo -e "使用方法: $0 <环境>"
    exit 1
fi

ENV=$1
INVENTORY="deploy/$ENV/hosts"

# 检查inventory文件是否存在
if [ ! -f "$INVENTORY" ]; then
    echo -e "${RED}错误: 找不到inventory文件 $INVENTORY${NC}"
    exit 1
fi

# 检查二进制文件是否存在或编译新的
if [ ! -f "qor5boot" ] || [ "$2" == "--rebuild" ]; then
    echo -e "${YELLOW}正在为Linux AMD64平台编译项目...${NC}"
    # 使用交叉编译，指定目标为Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o qor5boot
    if [ $? -ne 0 ]; then
        echo -e "${RED}编译失败，退出部署${NC}"
        exit 1
    fi
    echo -e "${GREEN}编译完成${NC}"
fi

# 检查环境变量文件
if [ ! -f "dev_env" ]; then
    echo -e "${YELLOW}警告: 环境变量文件 'dev_env' 不存在${NC}"
    echo -e "${YELLOW}请从dev_env.example复制并修改${NC}"
    exit 1
fi

# 运行Ansible playbook
echo -e "${GREEN}开始部署到$ENV环境...${NC}"
echo -e "${GREEN}目标平台: Linux AMD64${NC}"
ansible-playbook -i "$INVENTORY" deploy/deploy.yml -v

echo -e "${GREEN}部署完成!${NC}" 