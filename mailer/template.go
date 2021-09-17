package mailer

var taskNotificationEmailTemplate = `
<h1>DozenPlans %s</h1>
<p>您好，现在到了您设置的任务:%s的提醒时间了</p>
`
var taskFailedEmailTemplate = `
<h1>DozenPlans 任务超时提醒</h1>
<p>您好，您设置的任务 %s 已经超时未完成，任务失败，请注意任务进度。</p>
<p>任务内容 %s</p>`
