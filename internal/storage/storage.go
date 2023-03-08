package storage

import (
	"github.com/pkg/errors"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"time"
)

type Storage struct {
	*gorm.DB
}

func NewStorage(masterConfig *DBConfig, replicaConfigs ...*DBConfig) *Storage {
	mainDB, err := gorm.Open(pg.Open(masterConfig.GetDSN()), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(replicaConfigs) != 0 && replicaConfigs[0].Host != "" {
		var replicas []gorm.Dialector
		for _, replicaConfig := range replicaConfigs {
			replicas = append(replicas, pg.Open(replicaConfig.GetDSN()))
		}
		resolverConfig := dbresolver.Config{
			Replicas: replicas,
		}
		resolver := dbresolver.Register(resolverConfig).
			SetConnMaxIdleTime(5 * time.Minute).
			SetMaxIdleConns(3).
			SetMaxOpenConns(10)
		if err := mainDB.Use(resolver); err != nil {
			log.Fatal(err)
		}
	}

	return &Storage{
		mainDB,
	}
}

func (storage *Storage) Migrate() error {
	if err := storage.DB.AutoMigrate(&UserAccount{}); err != nil {
		return errors.Wrap(err, "failed to migrate UserAccount")
	}
	if err := storage.DB.AutoMigrate(&UnauthorizedToken{}); err != nil {
		return errors.Wrap(err, "failed to migrate UnauthorizedToken")
	}
	if err := storage.DB.AutoMigrate(&FavoriteList{}); err != nil {
		return errors.Wrap(err, "failed to migrate FavoriteList")
	}
	if err := storage.DB.AutoMigrate(&FavoriteListItem{}); err != nil {
		return errors.Wrap(err, "failed to migrate FavoriteListItem")
	}
	if err := storage.DB.AutoMigrate(&SearchHistory{}); err != nil {
		return errors.Wrap(err, "failed to migrate SearchHistory")
	}
	if err := storage.DB.AutoMigrate(&SearchHistoryResult{}); err != nil {
		return errors.Wrap(err, "failed to migrate SearchHistoryResult")
	}
	return nil
}

func (storage *Storage) Close() {
	sqlDB, err := storage.DB.DB()
	if err == nil {
		sqlDB.Close()
	}
}
