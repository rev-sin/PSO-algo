package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	vms := []*VM{
		NewVM(0, 100.0, 1024.0, 1000.0), // VM1: 100 CPU, 1GB RAM, 1Gbps
		NewVM(1, 150.0, 2048.0, 2000.0), // VM2: 150 CPU, 2GB RAM, 2Gbps
		NewVM(2, 80.0, 512.0, 500.0),    // VM3: 80 CPU, 512MB RAM, 500Mbps
	}

	// Initialize Tasks
	tasks := []*Task{
		NewTask(0, 20.0, 256.0, 200.0, 5.0),
		NewTask(1, 30.0, 512.0, 300.0, 7.0),
		NewTask(2, 15.0, 128.0, 100.0, 3.0),
		NewTask(3, 40.0, 768.0, 400.0, 10.0),
	}

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Run PSO Load Balancer")
		fmt.Println("2. Show VM Loads")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			runPSOLoadBalancer(vms, tasks)
		case 2:
			showVMLoads(vms)
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func runPSOLoadBalancer(vms []*VM, tasks []*Task) {
	swarm := NewSwarm(20, len(tasks), len(vms), vms, tasks)
	swarm.Optimize(100, 0.7, 2.0, 2.0)

	fmt.Printf("Best Fitness (Load Imbalance): %.2f\n", swarm.GBestFitness)
	fmt.Println("Best Task-to-VM Assignment:")
	for taskID, vmID := range swarm.GetBestAssignment() {
		fmt.Printf("Task %d -> VM %d\n", taskID, vmID)
	}
}

func showVMLoads(vms []*VM) {
	for _, vm := range vms {
		fmt.Printf("VM %d Load: %.2f%%\n", vm.ID, vm.CurrentLoad)
	}
}
