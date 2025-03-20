function exampleRestart() {
  echo "=================>"
  # 先尝试正常终止
  killall qor5example 2>/dev/null || true
  
  # 查找并终止进程，使用多种方法保证可靠性
  PID=$(ps -ef | grep "/tmp/qor5example" | grep -v grep | awk '{print $2}')
  if [ ! -z "$PID" ]; then
    echo "发现进程 $PID，正在终止..."
    kill -15 $PID 2>/dev/null || true
    sleep 0.5
    # 如果进程还在运行，使用强制终止
    if ps -p $PID > /dev/null; then
      echo "强制终止进程 $PID..."
      kill -9 $PID
    fi
  fi
  
  source dev_env
#  export DEV_PRESETS=1
  go build -o /tmp/qor5example main.go && /tmp/qor5example
}

export -f exampleRestart

find . -name "*.go" | entr -r bash -c "exampleRestart"
