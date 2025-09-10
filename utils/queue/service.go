package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"senkou-catalyst-be/utils/config"
	"senkou-catalyst-be/utils/mailer"
	"time"

	"github.com/hibiken/asynq"
)

type QueueService struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	inspector *asynq.Inspector
	handlers  map[string]asynq.HandlerFunc
}

type QueueConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
	Queues        map[string]int
}

func DefaultQueueConfig() *QueueConfig {
	host := config.GetEnv("REDIS_HOST", "localhost:6379")
	port := config.GetEnv("REDIS_PORT", "6379")
	password := config.GetEnv("REDIS_PASSWORD", "")
	db := config.GetEnvAsInt("REDIS_DB", 0)

	return &QueueConfig{
		RedisHost:     host,
		RedisPort:     port,
		RedisPassword: password,
		RedisDB:       db,
		Concurrency:   10,
		Queues: map[string]int{
			"critical": 6,
			"high":     4,
			"default":  3,
			"low":      1,
		},
	}
}

func NewQueueService(config *QueueConfig) (*QueueService, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency:         config.Concurrency,
		Queues:              config.Queues,
		LogLevel:            asynq.InfoLevel,
		RetryDelayFunc:      asynq.DefaultRetryDelayFunc,
		HealthCheckInterval: 15 * time.Second,
	})

	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
		LogLevel: asynq.InfoLevel,
	})

	inspector := asynq.NewInspector(redisOpt)

	return &QueueService{
		client:    client,
		server:    server,
		scheduler: scheduler,
		inspector: inspector,
		handlers:  make(map[string]asynq.HandlerFunc),
	}, nil
}

func (qs *QueueService) NewJobBuilder(taskName string) *JobBuilder {
	return NewJobBuilder(qs.client, taskName)
}

func (qs *QueueService) RegisterHandler(taskName string, handler asynq.HandlerFunc) {
	qs.handlers[taskName] = handler
}

func (qs *QueueService) RegisterHandlerFunc(taskName string, handler func(context.Context, *asynq.Task) error) {
	qs.handlers[taskName] = asynq.HandlerFunc(handler)
}

func (qs *QueueService) RegisterEmailHandlers() {
	qs.RegisterHandlerFunc("email:send_activation", qs.handleSendActivationEmail)
}

func (qs *QueueService) Start() error {
	mux := asynq.NewServeMux()
	for taskName, handler := range qs.handlers {
		mux.HandleFunc(taskName, handler)
	}

	go func() {
		if err := qs.scheduler.Run(); err != nil {
			log.Printf("Failed to start scheduler: %v", err)
		}
	}()

	return qs.server.Run(mux)
}

func (qs *QueueService) Stop() {
	qs.scheduler.Shutdown()
	qs.server.Shutdown()
	if err := qs.client.Close(); err != nil {
		log.Printf("Failed to close queue client: %v", err)
	}
}

func (qs *QueueService) GetQueueNames() ([]string, error) {
	return qs.inspector.Queues()
}

func (qs *QueueService) GetQueueInfo(qname string) (*asynq.QueueInfo, error) {
	return qs.inspector.GetQueueInfo(qname)
}

func (qs *QueueService) GetTaskInfo(queue, taskID string) (*asynq.TaskInfo, error) {
	return qs.inspector.GetTaskInfo(queue, taskID)
}

func (qs *QueueService) SchedulePeriodicTask(cronspec, taskName string, payload map[string]interface{}, opts ...asynq.Option) (string, error) {
	task := qs.NewJobBuilder(taskName).WithPayload(payload)
	payloadBytes, err := task.marshalPayload()
	if err != nil {
		return "", err
	}

	asynqTask := asynq.NewTask(taskName, payloadBytes, opts...)

	entryID, err := qs.scheduler.Register(cronspec, asynqTask)
	if err != nil {
		return "", err
	}

	return entryID, nil
}

func (qs *QueueService) DeletePeriodicTask(entryID string) error {
	return qs.scheduler.Unregister(entryID)
}

func (jb *JobBuilder) marshalPayload() ([]byte, error) {
	return jb.marshalPayloadInternal()
}

func (jb *JobBuilder) marshalPayloadInternal() ([]byte, error) {
	if len(jb.payload) == 0 {
		return []byte("{}"), nil
	}

	payloadBytes, err := json.Marshal(jb.payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

// handleSendActivationEmail handles sending activation emails
func (qs *QueueService) handleSendActivationEmail(ctx context.Context, task *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal activation email payload: %w", err)
	}

	email, ok := payload["email"].(string)
	if !ok {
		return fmt.Errorf("invalid email in payload")
	}

	activationLink, ok := payload["activation_link"].(string)
	if !ok {
		return fmt.Errorf("invalid activation_link in payload")
	}

	supportEmail, ok := payload["support_email"].(string)
	if !ok {
		return fmt.Errorf("invalid support_email in payload")
	}

	userName, _ := payload["user_name"].(string)

	// Create mailer service and send email
	mailerService, err := mailer.NewMailerService()
	if err != nil {
		return fmt.Errorf("failed to initialize mailer service: %w", err)
	}

	if !mailerService.TemplateExists("account-activation.html") {
		return fmt.Errorf("email template not found: account-activation.html")
	}

	templateData := map[string]interface{}{
		"ActivationLink": activationLink,
		"SupportEmail":   supportEmail,
		"UserName":       userName,
	}

	err = mailerService.SendTemplate(
		email,
		"Catalyst - Account Activation",
		"account-activation.html",
		templateData,
	)

	if err != nil {
		return fmt.Errorf("failed to send activation email to %s: %w", email, err)
	}

	log.Printf("Successfully sent activation email to %s", email)
	return nil
}
