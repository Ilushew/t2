package services

import "strings"

func getVerificationEmailHTML(code string) string {
	return strings.ReplaceAll(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Код подтверждения</title>
    <style>
        @media only screen and (max-width: 600px) {
            .container { width: 100% !important; }
            .content { padding: 32px 20px !important; }
            .code-box { padding: 20px !important; }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #0a0a0a; font-family: Inter, -apple-system, BlinkMacSystemFont, Arial, sans-serif; -webkit-font-smoothing: antialiased;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="background-color: #0a0a0a;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" class="container" style="max-width: 600px; width: 100%; background-color: #141414; border-radius: 24px; overflow: hidden; border: 1px solid #2a2a2a; box-shadow: 0 20px 60px rgba(0,0,0,0.6);">
                    <tr>
                        <td align="center" style="background-color: #ff40d0; background: linear-gradient(135deg, #ff40d0 0%, #d600a0 100%); padding: 32px 20px;">
                            <h1 style="margin: 0; font-size: 22px; font-weight: 800; color: #ffffff; text-transform: uppercase; letter-spacing: 2.5px;">Маршруты по Удмуртии</h1>
                        </td>
                    </tr>
                    <tr>
                        <td class="content" style="padding: 40px 32px; background-color: #141414;">
                            <h2 style="margin: 0 0 20px; font-size: 20px; font-weight: 700; color: #ffffff; text-transform: uppercase; letter-spacing: 1px;">Код подтверждения</h2>
                            <p style="margin: 0 0 16px; color: #cccccc; font-size: 15px; line-height: 1.6;">Здравствуйте!</p>
                            <p style="margin: 0 0 28px; color: #cccccc; font-size: 15px; line-height: 1.6;">Ваш код для подтверждения регистрации:</p>
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%" style="margin-bottom: 28px;">
                                <tr>
                                    <td align="center" class="code-box" style="padding: 28px; background-color: #1e1e1e; border: 2px dashed #ff40d0; border-radius: 16px;">
                                        <span style="font-size: 34px; font-weight: 800; color: #ff40d0; letter-spacing: 8px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;">{{CODE}}</span>
                                    </td>
                                </tr>
                            </table>
                            <p style="margin: 0 0 12px; color: #999999; font-size: 14px; line-height: 1.6;">Код действителен <strong style="color: #ffffff; font-weight: 600;">15 минут</strong>.</p>
                            <p style="margin: 0 0 32px; color: #777777; font-size: 13px; line-height: 1.5;">Если вы не регистрировались — просто проигнорируйте это письмо.</p>
                            <div style="border-top: 1px solid #2a2a2a; margin: 24px 0;"></div>
                            <p style="margin: 0; color: #888888; font-size: 13px; line-height: 1.5;">
                                С уважением,<br>
                                <span style="color: #ff40d0; font-weight: 600;">Команда Udmurtia Route</span>
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td align="center" style="padding: 24px 20px; background-color: #0f0f0f; border-top: 1px solid #222222;">
                            <p style="margin: 0; color: #555555; font-size: 11px; font-family: Inter, Arial, sans-serif; letter-spacing: 0.5px;">© 2026 Udmurtia AI Route. Хакатон проект.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, "{{CODE}}", code)
}