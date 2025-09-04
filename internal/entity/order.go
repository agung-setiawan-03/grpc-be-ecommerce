package entity

import "time"

type OrderDetail struct {
	OrderID         string
	CreatedAt       time.Time
	UserName        string
	UserEmail       string
	PaymentMethod   string
	ShippingAddress string
	ProductID       string
	ProductName     string
	CategoryName    string
	Quantity        int
	Price           float64
}

type VendorSalesReport struct {
	VendorID          string
	VendorName        string
	TotalOrders       int
	TotalRevenue      float64
	AvgQuantityPerTxn float64
	TopSellingProduct string
}
