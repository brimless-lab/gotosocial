package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/log"
	"time"
)

func (p *processor) VerifyResetPasswordToken(ctx context.Context, email string, resetPasswordToken string) gtserror.WithCode {
	log.Info("VerifyResetPasswordToken start")
	user, err := p.db.GetUserByEmailAddressFuzzy(ctx, email)
	if err != nil {
		if err == db.ErrNoEntries {
			return gtserror.NewErrorNotFound(err)
		}
		return gtserror.NewErrorInternalError(err)
	}

	// 获取最新的 user 对象，不从缓存中取
	if err = p.db.GetByID(context.Background(), user.ID, user); err != nil {
		return gtserror.NewErrorForbidden(err)
	}

	if user.ResetPasswordToken != resetPasswordToken || (user.ResetPasswordToken == resetPasswordToken && time.Now().After(user.ResetPasswordSentAt.Add(tenMinute))) {
		// token 无效 或 过期
		log.Infof("VerifyResetPasswordToken token: %s %s 无效 或 过期", user.ResetPasswordToken, user.ResetPasswordSentAt)
		return gtserror.NewErrorForbidden(errors.New("VerifyResetPasswordToken: No valid verification code"))
	}

	if !user.Account.SuspendedAt.IsZero() {
		log.Infof("VerifyResetPasswordToken: account %s is suspended", user.AccountID)
		return gtserror.NewErrorForbidden(fmt.Errorf("VerifyResetPasswordToken: account %s is suspended", user.AccountID))
	}

	updatingColumns := []string{"reset_password_token", "reset_password_sent_at", "last_emailed_at", "updated_at"}
	user.ResetPasswordToken = ""
	user.ResetPasswordSentAt = time.Time{} // 赋零值
	user.LastEmailedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := p.db.UpdateByID(ctx, user, user.ID, updatingColumns...); err != nil {
		log.Info("VerifyResetPasswordToken UpdateByID error")
		return gtserror.NewErrorInternalError(err)
	}
	return nil
}
