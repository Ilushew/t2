package services

import "fmt"

func getVerificationEmailHTML(code string) string {
	return fmt.Sprintf(`
			<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: #4F46E5; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
				.content { background: #f9f9f9; padding: 30px; border: 1px solid #ddd; }
				.code { font-size: 32px; font-weight: bold; color: #4F46E5; text-align: center; padding: 20px; background: white; border: 2px dashed #4F46E5; margin: 20px 0; border-radius: 8px; }
				.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Udmurtia AI Route</h1>
				</div>
				<div class="content">
					<h2>Код подтверждения</h2>
					<p>Здравствуйте!</p>
					<p>Ваш код для подтверждения регистрации:</p>
					<div class="code">%s</div>
					<p>Код действителен <strong>15 минут</strong>.</p>
					<p>Если вы не регистрировались — просто проигнорируйте это письмо.</p>
				</div>
				<div class="footer">
					<p>© 2026 Udmurtia AI Route. Хакатон проект.</p>
				</div>
			</div>
		</body>
		</html>
	`, code)
}
