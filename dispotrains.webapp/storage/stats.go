package storage

// StationStats holds operation statistics for a station.
type StationStats struct {
	Name            string
	Reports         int
	ReportDays      int
	Malfunctions    int
	MalfunctionDays int
	FunctionDays    int
	PercentFunction float64
	Elevators       map[string]ElevatorStats
}

// ElevatorStats holds operation statistics for an elevator. It is part of
// StationStats.
type ElevatorStats struct {
	Name            string
	Malfunctions    int
	MalfunctionDays int
	FunctionDays    int
	PercentFunction float64
}
