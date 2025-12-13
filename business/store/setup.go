package store

import "github.com/sirupsen/logrus"

func SetUpMYSQL(logger *logrus.Logger) (Store, Authenticator, error) {
	cfg := NewENVConfig()
	db, err := ConnectMYSQL(cfg.DB)
	if err != nil {
		return nil, nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}
	err = Migrate(db, cfg.DB.Database, cfg.DB.MigrationPath)
	if err != nil {
		return nil, nil, err
	}

	return NewDB(db), &Auth{
		Conn:   db,
		Logger: logger,
	}, nil
}
