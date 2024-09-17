package decorator

import (
	"fmt"

	"github.com/miromax42/indicator/v2/asset"
	"github.com/miromax42/indicator/v2/helper"
	"github.com/miromax42/indicator/v2/strategy"
)

type DelayStrategy struct {
	strategy.Strategy

	InnertStrategy strategy.Strategy
	Period         int
}

func NewDelayStrategy(period int, innerStrategy strategy.Strategy) *DelayStrategy {
	return &DelayStrategy{
		InnertStrategy: innerStrategy,
		Period:         period,
	}
}

// Name returns the name of the strategy.
func (n *DelayStrategy) Name() string {
	return fmt.Sprintf("D(%s)", n.InnertStrategy.Name())
}

// Compute processes the provided asset snapshots and generates a stream of actionable recommendations.
func (n *DelayStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	snapshotsSplice := helper.Duplicate(snapshots, 2)

	innerActions := n.InnertStrategy.Compute(snapshotsSplice[0])
	closings := asset.SnapshotsAsClosings(snapshotsSplice[1])
	var cnt int
	return helper.Operate(innerActions, closings, func(action strategy.Action, closing float64) strategy.Action {
		cnt++
		if cnt <= n.Period {
			return strategy.Hold
		}

		return action
	})
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (n *DelayStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := asset.SnapshotsAsClosings(snapshots[1])

	actions, outcomes := strategy.ComputeWithOutcome(n, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(n.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
