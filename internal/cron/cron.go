package cron

import (
	"context"
	"time"

	"github.com/Tel3scop/chat-client/internal/config"
	"github.com/Tel3scop/chat-client/internal/connector/auth"
	"github.com/Tel3scop/helpers/logger"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

var (
	Scheduler gocron.Scheduler
)

type Cron struct {
	authClient *auth.Client
	config     *config.Config
}

// NewCron функция возвращает новый крон
func NewCron(
	authClient *auth.Client,
	config *config.Config,
) *Cron {
	return &Cron{
		authClient: authClient,
		config:     config,
	}
}

func (c *Cron) StartCron() {
	var err error
	Scheduler, err = gocron.NewScheduler()
	if err != nil {
		logger.Error("Failed to create scheduler", zap.Error(err))
		return
	}

	c.setRefreshToken()
	c.setAccessToken()

	Scheduler.Start()
	logger.Info("Scheduler started successfully", zap.Int("count_jobs", len(Scheduler.Jobs())))

	select {}
}

func StopCron() {
	logger.Info("StartCron: shutdown scheduler")
	err := Scheduler.Shutdown()
	if err != nil {
		logger.Error("StartCron: Error shutdown jobs", zap.Error(err))
	}
}

func (c *Cron) setRefreshToken() {
	fName := "set_refresh_token"
	job, err := Scheduler.NewJob(
		gocron.DurationJob(c.config.Encrypt.RefreshTokenExpiration),
		gocron.NewTask(
			func() {
				ctx := context.Background()
				refreshToken := auth.GetRefreshToken()
				logger.Info("starting at", zap.Time("start_at", time.Now()), zap.String("func", fName))
				if refreshToken == "" {
					logger.Info("refresh token not found, skip task", zap.Time("start_at", time.Now()), zap.String("func", fName))
					return
				}
				err := c.authClient.SetRefreshToken(ctx)
				if err != nil {
					logger.Error("executing func", zap.String("func", fName), zap.Error(err))
				}
				logger.Info("finished at", zap.Time("end_at", time.Now()), zap.String("func", fName))
			},
		),
	)
	if err != nil {
		logger.Error("can not create new Job", zap.Error(err))
	}
	logger.Info("new Job", zap.String("func", fName), zap.Uint32("id", job.ID().ID()))
}

func (c *Cron) setAccessToken() {
	fName := "set_access_token"
	job, err := Scheduler.NewJob(
		gocron.DurationJob(c.config.Encrypt.AccessTokenExpiration),
		gocron.NewTask(
			func() {
				ctx := context.Background()
				refreshToken := auth.GetRefreshToken()
				logger.Info("starting at", zap.Time("start_at", time.Now()), zap.String("func", fName))
				if refreshToken == "" {
					logger.Info("refresh token not found, skip task", zap.Time("start_at", time.Now()), zap.String("func", fName))
					return
				}
				err := c.authClient.SetAccessToken(ctx)
				if err != nil {
					logger.Error("executing func", zap.String("func", fName), zap.Error(err))
				}
				logger.Info("finished at", zap.Time("end_at", time.Now()), zap.String("func", fName))
			},
		),
	)
	if err != nil {
		logger.Error("can not create new Job", zap.Error(err))
	}
	logger.Info("new Job", zap.String("func", fName), zap.Uint32("id", job.ID().ID()))
}
