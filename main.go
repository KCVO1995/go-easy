package main

import (
	"github.com/zserge/lorca"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var ui lorca.UI
	// 打开一个 chrome 窗口
	ui, _ = lorca.New("https://www.zhihu.com", "", 800, 600, "--disable-sync", "--disable-translate")
	// 创建一个频道连接系统信号
	chSignal := make(chan os.Signal, 1)
	// 监听系统信号，如果是 interrupt 或者 termination
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	// 等待 ui 结束，或者 chSignal 收到信号
	select {
	case <-ui.Done():
	case <-chSignal:
	}
	// 等待完毕，关闭 UI，结束主线程
	ui.Close()
}
