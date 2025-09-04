package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
)

type IOrderRepository interface {
	GetOrderDetail(ctx context.Context, orderID string) ([]entity.OrderDetail, error)
	GetMonthlyVendorSales(ctx context.Context, month time.Time) ([]entity.VendorSalesReport, error)
}

type orderRepository struct {
	db *sql.DB
}

func (or *orderRepository) GetOrderDetail(ctx context.Context, orderID string) ([]entity.OrderDetail, error) {
	query := `
	SELECT
		o.order_id,
		o.created_at,
		u.full_name,
		u.email,
		pm.method_name AS payment_method,
		sa.address_line AS shipping_address,
		p.product_id,
		p.name AS product_name,
		c.name AS category_name,
		oi.quantity,
		oi.price
	FROM orders o
	JOIN users u ON o.user_id = u.user_id
	JOIN payment_methods pm ON o.payment_method_id = pm.payment_method_id
	JOIN shipping_addresses sa ON o.shipping_address_id = sa.shipping_address_id
	JOIN order_items oi ON o.order_id = oi.order_id
	JOIN products p ON oi.product_id = p.product_id
	JOIN categories c ON p.category_id = c.category_id
	WHERE o.order_id = $1;
	`

	rows, err := or.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order details: %w", err)
	}
	defer rows.Close()

	var details []entity.OrderDetail
	for rows.Next() {
		var od entity.OrderDetail
		if err := rows.Scan(
			&od.OrderID, &od.CreatedAt, &od.UserName, &od.UserEmail,
			&od.PaymentMethod, &od.ShippingAddress, &od.ProductID,
			&od.ProductName, &od.CategoryName, &od.Quantity, &od.Price,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		details = append(details, od)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return details, nil
}

func (or *orderRepository) GetMonthlyVendorSales(ctx context.Context, month time.Time) ([]entity.VendorSalesReport, error) {
	query := `
	WITH vendor_sales AS (
		SELECT
			v.vendor_id,
			v.name AS vendor_name,
			COUNT(DISTINCT o.order_id) AS total_orders,
			SUM(oi.price * oi.quantity) AS total_revenue,
			AVG(oi.quantity) AS avg_quantity_per_transaction
		FROM vendors v
		JOIN products p ON p.vendor_id = v.vendor_id
		JOIN order_items oi ON oi.product_id = p.product_id
		JOIN orders o ON o.order_id = oi.order_id
		WHERE DATE_TRUNC('month', o.created_at) = DATE_TRUNC('month', $1::date)
		GROUP BY v.vendor_id, v.name
	),
	top_products AS (
		SELECT
			v.vendor_id,
			p.name AS top_product,
			SUM(oi.quantity) AS total_quantity,
			ROW_NUMBER() OVER (PARTITION BY v.vendor_id ORDER BY SUM(oi.quantity) DESC) AS rn
		FROM vendors v
		JOIN products p ON p.vendor_id = v.vendor_id
		JOIN order_items oi ON oi.product_id = p.product_id
		JOIN orders o ON o.order_id = oi.order_id
		WHERE DATE_TRUNC('month', o.created_at) = DATE_TRUNC('month', $1::date)
		GROUP BY v.vendor_id, p.name
	)
	SELECT
		vs.vendor_id,
		vs.vendor_name,
		vs.total_orders,
		vs.total_revenue,
		vs.avg_quantity_per_transaction,
		tp.top_product
	FROM vendor_sales vs
	LEFT JOIN top_products tp ON vs.vendor_id = tp.vendor_id AND tp.rn = 1;
	`

	rows, err := or.db.QueryContext(ctx, query, month)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var reports []entity.VendorSalesReport
	for rows.Next() {
		var r entity.VendorSalesReport
		if err := rows.Scan(
			&r.VendorID,
			&r.VendorName,
			&r.TotalOrders,
			&r.TotalRevenue,
			&r.AvgQuantityPerTxn,
			&r.TopSellingProduct,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		reports = append(reports, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func NewOrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}
