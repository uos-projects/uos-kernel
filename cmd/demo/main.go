package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/uos-projects/uos-kernel/internal/importer"
	"github.com/uos-projects/uos-kernel/kernel"
)

func main() {
	csvPath := "data/defects.csv"
	if len(os.Args) > 1 {
		csvPath = os.Args[1]
	}

	fmt.Println("=== UOS Kernel Demo: 人机物统一资源 + 数字孪生 + 世界模型 ===")
	fmt.Println()

	// 1. 加载数据
	fmt.Printf("加载数据: %s\n", csvPath)
	world, stats, err := importer.LoadDefects(csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n--- 世界模型统计 ---\n")
	fmt.Printf("  人力资源:   %d 人\n", stats.Persons)
	fmt.Printf("  区域:       %d 个\n", stats.Areas)
	fmt.Printf("  变电站:     %d 座\n", stats.Substations)
	fmt.Printf("  馈线:       %d 条\n", stats.Feeders)
	fmt.Printf("  开关站:     %d 座\n", stats.SwitchStations)
	fmt.Printf("  终端设备:   %d 台\n", stats.Terminals)
	fmt.Printf("  缺陷工单:   %d 张\n", stats.Tickets)
	fmt.Printf("  拓扑关系:   %d 条\n", stats.Edges)
	fmt.Println()

	ctx := context.Background()

	// 2. 打开一张工单，查看孪生
	fmt.Println("--- Demo 1: 读取工单孪生 ---")
	fd, err := world.Rosix.Open(ctx, "ticket:00001")
	if err != nil {
		fmt.Printf("  Open失败: %v\n", err)
	} else {
		twin, _ := world.Rosix.Read(ctx, fd)
		fmt.Printf("  工单: %s\n", twin.Resource.Name)
		fmt.Printf("  当前状态: %s\n", twin.Current["status"])
		fmt.Printf("  缺陷类型: %s\n", twin.Resource.Attributes["defect_type"])
		fmt.Printf("  解决过程: %s\n", truncate(twin.Resource.Attributes["resolution"], 60))

		// 查看时间线
		events, _ := world.Rosix.History(ctx, fd, time.Time{}, time.Now())
		fmt.Printf("  状态变迁 (%d步):\n", len(events))
		for _, e := range events {
			fmt.Printf("    %s  %s → %s  (%s)\n",
				e.Timestamp.Format("01-02 15:04"), e.OldValue, e.NewValue, truncate(e.Cause, 30))
		}
		world.Rosix.Close(fd)
	}
	fmt.Println()

	// 3. 查询人员负责的工单
	fmt.Println("--- Demo 2: 拓扑遍历 — 人员 → 工单 ---")
	personID := kernel.ResourceID("person:周留康")
	pfd, err := world.Rosix.Open(ctx, personID)
	if err != nil {
		fmt.Printf("  Open失败: %v\n", err)
	} else {
		twin, _ := world.Rosix.Read(ctx, pfd)
		fmt.Printf("  人员: %s (区域: %s)\n", twin.Resource.Name, twin.Resource.Attributes["area"])

		ticketFDs, _ := world.Rosix.Traverse(ctx, pfd, kernel.Inbound, kernel.RelAssigned)
		fmt.Printf("  负责工单数: %d\n", len(ticketFDs))
		for i, tfd := range ticketFDs {
			if i >= 3 {
				fmt.Printf("    ... 还有 %d 张\n", len(ticketFDs)-3)
				break
			}
			tt, _ := world.Rosix.Read(ctx, tfd)
			fmt.Printf("    [%s] %s\n", tt.Current["status"], truncate(tt.Resource.Name, 50))
			world.Rosix.Close(tfd)
		}
		world.Rosix.Close(pfd)
	}
	fmt.Println()

	// 4. 查询变电站物理拓扑
	fmt.Println("--- Demo 3: 拓扑遍历 — 变电站 → 设备 ---")
	subID := kernel.ResourceID("substation:苏州.潘阳变")
	sfd, err := world.Rosix.Open(ctx, subID)
	if err != nil {
		fmt.Printf("  Open失败: %v\n", err)
	} else {
		twin, _ := world.Rosix.Read(ctx, sfd)
		fmt.Printf("  变电站: %s\n", twin.Resource.Name)

		// 直接子节点（馈线）
		feederFDs, _ := world.Rosix.Traverse(ctx, sfd, kernel.Outbound, kernel.RelContains)
		fmt.Printf("  下辖馈线: %d 条\n", len(feederFDs))
		for i, ffd := range feederFDs {
			if i >= 5 {
				fmt.Printf("    ... 还有 %d 条\n", len(feederFDs)-5)
				break
			}
			ft, _ := world.Rosix.Read(ctx, ffd)
			fmt.Printf("    %s\n", ft.Resource.Name)
			world.Rosix.Close(ffd)
		}
		world.Rosix.Close(sfd)
	}
	fmt.Println()

	// 5. 模拟写入新状态
	fmt.Println("--- Demo 4: 写入新状态（模拟消缺） ---")
	fd2, err := world.Rosix.Open(ctx, "ticket:00002")
	if err != nil {
		fmt.Printf("  Open失败: %v\n", err)
	} else {
		before, _ := world.Rosix.Read(ctx, fd2)
		fmt.Printf("  工单: %s\n", before.Resource.Name)
		fmt.Printf("  当前状态: %s\n", before.Current["status"])

		err = world.Rosix.Write(ctx, fd2, "status", "re-repaired", "二次远程重启", "person:周留康")
		if err != nil {
			fmt.Printf("  写入失败: %v\n", err)
		} else {
			after, _ := world.Rosix.Read(ctx, fd2)
			fmt.Printf("  写入后状态: %s\n", after.Current["status"])

			events, _ := world.Rosix.History(ctx, fd2, time.Time{}, time.Now())
			fmt.Printf("  时间线现在有 %d 个事件\n", len(events))
		}
		world.Rosix.Close(fd2)
	}

	fmt.Println()
	fmt.Println("=== Demo 完成 ===")
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
