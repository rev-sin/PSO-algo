package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"

    ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)

// PerformanceMetrics stores statistics across iterations
type PerformanceMetrics struct {
    BestFitness []float64
    StartTime   time.Time
    EndTime     time.Time
}

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
    metrics := &PerformanceMetrics{BestFitness: []float64{}, StartTime: time.Now()}

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
        numIterations := 100
        for iter := 0; iter < numIterations; iter++ {
            swarm.Optimize(1, 0.7, 2.0, 2.0)
            bestFitness := swarm.GetBestFitness()
            metrics.BestFitness = append(metrics.BestFitness, bestFitness)

            // Update VM load chart
            vmLoads := make([]float64, len(vms))
            for i, vm := range vms {
                vmLoads[i] = math.Round(vm.CurrentLoad*100) / 100
            }
            vmLoadChart.Data = vmLoads

            // Update task assignments
            assignments := ""
            for taskID, vmID := range swarm.GetBestAssignment() {
                assignments += fmt.Sprintf("Task %d -> VM %d\n", taskID, vmID)
            }
            taskAssignment.Text = fmt.Sprintf("Iteration: %d\nBest Fitness: %.4f\n\n%s", iter+1, bestFitness, assignments)

            ui.Render(grid)
            time.Sleep(500 * time.Millisecond)
        }

        metrics.EndTime = time.Now()
        printMetricsSummary(metrics)
    }()

    for e := range ui.PollEvents() {
        if e.Type == ui.KeyboardEvent && e.ID == "q" {
            break
        }
    }
}

// Compute and print benchmark metrics
func printMetricsSummary(m *PerformanceMetrics) {
    n := float64(len(m.BestFitness))
    if n == 0 {
        fmt.Println("No fitness data recorded.")
        return
    }

    var sum, best, worst float64
    best = math.MaxFloat64
    for _, f := range m.BestFitness {
        sum += f
        if f < best {
            best = f
        }
        if f > worst {
            worst = f
        }
    }
    mean := sum / n

    var variance float64
    for _, f := range m.BestFitness {
        variance += math.Pow(f-mean, 2)
    }
    stddev := math.Sqrt(variance / n)
    duration := m.EndTime.Sub(m.StartTime)

    fmt.Println("\n==================== PSO PERFORMANCE METRICS ====================")
    fmt.Printf("Total Iterations     : %d\n", int(n))
    fmt.Printf("Best Fitness Found   : %.6f\n", best)
    fmt.Printf("Worst Fitness Found  : %.6f\n", worst)
    fmt.Printf("Average Fitness      : %.6f\n", mean)
    fmt.Printf("Std. Deviation       : %.6f\n", stddev)
    fmt.Printf("Total Runtime        : %.3fs\n", duration.Seconds())
    fmt.Println("================================================================")
}
