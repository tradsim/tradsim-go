package data

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/models"
	commonmodels "github.com/tradsim/tradsim-go/models"
)

// OrderRepository interface
type OrderRepository interface {
	GetOrders() ([]models.Order, error)
	GetPositions() ([]models.Position, error)
}

// OrderRepositoryImpl order repository
type OrderRepositoryImpl struct {
	db           *sql.DB
	orderStmt    *sql.Stmt
	logStmt      *sql.Stmt
	tradeStmt    *sql.Stmt
	positionStmt *sql.Stmt
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) (*OrderRepositoryImpl, error) {

	orderStmt, err := db.Prepare(`SELECT "order".* FROM public."order"`)
	if err != nil {
		return nil, err
	}

	logStmt, err := db.Prepare(`SELECT * FROM order_log`)
	if err != nil {
		return nil, err
	}

	tradeStmt, err := db.Prepare(`SELECT * FROM trade`)
	if err != nil {
		return nil, err
	}

	positionStmt, err := db.Prepare(`SELECT * FROM position`)
	if err != nil {
		return nil, err
	}

	return &OrderRepositoryImpl{db, orderStmt, logStmt, tradeStmt, positionStmt}, nil
}

// GetOrders returns all orders from storage
func (or *OrderRepositoryImpl) GetOrders() ([]models.Order, error) {

	orders, err := or.getOrders()
	if err != nil {
		return nil, err
	}

	logs, err := or.getOrderLogs()
	if err != nil {
		return nil, err
	}

	trades, err := or.getTrades()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orders); i++ {

		for _, log := range logs {
			if log.OrderID != orders[i].ID {
				continue
			}
			orders[i].Logs = append(orders[i].Logs, log)
		}

		for _, trade := range trades {
			if trade.OrderID != orders[i].ID {
				continue
			}
			orders[i].Trades = append(orders[i].Trades, trade)
		}
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) getOrders() ([]models.Order, error) {
	rows, err := or.orderStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	var id uuid.UUID
	var symbol string
	var price float64
	var quantity uint
	var direction commonmodels.TradeDirection
	var status commonmodels.OrderStatus
	var created time.Time

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *models.NewOrder(id, symbol, price, quantity, direction, status, created))
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepositoryImpl) getOrderLogs() ([]models.OrderLog, error) {
	rows, err := or.logStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.OrderLog

	var action string
	var occured time.Time
	var orderID uuid.UUID
	var id int64

	for rows.Next() {
		err := rows.Scan(&id, &orderID, &action, &occured)
		if err != nil {
			return nil, err
		}

		logs = append(logs, models.OrderLog{id, orderID, action, occured})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (or *OrderRepositoryImpl) getTrades() ([]models.Trade, error) {
	rows, err := or.tradeStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	var price float64
	var quantity uint
	var occured time.Time
	var orderID uuid.UUID
	var id int64

	for rows.Next() {
		err := rows.Scan(&id, &orderID, &price, &quantity, &occured)
		if err != nil {
			return nil, err
		}

		trades = append(trades, models.Trade{id, orderID, price, quantity, occured})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return trades, nil
}

// GetPositions returns the positions
func (or *OrderRepositoryImpl) GetPositions() ([]models.Position, error) {

	rows, err := or.positionStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []models.Position
	var id int64
	var symbol string
	var quantity int
	var updated time.Time

	for rows.Next() {
		err := rows.Scan(&id, &symbol, &quantity, &updated)
		if err != nil {
			return nil, err
		}

		positions = append(positions, models.Position{id, symbol, quantity, updated})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return positions, nil
}
