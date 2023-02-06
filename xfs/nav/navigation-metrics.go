package nav

type MetricEnum uint

const (
	// No of files invoked for during traversal
	//
	MetricNoFilesEn MetricEnum = iota

	// No of folders invoked for during traversal
	//
	MetricNoFoldersEn
)

// Metric
type Metric struct {
	Name  string
	Count uint
}

// MetricCollection
type MetricCollection map[MetricEnum]*Metric

type navigationMetrics struct {
	_metrics MetricCollection
}

func (m *navigationMetrics) tick(metricEn MetricEnum) {
	m._metrics[metricEn].Count++
}

func (m *navigationMetrics) save(active *ActiveState) {
	active.Metrics = &m._metrics
}

func (m *navigationMetrics) load(active *ActiveState) {
	m._metrics = *active.Metrics
}

type navigationMetricsFactory struct{}

func (f *navigationMetricsFactory) construct() *navigationMetrics {
	instance := &navigationMetrics{
		_metrics: make(MetricCollection),
	}
	instance._metrics[MetricNoFilesEn] = &Metric{Name: "files"}
	instance._metrics[MetricNoFoldersEn] = &Metric{Name: "folders"}

	return instance
}
