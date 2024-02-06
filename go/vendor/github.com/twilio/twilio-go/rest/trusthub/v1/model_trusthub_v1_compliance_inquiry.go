/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Trusthub
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

// TrusthubV1ComplianceInquiry struct for TrusthubV1ComplianceInquiry
type TrusthubV1ComplianceInquiry struct {
	// The unique ID used to start an embedded compliance registration session.
	InquiryId *string `json:"inquiry_id,omitempty"`
	// The session token used to start an embedded compliance registration session.
	InquirySessionToken *string `json:"inquiry_session_token,omitempty"`
	// The CustomerID matching the Customer Profile that should be resumed or resubmitted for editing.
	CustomerId *string `json:"customer_id,omitempty"`
	// The URL of this resource.
	Url *string `json:"url,omitempty"`
}
