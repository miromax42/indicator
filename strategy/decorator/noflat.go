package decorator

import (
	"fmt"

	"github.com/miromax42/indicator/v2/asset"
	"github.com/miromax42/indicator/v2/helper"
	"github.com/miromax42/indicator/v2/strategy"
)

// NoFlatStrategy is a decorator that prevents trading when the market is flat.
type NoFlatStrategy struct {
	strategy.Strategy

	InnerStrategy strategy.Strategy
	Period        int     // Period over which to measure flatness
	FlatThreshold float64 // Threshold for flatness (e.g., 1.0 for 1%)
}

// NewNoFlatStrategy creates a new instance of NoFlatStrategy.
func NewNoFlatStrategy(period int, flatThreshold float64, innerStrategy strategy.Strategy) *NoFlatStrategy {
	return &NoFlatStrategy{
		InnerStrategy: innerStrategy,
		Period:        period,
		FlatThreshold: flatThreshold,
	}
}

// Name returns the name of the strategy.
func (n *NoFlatStrategy) Name() string {
	return fmt.Sprintf("NoFlat(%d,%.f,%s)", n.Period, n.FlatThreshold, n.InnerStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (n *NoFlatStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	// Duplicate snapshots into two channels: one for the inner strategy, one for flatness calculation.
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	// Get the actions from the inner strategy.
	innerActions := n.InnerStrategy.Compute(snapshotsSplice[0])

	// Get the snapshots for computing flatness.
	flatnessSnapshots := snapshotsSplice[1]

	// Create the output channel.
	actions := make(chan strategy.Action)

	go func() {
		defer close(actions)

		// Initialize a slice to hold closing prices for the flatness calculation.
		closingWindow := make([]float64, 0, n.Period)

		// Read from the innerActions and flatnessSnapshots channels.
		for {
			action, ok1 := <-innerActions
			snapshot, ok2 := <-flatnessSnapshots

			if !ok1 || !ok2 {
				break
			}

			// Append the current closing price to the window.
			closingWindow = append(closingWindow, snapshot.Close)

			// Maintain the window size.
			if len(closingWindow) > n.Period {
				closingWindow = closingWindow[1:]
			}

			// Initialize variables for flatness.
			var isFlat bool

			// Perform the flatness calculation if enough data is available.
			if len(closingWindow) == n.Period {
				_, isFlat = n.calculateFlatness(closingWindow)
			} else {
				// Not enough data yet; default to not flat.
				isFlat = false
			}

			// If the market is flat, issue a Hold action.
			if isFlat {
				actions <- strategy.Hold
			} else {
				// Otherwise, pass through the inner strategy's action.
				actions <- action
			}
		}
	}()

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (n *NoFlatStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 4)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(n, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	// For reporting flatness, compute the flatness measure over the snapshots.
	flatnessMeasuresChan := n.computeFlatnessMeasures(snapshots[3])

	report := helper.NewReport(n.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))
	report.AddColumn(helper.NewNumericReportColumn("Flatness", flatnessMeasuresChan))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}

// computeFlatnessMeasures computes the flatness measures over the snapshots for reporting.
func (n *NoFlatStrategy) computeFlatnessMeasures(snapshots <-chan *asset.Snapshot) <-chan float64 {
	flatnessMeasures := make(chan float64)

	go func() {
		defer close(flatnessMeasures)
		closingWindow := make([]float64, 0, n.Period)

		for snapshot := range snapshots {
			closingWindow = append(closingWindow, snapshot.Close)
			if len(closingWindow) > n.Period {
				closingWindow = closingWindow[1:]
			}

			var flatnessMeasure float64

			if len(closingWindow) == n.Period {
				flatnessMeasure, _ = n.calculateFlatness(closingWindow)
			} else {
				// Not enough data yet; set flatness measure to zero.
				flatnessMeasure = 0
			}

			flatnessMeasures <- flatnessMeasure
		}
	}()

	return flatnessMeasures
}

// calculateFlatness calculates the flatness measure and determines if the market is flat.
func (n *NoFlatStrategy) calculateFlatness(closingWindow []float64) (float64, bool) {
	minPrice := closingWindow[0]
	maxPrice := closingWindow[0]
	totalPrice := 0.0

	for _, price := range closingWindow {
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
		totalPrice += price
	}

	avgPrice := totalPrice / float64(len(closingWindow))
	priceRange := maxPrice - minPrice
	flatnessMeasure := priceRange / avgPrice * 100 // Multiply by 100 for percentage

	isFlat := flatnessMeasure < n.FlatThreshold // Threshold is already in percentage

	return flatnessMeasure, isFlat
}
