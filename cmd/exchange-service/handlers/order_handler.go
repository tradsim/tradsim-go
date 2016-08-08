package handlers

import (
	"github.com/mantzas/adaptlog"
)

// OrderHandler handles orders
type OrderHandler struct {
	logger adaptlog.LevelLogger
	//processor trading.Processor
}

// // OrderCreateHandle is the handler for the orders
// func (oh *OrderHandler) OrderCreateHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

// 	var orderCreate models.OrderCreate

// 	err := json.NewDecoder(r.Body).Decode(&orderCreate)

// 	if err != nil {

// 		oh.logger.Errorf("Failed to bind model! %s", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	price, err := models.NewPrice(orderCreate.Price, uint8(2))

// 	if err != nil {
// 		oh.logger.Errorf("Failed to getting price! %s", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	direction, err := models.TradeDirectionFromString(orderCreate.Direction)
// 	if err != nil {
// 		oh.logger.Errorf("Failed to getting trade direction! %s", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	orderID, err := uuid.FromString(orderCreate.ID)
// 	if err != nil {
// 		oh.logger.Errorf("Failed to getting order id! %s", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	order := trading.NewOrder(orderID, orderCreate.Symbol, *price, orderCreate.Quantity, direction)

// 	err = oh.processor.Process(order)
// 	if err != nil {
// 		oh.logger.Errorf("Failed to process order! %s", err)
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusAccepted)
// }
