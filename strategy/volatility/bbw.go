package volatility

import (
	"fmt"

	"github.com/miromax42/indicator/v2/asset"
	"github.com/miromax42/indicator/v2/helper"
	"github.com/miromax42/indicator/v2/strategy"
	"github.com/miromax42/indicator/v2/volatility"
)

type BollingerBandsWidthStrategy struct {
	BollingerBands *volatility.BollingerBands[float64]
	Sensitivity    float64 //0-1
}

func NewBollingerBandsWidthStrategy(p int, s float64) *BollingerBandsWidthStrategy {
	b := volatility.NewBollingerBands[float64]()
	b.Period = p

	return &BollingerBandsWidthStrategy{
		BollingerBands: b,
		Sensitivity:    s,
	}
}

func (b *BollingerBandsWidthStrategy) Name() string {
	return fmt.Sprintf("BBW(%d,%.3f)", b.BollingerBands.Period, b.Sensitivity)
}

func (b *BollingerBandsWidthStrategy) Compute(snapshots <-chan *asset.Snapshot) <-chan strategy.Action {
	closings := helper.Duplicate(
		asset.SnapshotsAsClosings(snapshots),
		2,
	)

	uppers, middles, lowers := b.BollingerBands.Compute(closings[0])
	go helper.Drain(middles)

	closings[1] = helper.Skip(closings[1], b.BollingerBands.IdlePeriod())

	actions := helper.Operate3(uppers, lowers, closings[1], func(upper, lower, closing float64) strategy.Action {
		if closing == 0 {
			return strategy.Hold
		}
		width := (upper - lower) / closing
		if width >= b.Sensitivity {
			return strategy.Buy
		}

		return strategy.Hold
	})

	actions = helper.Shift(actions, b.BollingerBands.IdlePeriod(), strategy.Hold)

	return actions
}

// Report processes the provided asset snapshots and generates a report annotated with the recommended actions.
func (b *BollingerBandsWidthStrategy) Report(c <-chan *asset.Snapshot) *helper.Report {
	//
	// snapshots[0] -> dates
	// snapshots[1] -> closings[0] -> closings
	//                 closings[1] -> upper
	//                             -> middle
	//                             -> lower
	// snapshots[2] -> actions     -> annotations
	//              -> outcomes
	//
	snapshots := helper.Duplicate(c, 3)

	dates := asset.SnapshotsAsDates(snapshots[0])
	closings := helper.Duplicate(asset.SnapshotsAsClosings(snapshots[1]), 2)

	uppers, middles, lowers := b.BollingerBands.Compute(closings[0])
	uppers = helper.Shift(uppers, b.BollingerBands.IdlePeriod(), 0)
	middles = helper.Shift(middles, b.BollingerBands.IdlePeriod(), 0)
	lowers = helper.Shift(lowers, b.BollingerBands.IdlePeriod(), 0)

	actions, outcomes := strategy.ComputeWithOutcome(b, snapshots[2])
	annotations := strategy.ActionsToAnnotations(actions)
	outcomes = helper.MultiplyBy(outcomes, 100)

	report := helper.NewReport(b.Name(), dates)
	report.AddChart()

	report.AddColumn(helper.NewNumericReportColumn("Close", closings[1]))
	report.AddColumn(helper.NewNumericReportColumn("Upper", uppers))
	report.AddColumn(helper.NewNumericReportColumn("Middle", middles))
	report.AddColumn(helper.NewNumericReportColumn("Lower", lowers))
	report.AddColumn(helper.NewAnnotationReportColumn(annotations))

	report.AddColumn(helper.NewNumericReportColumn("Outcome", outcomes), 1)

	return report
}
