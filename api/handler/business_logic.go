package handler

import (
	"context"
	"fmt"
	"market_system/models"
	"market_system/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SaleScanBarcode(c *gin.Context) {

	var (
		saleID   = c.Query("sale_id")
		branchID = c.Query("branch_id")
		barcode  = c.Query("barcode")
	)
	if !helpers.IsValidUUID(saleID) {
		handleResponse(c, http.StatusBadRequest, "sale id is not uuid")
		return
	}

	if !helpers.IsValidUUID(branchID) {
		handleResponse(c, http.StatusBadRequest, "branch id is not uuid")
		return
	}

	remainingTableProduct, err := h.strg.Remainder().GetList(context.Background(), &models.GetListRemainderRequest{
		Limit: 1,
		Query: fmt.Sprintf(" AND bracode = %s AND branch_id = %s", barcode, branchID),
	})

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(remainingTableProduct.Remainder) <= 0 {
		handleResponse(c, http.StatusBadRequest, "Товар не найден")
		return
	}

	saleProduct, err := h.strg.Sale_Product().GetList(context.Background(), &models.GetListSaleProductRequest{
		Limit: 1,
		Query: fmt.Sprintf(" AND bracode = %s AND sale_id = %s", barcode, saleID),
	})

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(saleProduct.SaleProducts) <= 0 {
		var product = remainingTableProduct.Remainder[0]
		_, err := h.strg.Sale_Product().Create(context.Background(), &models.CreateSaleProduct{
			SaleID:            saleID,
			CategoryID:        product.CategoryID,
			ProductName:       product.ProductName,
			Barcode:           product.Barcode,
			RemainingQuantity: product.Quantity,
			Quantity:          1,
			AllowDiscount:     false,
			DiscountType:      "",
			Discount:          0,
			Price:             product.PriceIncome,
			TotalAmount:       product.PriceIncome,
		})
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		handleResponse(c, http.StatusCreated, "Успешно")
		return
	}

	if remainingTableProduct.Remainder[0].Quantity-saleProduct.SaleProducts[0].Quantity < 0 {
		handleResponse(c, http.StatusBadRequest, "Максималный лимит")
		return
	}

	_, err = h.strg.Sale_Product().Update(context.Background(), &models.UpdateSaleProduct{
		Id:                saleProduct.SaleProducts[0].Id,
		RemainingQuantity: saleProduct.SaleProducts[0].RemainingQuantity,
		Quantity:          saleProduct.SaleProducts[0].Quantity + 1,
		AllowDiscount:     saleProduct.SaleProducts[0].AllowDiscount,
		DiscountType:      saleProduct.SaleProducts[0].DiscountType,
		Discount:          saleProduct.SaleProducts[0].Discount,
		Price:             saleProduct.SaleProducts[0].Price,
		TotalAmount:       saleProduct.SaleProducts[0].TotalAmount + saleProduct.SaleProducts[0].Price,
	})

	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, http.StatusCreated, "Успешно")
	return
}

func (h *Handler) Dosale(c *gin.Context) {
	var (
		saleID   = c.Query("sale_id")
		branchID = c.Query("branch_id")
	)
	if !helpers.IsValidUUID(saleID) {
		handleResponse(c, http.StatusBadRequest, "sale id is not uuid")
		return
	}

	if !helpers.IsValidUUID(branchID) {
		handleResponse(c, http.StatusBadRequest, "branch id is not uuid")
		return
	}

	saleData, err := h.strg.Sale().GetByID(context.Background(), &models.SalePrimaryKey{Id: saleID})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	salePaymentResponse, err := h.strg.Payment().GetList(context.Background(), &models.GetListPaymentRequest{
		Limit: 100,
		Query: fmt.Sprintf(" AND sale_id = %s", saleID),
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	cashTransactionResponse, err := h.strg.Transaction().GetList(context.Background(), &models.GetListTransactonRequest{
		Limit: 100,
		Query: fmt.Sprintf(" AND shift_id = %s", saleData.ShiftID),
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(cashTransactionResponse.Transactions) <= 0 {
		handleResponse(c, http.StatusBadRequest, "не найден транзакции")
		return
	}

	var (
		salePayment     = salePaymentResponse.Payments[0]
		cashTransaction = cashTransactionResponse.Transactions[0]
	)
	_, err = h.strg.Transaction().Update(context.Background(), &models.UpdateTransaction{
		Id:          cashTransaction.Id,
		Cash:        cashTransaction.Cash + salePayment.Cash,
		Uzcard:      cashTransaction.Uzcard + salePayment.Uzcard,
		Payme:       cashTransaction.Payme + salePayment.Payme,
		Click:       cashTransaction.Click + salePayment.Click,
		Humo:        cashTransaction.Humo + salePayment.Humo,
		Apelsin:     cashTransaction.Apelsin + salePayment.Apelsin,
		TotalAmount: cashTransaction.TotalAmount + salePayment.TotalAmount,
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	saleProductResponse, err := h.strg.Sale_Product().GetList(context.Background(), &models.GetListSaleProductRequest{
		Limit: 1000,
		Query: fmt.Sprintf(" AND sale_id = %s", saleID),
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var remainderQuery = fmt.Sprintf(" AND branch_id = %s AND barcode IN (", branchID)
	for _, saleProduct := range saleProductResponse.SaleProducts {
		remainderQuery += fmt.Sprintf("'%s',", saleProduct.Barcode)
	}
	remainderQuery = remainderQuery[:len(remainderQuery)-1]
	remainderQuery += ")"

	remainderResponse, err := h.strg.Remainder().GetList(context.Background(), &models.GetListRemainderRequest{
		Limit: 1000,
		Query: remainderQuery,
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, saleProduct := range saleProductResponse.SaleProducts {
		for _, remainder := range remainderResponse.Remainder {
			if saleProduct.Barcode == remainder.Barcode {
				_, err := h.strg.Remainder().Update(context.Background(), &models.UpdateRemainder{
					Id:          remainder.Id,
					ProductName: remainder.ProductName,
					Barcode:     remainder.Barcode,
					PriceIncome: remainder.PriceIncome,
					Quantity:    remainder.Quantity - saleProduct.Quantity,
				})
				if err != nil {
					handleResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
	}

	_, err = h.strg.Sale().Update(context.Background(), &models.UpdateSale{
		Id:          saleData.Id,
		BranchID:    saleData.BranchID,
		SalePointID: saleData.SalePointID,
		ShiftID:     saleData.ShiftID,
		EmployeeID:  saleData.EmployeeID,
		Barcode:     saleData.Barcode,
		Status:      "finished",
	})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, http.StatusCreated, "Успешно")
	return
}
