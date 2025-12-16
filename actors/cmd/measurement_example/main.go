package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uos-projects/uos-kernel/actors"
	"github.com/uos-projects/uos-kernel/actors/capacities"
	"github.com/uos-projects/uos-kernel/actors/mq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建 MQ 消费者（模拟）
	mqConsumer := mq.NewMockMQConsumer()
	defer mqConsumer.Close()

	// 创建 Actor 系统
	system := actors.NewSystem(ctx)
	defer system.Shutdown()

	// 创建 PowerSystemResourceActor（代表一个发电机）
	generator := actors.NewPowerSystemResourceActor(
		"GEN_001",
		"SynchronousMachine",
		nil,
	)

	// 添加 Control Capacity（设定点控制）
	setPointCapacity := capacities.NewSetPointCapacity("SP_001")
	generator.AddCapacity(setPointCapacity)

	// 添加 Measurement Capacity（功率测量）
	powerMeasurement := capacities.NewAnalogMeasurementCapacity(
		"MEAS_POWER_001",
		"ThreePhaseActivePower",
		mqConsumer,
	)
	generator.AddCapacity(powerMeasurement)

	// 注册到系统
	if err := system.Register(generator); err != nil {
		log.Fatalf("Failed to register actor: %v", err)
	}

	// 启动 Actor（会自动启动所有 Measurement Capacity 的订阅）
	if err := generator.Start(ctx); err != nil {
		log.Fatalf("Failed to start actor: %v", err)
	}

	// 模拟外部系统发布测量数据
	go func() {
		time.Sleep(100 * time.Millisecond) // 等待订阅建立

		// 发布功率测量值
		for i := 0; i < 5; i++ {
			msg := &mq.MQMessage{
				Value:     1000.0 + float64(i)*100.0, // 功率值（MW）
				Timestamp: time.Now(),
				Quality:   mq.QualityGood,
				Source:    "SCADA",
			}

			if err := mqConsumer.Publish("measurement/MEAS_POWER_001", msg); err != nil {
				log.Printf("Failed to publish message: %v", err)
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	// 等待一段时间让消息处理
	time.Sleep(2 * time.Second)

	// 获取当前测量值
	if measCap, ok := generator.GetCapacity("AnalogMeasurementCapacity"); ok {
		if analogCap, ok := measCap.(*capacities.AnalogMeasurementCapacity); ok {
			currentValue := analogCap.GetCurrentValue()
			if currentValue != nil {
				fmt.Printf("当前功率测量值: %.2f MW\n", currentValue.Value)
				fmt.Printf("时间戳: %s\n", currentValue.Timestamp.Format(time.RFC3339))
				fmt.Printf("质量: %v\n", currentValue.Quality)
			}
		}
	}

	// 发送控制命令（设定点）
	setPointMsg := &capacities.SetPointMessage{Value: 1200.0}
	generator.Send(setPointMsg)

	// 等待处理
	time.Sleep(100 * time.Millisecond)

	// 停止 Actor
	if err := generator.Stop(); err != nil {
		log.Fatalf("Failed to stop actor: %v", err)
	}

	fmt.Println("示例完成")
}

