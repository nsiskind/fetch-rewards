package server

import (
	"testing"

	"fetch-challenge/models"
)

func TestNamePoints(t *testing.T) {
	var tests = []struct {
		name     string
		retailer string
		expected int
	}{
		{
			"simple name",
			"target123",
			9,
		},
		{
			"complex name",
			"t!t!",
			2,
		},
		{
			"no alphanumeric",
			"*(",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getRetailerPoints(tt.retailer)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestTotalIsRoundPoints(t *testing.T) {
	var tests = []struct {
		name     string
		total    string
		expected int
	}{
		{
			"round total",
			"12",
			50,
		},
		{
			"not round",
			"12.234",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getTotalIsRoundPoints(tt.total)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestTotalIsMultipleTwentyFiveCentsPoints(t *testing.T) {
	var tests = []struct {
		name     string
		total    string
		expected int
	}{
		{
			"multiple",
			"12",
			25,
		},
		{
			"not multiple",
			"12.234",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getTotalIsMultipleTwentyFiveCentsPoints(tt.total)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestItemCountPoints(t *testing.T) {
	var tests = []struct {
		name     string
		items    []*models.Item
		expected int
	}{
		{
			"multiple items",
			[]*models.Item{
				{
					ShortDescription: "foo",
					Price:            "123",
				},
				{
					ShortDescription: "bar",
					Price:            "456",
				},
				{
					ShortDescription: "baz",
					Price:            "789",
				},
			},
			5,
		},
		{
			"no items",
			[]*models.Item{},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getPointsForEveryTwoItems(tt.items)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestItemDescriptionPoints(t *testing.T) {
	var tests = []struct {
		name     string
		items    []*models.Item
		expected int
	}{
		{
			"some descriptions",
			[]*models.Item{
				{
					ShortDescription: "foobar",
					Price:            "1",
				},
				{
					ShortDescription: "foo",
					Price:            "1",
				},
			},
			2,
		},
		{
			"not a multiple of 3",
			[]*models.Item{
				{
					ShortDescription: "fooba",
					Price:            "1",
				},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getDescriptionPoints(tt.items)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestDatePoints(t *testing.T) {
	var tests = []struct {
		name     string
		date     string
		expected int
	}{
		{
			"odd date",
			"2012-01-03",
			6,
		},
		{
			"even date",
			"2020-03-04",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getDatePoints(tt.date)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestTimePoints(t *testing.T) {
	var tests = []struct {
		name     string
		time     string
		expected int
	}{
		{
			"between 2 and 4",
			"15:00",
			10,
		},
		{
			"outside of 2 and 4",
			"01:00",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := getTimePoints(tt.time)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}

func TestComputePoints(t *testing.T) {
	var tests = []struct {
		name     string
		receipt  *models.Receipt
		expected int
	}{
		{
			"target",
			&models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []*models.Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					}, {
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					}, {
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					}, {
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					}, {
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
				Total: "35.35",
			},
			28,
		},
		{
			"M&M",
			&models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []*models.Item{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
				Total: "9.00",
			},
			109,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := computePoints(tt.receipt)
			if ans != tt.expected {
				t.Errorf("got %v, want %v", ans, tt.expected)
			}
		})
	}
}
