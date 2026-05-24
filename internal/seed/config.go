package seed

type Config struct {
	MonthlyRevenueUSD              int
	ActiveCustomers                int
	MonthlyFeeFrom                 int
	MonthlyFeeTo                   int
	CurrentDelinquencyPercentage   int
	DataFromYear                   int
	CustomerStartVariationMonths   int
	EnterpriseCompanyWeight        int
	PymeCompanyWeight              int
	StartupCompanyWeight           int
	EnterprisePaymentDelayFromDays int
	EnterprisePaymentDelayToDays   int
	GeneralPaymentDelayFromDays    int
	GeneralPaymentDelayToDays      int
}

var config = Config{
	// Business goals
	MonthlyRevenueUSD:            380000,
	ActiveCustomers:              420,
	CurrentDelinquencyPercentage: 14,

	// Customer billing
	MonthlyFeeFrom:               200,
	MonthlyFeeTo:                 15000,
	DataFromYear:                 2024,
	CustomerStartVariationMonths: 6,

	// Company type distribution
	EnterpriseCompanyWeight: 20,
	PymeCompanyWeight:       40,
	StartupCompanyWeight:    40,

	// Payment delay behavior
	EnterprisePaymentDelayFromDays: 0,
	EnterprisePaymentDelayToDays:   100,
	GeneralPaymentDelayFromDays:    0,
	GeneralPaymentDelayToDays:      45,
}
