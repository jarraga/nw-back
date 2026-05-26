package seed

import (
	"context"
	"log"
	"sort"
	"time"

	"nw-back/internal/postgres/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

func createCustomerPayments(ctx context.Context, queries *db.Queries, customers []db.Customer, config Config) error {
	dataTo := dataToMonth()
	now := time.Now().UTC()
	debtorStartMonths := randomDebtorStartMonths(customers, config, dataTo, now)
	latePaymentMonths := randomLatePaymentMonths(customers, config, dataTo, now, debtorStartMonths)
	paymentsCreated := 0

	for _, customer := range customers {
		dataFrom := monthStart(customer.BillingStartedAt.Time)
		payments, err := randomPayments(customer, dataFrom, dataTo, debtorStartMonths[customer.ID], latePaymentMonths, now, config.DueDay)
		if err != nil {
			return err
		}

		for _, payment := range payments {
			_, err = queries.CreateCustomerPayment(ctx, payment)
			if err != nil {
				return err
			}

			paymentsCreated++
		}
	}

	log.Printf("%d customer payments created", paymentsCreated)
	log.Printf("%d customers left with overdue debt", len(debtorStartMonths))
	return nil
}

func randomPayments(customer db.Customer, dataFrom time.Time, dataTo time.Time, firstUnpaidMonth time.Time, latePaymentMonths map[int]map[int64]bool, now time.Time, dueDay int) ([]db.CreateCustomerPaymentParams, error) {
	payments := []db.CreateCustomerPaymentParams{}
	currentMonth := dataFrom
	lastPaidAt := time.Time{}

	for !currentMonth.After(dataTo) {
		if !firstUnpaidMonth.IsZero() && !currentMonth.Before(firstUnpaidMonth) {
			break
		}

		payment, exists := randomPayment(customer, currentMonth, lastPaidAt, now, dueDay, isLatePaymentMonth(customer.ID, currentMonth, latePaymentMonths))
		if !exists {
			break
		}

		payments = append(payments, payment)
		lastPaidAt = payment.PaidAt.Time

		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	return payments, nil
}

func randomPayment(customer db.Customer, month time.Time, lastPaidAt time.Time, now time.Time, dueDay int, late bool) (db.CreateCustomerPaymentParams, bool) {
	paidAt := randomPaidAt(month, dueDay, late)
	if !lastPaidAt.IsZero() && !paidAt.After(lastPaidAt) {
		paidAt = lastPaidAt.AddDate(0, 0, gofakeit.Number(1, 5))
	}

	if paidAt.After(now) {
		if !isOverdueMonth(month, now, dueDay) {
			return db.CreateCustomerPaymentParams{}, false
		}

		paidAt = now
	}

	timestamptz := pgtype.Timestamptz{
		Time:  paidAt,
		Valid: true,
	}

	return db.CreateCustomerPaymentParams{
		CustomerID: customer.ID,
		Year:       int32(month.Year()),
		Month:      int32(month.Month()),
		Status:     db.PaymentStatusPaid,
		PaidAt:     timestamptz,
	}, true
}

func randomPaidAt(month time.Time, dueDay int, late bool) time.Time {
	if late {
		daysAfterDueDate := gofakeit.Number(1, 45)
		return dueDate(month, dueDay).AddDate(0, 0, daysAfterDueDate)
	}

	daysBeforeDueDate := gofakeit.Number(0, 7)
	return dueDate(month, dueDay).AddDate(0, 0, -daysBeforeDueDate)
}

func dataToMonth() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func monthStart(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func randomDebtorStartMonths(customers []db.Customer, config Config, dataTo time.Time, now time.Time) map[int64]time.Time {
	debtorStartMonths := map[int64]time.Time{}
	activeCustomers := activeCustomersForMonth(customers, dataTo)
	targetDebtors := targetDebtorsForMonth(len(activeCustomers), dataTo, config.DelinquencyTrend)

	for len(debtorStartMonths) < targetDebtors {
		customer, ok := randomNonDebtorCustomer(activeCustomers, debtorStartMonths)
		if !ok {
			break
		}

		firstUnpaidMonth := randomOverdueMonthByAge(monthStart(customer.BillingStartedAt.Time), dataTo, now, config)
		if firstUnpaidMonth.IsZero() {
			continue
		}

		debtorStartMonths[customer.ID] = firstUnpaidMonth
	}

	return debtorStartMonths
}

func activeCustomersForMonth(customers []db.Customer, month time.Time) []db.Customer {
	activeCustomers := []db.Customer{}

	for _, customer := range customers {
		if !monthStart(customer.BillingStartedAt.Time).After(month) {
			activeCustomers = append(activeCustomers, customer)
		}
	}

	return activeCustomers
}

func randomLatePaymentMonths(customers []db.Customer, config Config, dataTo time.Time, now time.Time, debtorStartMonths map[int64]time.Time) map[int]map[int64]bool {
	latePaymentMonths := map[int]map[int64]bool{}
	dataFrom := time.Date(config.DataFromYear, time.January, 1, 0, 0, 0, 0, time.UTC)

	for currentMonth := dataFrom; !currentMonth.After(dataTo); currentMonth = currentMonth.AddDate(0, 1, 0) {
		if !isOverdueMonth(currentMonth, now, config.DueDay) {
			continue
		}

		activeCustomers := activeCustomersForMonth(customers, currentMonth)
		targetDebtors := targetDebtorsForMonth(len(activeCustomers), currentMonth, config.DelinquencyTrend)
		unpaidDebtors := currentDebtorsForMonth(activeCustomers, debtorStartMonths, currentMonth)
		latePaymentsNeeded := targetDebtors - unpaidDebtors

		if latePaymentsNeeded <= 0 {
			continue
		}

		key := monthKey(currentMonth)
		latePaymentMonths[key] = map[int64]bool{}

		for len(latePaymentMonths[key]) < latePaymentsNeeded {
			customer, ok := randomCustomerWithPaymentInMonth(activeCustomers, debtorStartMonths, latePaymentMonths[key], currentMonth)
			if !ok {
				break
			}

			latePaymentMonths[key][customer.ID] = true
		}
	}

	return latePaymentMonths
}

func currentDebtorsForMonth(customers []db.Customer, debtorStartMonths map[int64]time.Time, month time.Time) int {
	total := 0

	for _, customer := range customers {
		startMonth, ok := debtorStartMonths[customer.ID]
		if ok && !startMonth.After(month) {
			total++
		}
	}

	return total
}

func randomCustomerWithPaymentInMonth(customers []db.Customer, debtorStartMonths map[int64]time.Time, selectedCustomers map[int64]bool, month time.Time) (db.Customer, bool) {
	candidates := []db.Customer{}

	for _, customer := range customers {
		if selectedCustomers[customer.ID] {
			continue
		}

		startMonth, isCurrentDebtor := debtorStartMonths[customer.ID]
		if isCurrentDebtor && !startMonth.After(month) {
			continue
		}

		candidates = append(candidates, customer)
	}

	if len(candidates) == 0 {
		return db.Customer{}, false
	}

	return candidates[gofakeit.Number(0, len(candidates)-1)], true
}

func isLatePaymentMonth(customerID int64, month time.Time, latePaymentMonths map[int]map[int64]bool) bool {
	customers, ok := latePaymentMonths[monthKey(month)]
	if !ok {
		return false
	}

	return customers[customerID]
}

func randomNonDebtorCustomer(customers []db.Customer, debtorStartMonths map[int64]time.Time) (db.Customer, bool) {
	nonDebtors := []db.Customer{}

	for _, customer := range customers {
		if _, ok := debtorStartMonths[customer.ID]; !ok {
			nonDebtors = append(nonDebtors, customer)
		}
	}

	if len(nonDebtors) == 0 {
		return db.Customer{}, false
	}

	return nonDebtors[gofakeit.Number(0, len(nonDebtors)-1)], true
}

func targetDebtorsForMonth(activeCustomers int, month time.Time, trend []DelinquencyPoint) int {
	percentage := delinquencyPercentageForMonth(month, trend)
	target := (activeCustomers*percentage + 50) / 100

	if target < 0 {
		return 0
	}

	if target > activeCustomers {
		return activeCustomers
	}

	return target
}

func delinquencyPercentageForMonth(month time.Time, trend []DelinquencyPoint) int {
	if len(trend) == 0 {
		return 0
	}

	points := append([]DelinquencyPoint{}, trend...)
	sort.Slice(points, func(i int, j int) bool {
		return delinquencyPointIndex(points[i]) < delinquencyPointIndex(points[j])
	})

	monthIndex := yearMonthIndex(month.Year(), int(month.Month()))
	firstPoint := points[0]
	firstPointIndex := delinquencyPointIndex(firstPoint)

	if monthIndex <= firstPointIndex {
		return normalizedPercentage(firstPoint.Percentage)
	}

	for i := 1; i < len(points); i++ {
		previous := points[i-1]
		next := points[i]
		previousIndex := delinquencyPointIndex(previous)
		nextIndex := delinquencyPointIndex(next)

		if monthIndex <= nextIndex {
			monthsBetween := nextIndex - previousIndex
			if monthsBetween <= 0 {
				return normalizedPercentage(next.Percentage)
			}

			monthsElapsed := monthIndex - previousIndex
			percentage := previous.Percentage + (next.Percentage-previous.Percentage)*monthsElapsed/monthsBetween
			return normalizedPercentage(percentage)
		}
	}

	return normalizedPercentage(points[len(points)-1].Percentage)
}

func delinquencyPointIndex(point DelinquencyPoint) int {
	return yearMonthIndex(point.Year, point.Month)
}

func yearMonthIndex(year int, month int) int {
	return year*12 + month
}

func monthKey(month time.Time) int {
	return month.Year()*100 + int(month.Month())
}

func normalizedPercentage(percentage int) int {
	if percentage < 0 {
		return 0
	}

	if percentage > 100 {
		return 100
	}

	return percentage
}

func randomOverdueMonthByAge(dataFrom time.Time, dataTo time.Time, now time.Time, config Config) time.Time {
	bucket := randomOverdueAgeBucket(config.OverdueAgeBuckets)
	candidates := overdueMonthsForAgeBucket(dataFrom, dataTo, now, config.DueDay, bucket)

	if len(candidates) == 0 {
		return randomOldestAvailableOverdueMonth(dataFrom, dataTo, now, config.DueDay)
	}

	return candidates[gofakeit.Number(0, len(candidates)-1)]
}

func randomOverdueAgeBucket(buckets []OverdueAgeBucket) OverdueAgeBucket {
	totalWeight := 0

	for _, bucket := range buckets {
		if bucket.Weight > 0 {
			totalWeight += bucket.Weight
		}
	}

	if totalWeight == 0 {
		return OverdueAgeBucket{FromDays: 1, ToDays: 30, Weight: 1}
	}

	randomWeight := gofakeit.Number(1, totalWeight)
	runningWeight := 0

	for _, bucket := range buckets {
		if bucket.Weight <= 0 {
			continue
		}

		runningWeight += bucket.Weight
		if randomWeight <= runningWeight {
			return bucket
		}
	}

	return buckets[len(buckets)-1]
}

func overdueMonthsForAgeBucket(dataFrom time.Time, dataTo time.Time, now time.Time, dueDay int, bucket OverdueAgeBucket) []time.Time {
	months := []time.Time{}

	for currentMonth := dataFrom; !currentMonth.After(dataTo); currentMonth = currentMonth.AddDate(0, 1, 0) {
		if !isOverdueMonth(currentMonth, now, dueDay) {
			continue
		}

		ageDays := int(currentDate(now).Sub(dueDate(currentMonth, dueDay)).Hours() / 24)
		if ageDays >= bucket.FromDays && ageDays <= bucket.ToDays {
			months = append(months, currentMonth)
		}
	}

	return months
}

func randomOldestAvailableOverdueMonth(dataFrom time.Time, dataTo time.Time, now time.Time, dueDay int) time.Time {
	months := []time.Time{}

	for currentMonth := dataFrom; !currentMonth.After(dataTo); currentMonth = currentMonth.AddDate(0, 1, 0) {
		if isOverdueMonth(currentMonth, now, dueDay) {
			months = append(months, currentMonth)
		}
	}

	if len(months) == 0 {
		return time.Time{}
	}

	return months[gofakeit.Number(0, len(months)-1)]
}

func isOverdueMonth(month time.Time, now time.Time, dueDay int) bool {
	return dueDate(month, dueDay).Before(currentDate(now))
}

func dueDate(month time.Time, dueDay int) time.Time {
	lastDayOfMonth := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dueDay = normalizedDueDay(dueDay)

	if dueDay > lastDayOfMonth {
		dueDay = lastDayOfMonth
	}

	return time.Date(month.Year(), month.Month(), dueDay, 0, 0, 0, 0, time.UTC)
}

func currentDate(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func normalizedDueDay(dueDay int) int {
	if dueDay < 1 {
		return 10
	}

	if dueDay > 31 {
		return 31
	}

	return dueDay
}
