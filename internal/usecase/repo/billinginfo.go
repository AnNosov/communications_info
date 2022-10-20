package repo

import (
	"fmt"
	"math"
	"strings"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/filereader"
)

func getBillingNumber(bitArray []string) int {
	var result float64 = 0
	for i := len(bitArray) - 1; i >= 0; i-- {
		if bitArray[i] == "1" {
			result += math.Pow(2, float64(len(bitArray)-i-1))
		}
	}
	return int(result)
}

func haveBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func makeBillingProfile(num int) (billingData entity.BillingData) {

	for i := 0; i < 7; i++ {
		switch i {
		case 0:
			billingData.CreateCustomer = haveBit(num, uint(i))
		case 1:
			billingData.Purchase = haveBit(num, uint(i))
		case 2:
			billingData.Payout = haveBit(num, uint(i))
		case 3:
			billingData.Recurring = haveBit(num, uint(i))
		case 4:
			billingData.FraudControl = haveBit(num, uint(i))
		case 5:
			billingData.CheckoutPage = haveBit(num, uint(i))
		}

	}
	return billingData
}

func GetBillingInfo(filePath, separator string) (entity.BillingData, error) {

	buf, err := filereader.GetCsvContent(filePath)
	if err != nil {
		return entity.BillingData{}, fmt.Errorf("billinginfo - GetBillingInfo: %w", err)
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 1 {
		return entity.BillingData{}, fmt.Errorf("getBillingInfo: incorrect lines in file: %v", len(lines))
	}
	line := strings.Split(lines[0], separator)
	if len(line) != 6 {
		return entity.BillingData{}, fmt.Errorf("getBillingInfo: incorrect size of bit mask: %v", len(line))
	}
	billingNum := getBillingNumber(line)
	billingData := makeBillingProfile(billingNum)

	return billingData, nil
}
