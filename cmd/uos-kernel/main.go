package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uos-projects/uos-kernel/kernel"
)

func main() {
	// 解析命令行参数
	typesystemPath := flag.String("typesystem", "typesystem.yaml", "Path to typesystem YAML file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建资源内核（自动创建 Actor System）
	// 使用简化方式：Kernel 内部管理 System 生命周期
	k := kernel.NewKernelWithContext(ctx)
	defer k.Shutdown()

	// 3. 加载类型系统定义
	if err := k.LoadTypeSystem(*typesystemPath); err != nil {
		log.Fatalf("Failed to load type system from %s: %v", *typesystemPath, err)
	}

	log.Printf("Type system loaded successfully from %s", *typesystemPath)

	// 4. 打印已加载的类型
	types := k.GetTypeRegistry().List()
	log.Printf("Loaded %d resource types: %v", len(types), types)

	// 5. 设置信号处理，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 6. 示例：创建一个资源（可选）
	// 这里可以添加初始化逻辑，例如从配置文件加载初始资源

	fmt.Println("\nUOS Kernel started successfully!")
	fmt.Println("Press Ctrl+C to shutdown...")

	// 等待信号
	<-sigChan
	fmt.Println("\nShutting down...")
}
