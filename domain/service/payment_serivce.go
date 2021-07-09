package service

import (
	"github.com/asveg/payment/domain/model"
	"github.com/asveg/payment/domain/repository"
)

type IPaymentDataService interface {
	AddPayment(payment *model.Payment) (paymentid int64, err error)
	DeletePaymentByID(paymentid int64) error
	UpdatePayment(payment *model.Payment) error
	FindPaymentByID(paymentid int64) (payment *model.Payment, err error)
	FindPaymentAll() ([]model.Payment, error)
}

type PaymentDataService struct {
	PaymentRepository repository.IPaymentRepostory
}

func NewPaymentDataService(paymentRepository repository.IPaymentRepostory) IPaymentDataService {
	return &PaymentDataService{PaymentRepository: paymentRepository}
}

func (p *PaymentDataService) AddPayment(payment *model.Payment) (paymentid int64, err error) {
	return p.PaymentRepository.CreatePayment(payment)
}

func (p *PaymentDataService) DeletePaymentByID(paymentid int64) error {
	return p.PaymentRepository.DeletePaymentByID(paymentid)
}

func (p *PaymentDataService) UpdatePayment(payment *model.Payment) error {
	return p.PaymentRepository.UpdatePayment(payment)
}

func (p *PaymentDataService) FindPaymentByID(paymentid int64) (payment *model.Payment, err error) {
	return p.PaymentRepository.FindPaymentByID(paymentid)
}

func (p *PaymentDataService) FindPaymentAll() ([]model.Payment, error) {
	return p.PaymentRepository.FindPaymentAllByID()
}
