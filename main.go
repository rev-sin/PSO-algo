package main

import (
    "fmt"
    "math/rand"
    "time"
    "math"

    ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)

func main() {
    rand.Seed(time.Now().UnixNano())

    vms := []*VM{
        NewVM(0, 100.0, 1024.0, 1000.0),
        NewVM(1, 150.0, 2048.0, 2000.0),
        NewVM(2, 80.0, 512.0, 500.0),
    }

    tasks := []*Task{
        NewTask(0, 20.0, 256.0, 200.0, 5.0),
        NewTask(1, 30.0, 512.0, 300.0, 7.0),
        NewTask(2, 15.0, 128.0, 100.0, 3.0),
        NewTask(3, 40.0, 768.0, 400.0, 10.0),
    }

    if err := ui.Init(); err != nil {
        fmt.Printf("Failed to initialize termui: %v\n", err)
        return
    }
    defer ui.Close()

    runRealTimeSimulation(vms, tasks)
}

func runRealTimeSimulation(vms []*VM, tasks []*Task) {
    swarm := NewSwarm(20, len(tasks), len(vms), vms, tasks)

    vmLoadChart := widgets.NewBarChart()
    vmLoadChart.Title = "VM Loads"
    vmLoadChart.Labels = []string{"VM0", "VM1", "VM2"}
    vmLoadChart.BarWidth = 5
    vmLoadChart.BarColors = []ui.Color{ui.ColorGreen, ui.ColorYellow, ui.ColorRed}
    vmLoadChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
    vmLoadChart.BarGap = 2

    taskAssignment := widgets.NewParagraph()
    taskAssignment.Title = "Task Assignments"
    taskAssignment.Text = ""

    grid := ui.NewGrid()
    termWidth, termHeight := ui.TerminalDimensions()
    grid.SetRect(0, 0, termWidth, termHeight)
    grid.Set(
        ui.NewRow(0.5, vmLoadChart),
        ui.NewRow(0.5, taskAssignment),
    )

    ui.Render(grid)

    go func() {
        for iter := 0; iter < 100; iter++ {
            swarm.Optimize(1, 0.7, 2.0, 2.0)

            vmLoads := make([]float64, len(vms))
            for i, vm := range vms {
                vmLoads[i] = math.Round(vm.CurrentLoad*100) / 100
            }
            vmLoadChart.Data = vmLoads

            assignments := ""
            for taskID, vmID := range swarm.GetBestAssignment() {
                assignments += fmt.Sprintf("Task %d -> VM %d\n", taskID, vmID)
            }
            taskAssignment.Text = assignments

            ui.Render(grid)
            
            time.Sleep(1000 * time.Millisecond)
        }
    }()

    for e := range ui.PollEvents() {
        if e.Type == ui.KeyboardEvent && e.ID == "q" {
            break
        }
    }
}