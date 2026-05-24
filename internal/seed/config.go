package seed

type Config struct {
	MonthlyRevenueUSD            int
	ActiveCustomers              int
	MonthlyFeeFrom               int
	MonthlyFeeTo                 int
	CurrentDelinquencyPercentage int
	DataFromYear                 int
}

var config = Config{
	MonthlyRevenueUSD:            380000,
	ActiveCustomers:              420,
	MonthlyFeeFrom:               200,
	MonthlyFeeTo:                 15000,
	CurrentDelinquencyPercentage: 14,
	DataFromYear:                 2014,
}
