package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"infotecs-test-task/internal/lib/config"
	"infotecs-test-task/internal/models"
	"infotecs-test-task/internal/repositories"
	"infotecs-test-task/internal/storage"
)

type WalletProvider interface {
	repositories.TX
	InsertOne(wallet models.Wallet) error
	UpdateBalance(id string, delta float32) error
	FindOneByID(id string) (models.Wallet, error)
}

type TransactionsProvider interface {
	repositories.TX
	InsertOne(transaction models.Transaction) error
	FindAllByWalletID(id string) ([]models.Transaction, error)
}

type WalletService struct {
	wProvider WalletProvider
	tProvider TransactionsProvider
	cfg       config.API
}

func NewWalletService(cfg config.API, wProvider WalletProvider, tProvider TransactionsProvider) *WalletService {
	return &WalletService{
		wProvider: wProvider,
		tProvider: tProvider,
		cfg:       cfg,
	}
}

func (s *WalletService) CreateWallet() (models.Wallet, error) {
	wallet := models.Wallet{
		Id:      uuid.NewString(),
		Balance: s.cfg.EWalletInitBalance,
	}
	err := s.wProvider.InsertOne(wallet)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(id string) (models.Wallet, error) {
	return s.wProvider.FindOneByID(id)
}

func (s *WalletService) TransferMoney(transaction models.Transaction) error {
	tx, err := s.wProvider.GetTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	s.tProvider.SetTx(tx)

	err = s.tProvider.InsertOne(transaction)
	if err != nil {
		return err
	}

	err = s.wProvider.UpdateBalance(transaction.From, -transaction.Amount)
	if err != nil {
		if errors.Is(err, storage.ErrWalletNotFound) {
			return ErrSendingWalletNotFound
		}
		return err
	}

	err = s.wProvider.UpdateBalance(transaction.To, transaction.Amount)
	if err != nil {
		if errors.Is(err, storage.ErrWalletNotFound) {
			return ErrReceivingWalletNotFound
		}
		return err
	}

	tx.Commit()

	s.wProvider.SetTx(nil)
	s.tProvider.SetTx(nil)

	return nil
}

func (s *WalletService) GetTransactionsHistory(id string) ([]models.Transaction, error) {
	transactions, err := s.tProvider.FindAllByWalletID(id)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
