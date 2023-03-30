package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/email"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/util"
	"time"
)

var tenMinute = 10 * time.Minute

func (p *processor) ResetPasswordEmail(ctx context.Context, emailAddress string) (*gtsmodel.User, gtserror.WithCode) {
	if emailAddress == "" {
		return nil, gtserror.NewErrorNotFound(errors.New("no email provided"))
	}

	user, err := p.db.GetUserByEmailAddressFuzzy(ctx, emailAddress)
	if err != nil {
		if err == db.ErrNoEntries {
			return nil, gtserror.NewErrorNotFound(err)
		}
		return nil, gtserror.NewErrorInternalError(err)
	}

	// 获取最新的 user 对象，不从缓存中取
	if err = p.db.GetByID(context.Background(), user.ID, user); err != nil {
		return nil, gtserror.NewErrorForbidden(err)
	}

	// 十分钟内可以重新发送邮件
	//if user.ResetPasswordToken != "" && user.ResetPasswordSentAt.After(time.Now().Add(-tenMinute)) {
	//	// 十分钟内 token 有效
	//	log.Infof("%s, %s (十分钟内 token 有效)", user.ResetPasswordToken, user.ResetPasswordSentAt)
	//	return nil, gtserror.NewErrorForbidden(errors.New("ResetPasswordEmail: Repeated requests are not allowed within ten minutes"))
	//}

	// token := uuid.NewString()
	token := util.GenValidateCode(4)
	//link := uris.GenerateURIForEmailResetPassword(token)

	instance := &gtsmodel.Instance{}
	host := config.GetHost()
	if err := p.db.GetWhere(ctx, []db.Where{{Key: "domain", Value: host}}, instance); err != nil {
		return nil, gtserror.NewErrorForbidden(fmt.Errorf("ResetPasswordEmail: error getting instance: %s", err))
	}

	// assemble the email contents and send the email
	confirmData := email.ResetData{
		Username:     user.AccountID,
		InstanceURL:  instance.URI,
		InstanceName: instance.Title,
		//ResetLink:    link,
		ResetLink: token, // 只要验证码
	}
	if err := p.emailSender.SendResetEmail(user.Email, confirmData); err != nil {
		return nil, gtserror.NewErrorForbidden(fmt.Errorf("ResetPasswordEmail: error sending to email address %s belonging to user %s: %s", user.Email, user.AccountID, err))
	}

	updatingColumns := []string{"reset_password_token", "reset_password_sent_at", "last_emailed_at", "updated_at"}
	user.ResetPasswordToken = token
	user.ResetPasswordSentAt = time.Now()
	user.LastEmailedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := p.db.UpdateByID(ctx, user, user.ID, updatingColumns...); err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}

	return user, nil
}
