package tasks

import (
	"context"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
)

// VTask represents voda-task
type VTask struct {
	code         string
	scheduleTime time.Time

	client *cloudtasks.Client
	ctx    context.Context
}

// task = {
// 	'http_request': {
// 		'http_method': 'POST',
// 		'url': CONFIG['CLOUD_TASKS']['FUNCTION_URL'],
// 		'oidc_token': {
// 			'service_account_email': CONFIG['CLOUD_TASKS']['SERVICE_ACCOUNT_EMAIL'],
// 		},
// 		'headers': {
// 			'Content-Type': 'application/json'
// 		},
// 		'body': json.dumps(payload).encode()
// 	},
// 	'schedule_time': timestamp
// }
