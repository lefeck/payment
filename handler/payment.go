package handler

import (
	"context"
	"github.com/wangjinh/common"
	"github.com/wangjinh/payment/domain/model"
	"github.com/wangjinh/payment/domain/service"

	. "github.com/wangjinh/payment/proto/payment"
)

type Payment struct {
	PaymentDataService service.IPaymentDataService
}

/*
type PaymentHandler interface {
	//支付的方法
	AddPayment(context.Context, *PaymentInfo, *PaymentID) error
	//更新支付
	UpdatePayment(context.Context, *PaymentInfo, *Response) error
	DeletePaymentByID(context.Context, *PaymentID, *Response) error
	FindPaymentByID(context.Context, *PaymentID, *PaymentInfo) error
	FindAllPayment(context.Context, *All, *PaymentAll) error
}
*/

func (p *Payment) AddPayment(ctx context.Context, request *PaymentInfo, response *PaymentID) error {
	payment := &model.Payment{}

	if err := common.SwapTo(request, payment); err != nil {
		return err
	}
	paymentId, err := p.PaymentDataService.AddPayment(payment)
	if err != nil {
		return err
	}
	response.PaymentId = paymentId
	return nil
}

func (p *Payment) UpdatePayment(ctx context.Context, request *PaymentInfo, response *Response) error {
	payment := &model.Payment{}

	if err := common.SwapTo(request, payment); err != nil {
		return err
	}
	if err := p.PaymentDataService.UpdatePayment(payment); err != nil {
		return err
	}

	return nil
}

func (p *Payment) DeletePaymentByID(ctx context.Context, request *PaymentID, response *Response) error {
	if err := p.PaymentDataService.DeletePaymentByID(request.PaymentId); err != nil {
		return err
	}
	response.Msg = "Payment delete success"
	return nil
}

func (p *Payment) FindPaymentByID(ctx context.Context, request *PaymentID, response *PaymentInfo) error {
	paymentId, err := p.PaymentDataService.FindPaymentByID(request.PaymentId)
	if err != nil {
		return err
	}
	err = common.SwapTo(paymentId, response)
	if err != nil {
		return err
	}
	return nil
}

func (p *Payment) FindAllPayment(ctx context.Context, request *All, response *PaymentAll) error {
	paymentAll, err := p.PaymentDataService.FindPaymentAll()
	if err != nil {
		return err
	}
	for _, v := range paymentAll {
		payment := &PaymentInfo{}
		if err := common.SwapTo(request, v); err != nil {
			return err
		}
		response.PaymentInfo = append(response.PaymentInfo, payment)
	}
	return nil
}
