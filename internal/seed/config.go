package seed

type DelinquencyPoint struct {
	Year       int
	Month      int
	Percentage int
}

type OverdueAgeBucket struct {
	FromDays int
	ToDays   int
	Weight   int
}

type Config struct {
	MonthlyRevenueUSD            int
	ActiveCustomers              int
	MonthlyFeeFrom               int
	MonthlyFeeTo                 int
	DueDay                       int
	DelinquencyTrend             []DelinquencyPoint
	OverdueAgeBuckets            []OverdueAgeBucket
	DataFromYear                 int
	CustomerStartVariationMonths int
	ReviewedCustomersPercentage  int
	EnterpriseCompanyWeight      int
	PymeCompanyWeight            int
	StartupCompanyWeight         int
}

var config = Config{
	// Business goals
	MonthlyRevenueUSD: 380000,
	ActiveCustomers:   420,

	// Customer billing
	MonthlyFeeFrom:               200,
	MonthlyFeeTo:                 15000,
	DueDay:                       10,
	DataFromYear:                 2024,
	CustomerStartVariationMonths: 6,
	ReviewedCustomersPercentage:  20,

	// Company type distribution
	EnterpriseCompanyWeight: 20,
	PymeCompanyWeight:       40,
	StartupCompanyWeight:    40,

	// Delinquency behavior
	DelinquencyTrend: []DelinquencyPoint{
		{Year: 2024, Month: 1, Percentage: 3},
		{Year: 2025, Month: 5, Percentage: 6},
		{Year: 2026, Month: 5, Percentage: 14},
	},
	OverdueAgeBuckets: []OverdueAgeBucket{
		{FromDays: 1, ToDays: 30, Weight: 35},
		{FromDays: 31, ToDays: 75, Weight: 35},
		{FromDays: 76, ToDays: 90, Weight: 15},
		{FromDays: 91, ToDays: 140, Weight: 15},
	},
}
