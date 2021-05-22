package repository

import (
	"github.com/jinzhu/gorm"
	"payment/domain/model"
)

type IPaymentRepostory interface {
	InitTable()error
	FindPaymentByID(paymentid int64)(payment *model.Payment,err error)
	CreatePayment(payment *model.Payment)(paymentid int64, err error)
	DeletePaymentByID(paymentid int64) error
	UpdatePayment(payment *model.Payment)error
	FindPaymentAllByID()(paymentAll []model.Payment, err error)
}

type PaymentRepository struct {
	mysql *gorm.DB
}

func NewPaymentRepository (db *gorm.DB) IPaymentRepostory {
	return &PaymentRepository{mysql: db}
}

func (p *PaymentRepository) InitTable() error {
	return p.mysql.CreateTable(&model.Payment{}).Error
}

func (p *PaymentRepository) FindPaymentByID(paymentid int64) (payment *model.Payment, err error) {
	payment = &model.Payment{}
	return payment, p.mysql.First(payment,paymentid).Error
}

func (p *PaymentRepository) CreatePayment(payment *model.Payment) (paymentid int64, err error) {
	return paymentid,p.mysql.Create(payment).Error
}

func (p *PaymentRepository) DeletePaymentByID(paymentid int64) error {
	return p.mysql.Where("id=?",paymentid).Delete(paymentid).Error
}

func (p *PaymentRepository) UpdatePayment(payment *model.Payment) error {
	return p.mysql.Model(payment).UpdateColumn(payment).Error
}

func (p *PaymentRepository) FindPaymentAllByID() (paymentAll []model.Payment, err error) {
	return paymentAll,p.mysql.Find(paymentAll).Error
}
