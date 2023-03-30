/*
   GoToSocial
   Copyright (C) 2021-2023 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/superseriousbusiness/gotosocial/internal/processing"
)

const (
	// Params
	Email = "email"
	Token = "reset_password_token"

	// BasePath is the base URI path for this module, minus the 'api' prefix
	BasePath = "/v1/user"
	// PasswordChangePath is the path for POSTing a password change request.
	PasswordChangePath = BasePath + "/password_change"
	// 发送修改用户密码的验证码
	SendResetPasswordEmailPath = BasePath + "/send_reset_password_email"
	// 验证修改密码的邮箱和验证码
	VerifyResetPasswordToken = BasePath + "/verify_reset_password_token"
)

type Module struct {
	processor processing.Processor
}

func New(processor processing.Processor) *Module {
	return &Module{
		processor: processor,
	}
}

func (m *Module) Route(attachHandler func(method string, path string, f ...gin.HandlerFunc) gin.IRoutes) {
	attachHandler(http.MethodPost, PasswordChangePath, m.PasswordChangePOSTHandler)
	attachHandler(http.MethodPost, SendResetPasswordEmailPath, m.SendResetPasswordEmailPostHandler)
	attachHandler(http.MethodPost, VerifyResetPasswordToken, m.VerifyResetTokenPostHandler)
}
