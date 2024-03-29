/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Notify
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"time"
)

// NotifyV1Service struct for NotifyV1Service
type NotifyV1Service struct {
	// The unique string that we created to identify the Service resource.
	Sid *string `json:"sid,omitempty"`
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Service resource.
	AccountSid *string `json:"account_sid,omitempty"`
	// The string that you assigned to describe the resource.
	FriendlyName *string `json:"friendly_name,omitempty"`
	// The date and time in GMT when the resource was created specified in [RFC 2822](https://www.ietf.org/rfc/rfc2822.txt) format.
	DateCreated *time.Time `json:"date_created,omitempty"`
	// The date and time in GMT when the resource was last updated specified in [RFC 2822](https://www.ietf.org/rfc/rfc2822.txt) format.
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	// The SID of the [Credential](https://www.twilio.com/docs/notify/api/credential-resource) to use for APN Bindings.
	ApnCredentialSid *string `json:"apn_credential_sid,omitempty"`
	// The SID of the [Credential](https://www.twilio.com/docs/notify/api/credential-resource) to use for GCM Bindings.
	GcmCredentialSid *string `json:"gcm_credential_sid,omitempty"`
	// The SID of the [Credential](https://www.twilio.com/docs/notify/api/credential-resource) to use for FCM Bindings.
	FcmCredentialSid *string `json:"fcm_credential_sid,omitempty"`
	// The SID of the [Messaging Service](https://www.twilio.com/docs/sms/quickstart#messaging-services) to use for SMS Bindings. In order to send SMS notifications this parameter has to be set.
	MessagingServiceSid *string `json:"messaging_service_sid,omitempty"`
	// Deprecated.
	FacebookMessengerPageId *string `json:"facebook_messenger_page_id,omitempty"`
	// The protocol version to use for sending APNS notifications. Can be overridden on a Binding by Binding basis when creating a [Binding](https://www.twilio.com/docs/notify/api/binding-resource) resource.
	DefaultApnNotificationProtocolVersion *string `json:"default_apn_notification_protocol_version,omitempty"`
	// The protocol version to use for sending GCM notifications. Can be overridden on a Binding by Binding basis when creating a [Binding](https://www.twilio.com/docs/notify/api/binding-resource) resource.
	DefaultGcmNotificationProtocolVersion *string `json:"default_gcm_notification_protocol_version,omitempty"`
	// The protocol version to use for sending FCM notifications. Can be overridden on a Binding by Binding basis when creating a [Binding](https://www.twilio.com/docs/notify/api/binding-resource) resource.
	DefaultFcmNotificationProtocolVersion *string `json:"default_fcm_notification_protocol_version,omitempty"`
	// Whether to log notifications. Can be: `true` or `false` and the default is `true`.
	LogEnabled *bool `json:"log_enabled,omitempty"`
	// The absolute URL of the Service resource.
	Url *string `json:"url,omitempty"`
	// The URLs of the Binding, Notification, Segment, and User resources related to the service.
	Links *map[string]interface{} `json:"links,omitempty"`
	// Deprecated.
	AlexaSkillId *string `json:"alexa_skill_id,omitempty"`
	// Deprecated.
	DefaultAlexaNotificationProtocolVersion *string `json:"default_alexa_notification_protocol_version,omitempty"`
	// URL to send delivery status callback.
	DeliveryCallbackUrl *string `json:"delivery_callback_url,omitempty"`
	// Callback configuration that enables delivery callbacks, default false
	DeliveryCallbackEnabled *bool `json:"delivery_callback_enabled,omitempty"`
}
