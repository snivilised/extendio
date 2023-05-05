package nav

type MetricEnum uint

// if new metrics are added, ensure that navigationMetricsFactory.new is kept
// in sync.
const (
	// MetricNoFilesInvokedEn represents the no of files invoked for during traversal
	//
	MetricNoFilesInvokedEn MetricEnum = iota

	// MetricNoFilesFilteredOutEn
	//
	MetricNoFilesFilteredOutEn

	// MetricNoFoldersInvokedEn represents the no of folders invoked for during traversal
	//
	MetricNoFoldersInvokedEn

	// MetricNoFoldersFilteredOutEn
	//
	MetricNoFoldersFilteredOutEn

	// MetricNoChildFilesFoundEn represents the number of children files
	// of a particular directory that pass the compound filter when using the folders
	// with files subscription
	//
	MetricNoChildFilesFoundEn

	// MetricNoChildFilesFilteredOutEn represents the number of children files
	// of a particular directory that fail to pass the compound filter when using
	// the folders with files subscription
	//
	MetricNoChildFilesFilteredOutEn
)

// Metric
type Metric struct {
	Name  string
	Count uint
}

// MetricCollection
// TODO: make this private as it's internal implementation detail
type MetricCollection map[MetricEnum]*Metric

type NavigationMetrics struct {
	collection MetricCollection
}

func (m *NavigationMetrics) Count(metricEn MetricEnum) uint {
	var result uint

	if m, found := m.collection[metricEn]; found {
		result = m.Count
	}

	return result
}

func (m *NavigationMetrics) tick(metricEn MetricEnum) {
	m.collection[metricEn].Count++
}

func (m *NavigationMetrics) post(metricEn MetricEnum, count uint) {
	m.collection[metricEn].Count += count
}

func (m *NavigationMetrics) save(active *ActiveState) {
	active.Metrics = &m.collection
}

func (m *NavigationMetrics) load(active *ActiveState) {
	m.collection = *active.Metrics
}

type navigationMetricsFactory struct{}

func (f navigationMetricsFactory) new() *NavigationMetrics {
	instance := &NavigationMetrics{
		collection: make(MetricCollection),
	}
	instance.collection[MetricNoFilesInvokedEn] = &Metric{Name: "filesInvoked"}
	instance.collection[MetricNoFilesFilteredOutEn] = &Metric{Name: "filesFilteredOut"}
	instance.collection[MetricNoFoldersInvokedEn] = &Metric{Name: "foldersInvoked"}
	instance.collection[MetricNoFoldersFilteredOutEn] = &Metric{Name: "foldersFilteredOut"}
	instance.collection[MetricNoChildFilesFoundEn] = &Metric{Name: "childrenFound"}
	instance.collection[MetricNoChildFilesFilteredOutEn] = &Metric{Name: "childrenFilteredOut"}

	return instance
}
