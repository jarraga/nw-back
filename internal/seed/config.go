package seed

type Config struct {
	MonthlyRevenueUSD              int
	ActiveCustomers                int
	MonthlyFeeFrom                 int
	MonthlyFeeTo                   int
	CurrentDelinquencyPercentage   int
	DataFromYear                   int
	CustomerStartVariationMonths   int
	EnterprisePaymentDelayFromDays int
	EnterprisePaymentDelayToDays   int
	GeneralPaymentDelayFromDays    int
	GeneralPaymentDelayToDays      int
}

var config = Config{
	MonthlyRevenueUSD:              380000,
	ActiveCustomers:                420,
	MonthlyFeeFrom:                 200,
	MonthlyFeeTo:                   15000,
	CurrentDelinquencyPercentage:   14,
	DataFromYear:                   2024,
	CustomerStartVariationMonths:   6,
	EnterprisePaymentDelayFromDays: 0,
	EnterprisePaymentDelayToDays:   100,
	GeneralPaymentDelayFromDays:    0,
	GeneralPaymentDelayToDays:      45,
}
