package ussd

import (
	"fmt"
	"marketlens/internal/models"
	"strings"
)

// pageItems returns the slice of items for the given page, plus whether there is a next page.
func pageItems[T any](items []T, page, perPage int) (visible []T, hasNext bool) {
	start := page * perPage
	if start >= len(items) {
		return nil, false
	}
	end := start + perPage
	if end >= len(items) {
		return items[start:], false
	}
	return items[start:end], true
}

func renderMainMenu() string {
	return "CON Welcome to MarketLens\n1. Latest price\n2. Report price\n3. Help\n4. Exit"
}

func renderSelectState(states []string, page, perPage int) string {
	if len(states) == 0 {
		return "CON No states available.\n0. Back"
	}
	visible, hasNext := pageItems(states, page, perPage)
	if len(visible) == 0 {
		return "CON No more states.\n0. Back"
	}

	var b strings.Builder
	b.WriteString("CON Select your state:\n")
	for i, s := range visible {
		fmt.Fprintf(&b, "%d. %s\n", i+1, s)
	}
	if hasNext {
		b.WriteString("00. Next\n")
	}
	if page > 0 {
		b.WriteString("98. Previous\n")
	}
	b.WriteString("0. Back")
	return b.String()
}

func renderSelectCrop(crops []models.Crop, page, perPage int) string {
	if len(crops) == 0 {
		return "CON No crops available right now.\n0. Back"
	}
	visible, hasNext := pageItems(crops, page, perPage)
	if len(visible) == 0 {
		return "CON No more crops.\n0. Back"
	}

	var b strings.Builder
	b.WriteString("CON Select crop:\n")
	for i, c := range visible {
		fmt.Fprintf(&b, "%d. %s\n", i+1, c.Name)
	}
	if hasNext {
		b.WriteString("00. Next\n")
	}
	if page > 0 {
		b.WriteString("98. Previous\n")
	}
	b.WriteString("0. Back")
	return b.String()
}

func renderSelectMarket(markets []models.Market, page, perPage int) string {
	if len(markets) == 0 {
		return "CON No markets in this state.\n0. Back"
	}
	visible, hasNext := pageItems(markets, page, perPage)
	if len(visible) == 0 {
		return "CON No more markets.\n0. Back"
	}

	var b strings.Builder
	b.WriteString("CON Select market:\n")
	for i, m := range visible {
		fmt.Fprintf(&b, "%d. %s (%s)\n", i+1, m.Name, m.State)
	}
	if hasNext {
		b.WriteString("00. Next\n")
	}
	if page > 0 {
		b.WriteString("98. Previous\n")
	}
	b.WriteString("0. Back")
	return b.String()
}

func renderShowPrice(agg *models.AggregatedPrice) string {
	if agg == nil {
		return "END Price not available at the moment."
	}
	return fmt.Sprintf(
		"END %s @ %s\nPrice: %s%.2f - %s%.2f per %s\nTrend: %s\nConfidence: %s\nUpdated: %s",
		agg.CropName, agg.MarketName,
		agg.Currency, agg.PriceMin,
		agg.Currency, agg.PriceMax,
		agg.Unit,
		agg.Trend,
		agg.Confidence,
		agg.UpdatedAt.Format("2006-01-02 15:04"),
	)
}

func renderEnterPrice(cropName, marketName string) string {
	return fmt.Sprintf("CON Report price for\n%s @ %s\n\nEnter price in Naira:", cropName, marketName)
}

func renderConfirmPrice(cropName, marketName string, price float64) string {
	return fmt.Sprintf("CON You are reporting:\n%s @ %s\nPrice: NGN%.2f\n\n1. Confirm\n2. Re-enter price\n0. Cancel", cropName, marketName, price)
}

func renderPriceSubmitted() string {
	return "END Thank you! Your price report has been submitted and will be reviewed."
}

func renderHelp() string {
	return "END MarketLens helps you check daily crop prices.\nDial again and select option 1 to get started."
}

func renderError() string {
	return "CON Invalid choice. Please try again."
}

func renderGoodbye() string {
	return "END Thank you for using MarketLens!"
}
