package notifier

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	userDAO "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	NotifyProblemWF = "NotifyProblem"
	NotifierQueue   = "NotifierQueue"
)

const (
	notificationEmailNoAuthKey = `
		This email is automatically sent to every admin of the Welbridge system to notify them of 
		a certain operation problem.

		[Timestamp=%s] %s side's Authenticator key is currently not available in
		Welbridge [core] microservice, most likely due to a service restart during operation.
		Without this key, claim request signing wouldn't be available and users of the bridge
		wouldn't be able to claim their resources. Please visit Welbridge's administration
		portal at %s, login as your admin handle, and provide the system with the
		Authenticator key at %s.
	`
	notificationSubjectNoAuthKey = "[Welbridge system] No Authenticator key available for %s side"
)

type Notifier struct {
	tempCli client.Client
	userDAO userDAO.IUserDAO
	mailer  *manager.Mailer
	worker  worker.Worker
}

func MkNotifier(tempCli client.Client, daos *dao.DAOs, mailer *manager.Mailer) *Notifier {
	return &Notifier{tempCli: tempCli, userDAO: daos.User, mailer: mailer}
}

func (notifier *Notifier) NotifyWorkflow(ctx workflow.Context, problem error, role string) error {
	log := workflow.GetLogger(ctx)
	log.Info("[NotifierWF] start notifying role " + role + " of problem: " + problem.Error())
	ao := workflow.ActivityOptions{
		TaskQueue:              NotifierQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 100,
			MaximumAttempts: 32,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteLocalActivity(ctx, notifier.NotifyActivity, problem, role)
	if err := future.Get(ctx, nil); err != nil {
		log.Error("[NotifierWF] Failed to notify users with role " + role + " of problem: " + problem.Error())
		return err
	}

	log.Info("[NotifierWF] Notification succeeded")

	return nil
}

func (notifier *Notifier) NotifyActivity(ctx context.Context, problem error, role string) error {
	logger.Get().Info().Msgf("[Notifier] Problem to notify all users with role %s: %s", role, problem.Error())

	if err := notifier.handleHealthProblem(problem, role); err != nil {
		logger.Get().Err(err).Msgf("[Notifier] Failed to notify users with role %s about problem: %s", role, problem.Error())
		return err
	}

	logger.Get().Info().Msgf("[Notifier] Notified users with role %s about problem: %s", role, problem.Error())
	return nil
}

func (notifier *Notifier) handleHealthProblem(prob error, role string) error {
	switch prob {
	case model.ErrEthAuthenticatorKeyUnavailable:
		adminPortal := "<placeholder>"
		chain := "Ethereum"
		injectKeyURL := adminPortal + "/<placeholder>"
		mailBody := fmt.Sprintf(notificationEmailNoAuthKey, time.Now().String(), chain, adminPortal, injectKeyURL)
		subject := fmt.Sprint(notificationSubjectNoAuthKey, chain)
		return notifier.sendNotificationToRole(role, subject, mailBody)
	case model.ErrWelAuthenticatorKeyUnavailable:
		adminPortal := "<placeholder>"
		chain := "Welups"
		injectKeyURL := adminPortal + "/<placeholder>"
		mailBody := fmt.Sprintf(notificationEmailNoAuthKey, time.Now().String(), chain, adminPortal, injectKeyURL)
		subject := fmt.Sprint(notificationSubjectNoAuthKey, chain)
		return notifier.sendNotificationToRole(role, subject, mailBody)
	default:
		return nil
	}
	return nil
}

func (notifier *Notifier) sendNotificationToRole(role string, subject string, body string) error {
	users, err := notifier.userDAO.GetUsersWithRole(role, 0, 1000)
	log := logger.Get()
	if err != nil {
		log.Err(err).Msgf("[Notifier] Unable to fetch users with role %s", role)
		return err
	}
	fmt.Println("Users: ", users)

	mails := libs.Map(func(u model.User) string { return u.Email },
		libs.Filter(func(u model.User) bool { return u.Status == "ok" }, users))
	fmt.Println("Mails: ", mails)
	for _, mail := range mails {
		mess := notifier.mailer.MkPlainMessage(mail, subject, body)
		err := notifier.mailer.Send(mess)
		if err != nil {
			log.Err(err).Msgf("[Notifier] unable to send mail to address %s", mail) // best effort lol
		}
	}

	return nil
}

// Worker
func (notifier *Notifier) registerService(w worker.Worker) {
	w.RegisterActivity(notifier.NotifyActivity)

	w.RegisterWorkflowWithOptions(notifier.NotifyWorkflow, workflow.RegisterOptions{Name: NotifyProblemWF})
}

func (notifier *Notifier) StartService() error {
	w := worker.New(notifier.tempCli, NotifierQueue, worker.Options{})
	notifier.registerService(w)

	notifier.worker = w
	logger.Get().Info().Msgf("Starting Notifier")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting Notifier")
		return err
	}

	logger.Get().Info().Msgf("Notifier started")
	return nil
}

func (notifier *Notifier) StopService() {
	if notifier.worker != nil {
		notifier.worker.Stop()
	}
}
