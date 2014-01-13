// Package geotrigger_golang provides API access to the Geotrigger Service,
// a cloud based system of geofencing and push notifications. The library
// makes it easier to interact with the Geotrigger API as either a Device or
// an Application. This assumes you have a developer account on
// developers.arcgis.com, from which you can create an Application and obtain
// the necessary credentials to use with this golang library.
//
// For more information about the Geotrigger Service, please look here:
// https://developers.arcgis.com/en/geotrigger-service/
//
// Documentation for this library can be found on github:
// https://github.com/geoloqi/geotrigger_golang
package geotrigger_golang

// The client struct type. Has one, un-exported field for a session that handles
// auth for you. Make API requests with the "Request" method. This is the type
// you should use directly for interacting with the geotrigger API.
type Client struct {
	session session
}

// Create and register a new device associated with the provided client_id
// The channel that is returned will be written to once. If the read value is a nil,
// then the returned client pointer has been successfully inflated and is ready for use.
// Otherwise, the error will contain information about what went wrong.
func NewDeviceClient(clientId string) (*Client, chan error) {
	refreshStatusChecks := make(chan *refreshStatusCheck)
	device := &device{
		clientId:            clientId,
		refreshStatusChecks: refreshStatusChecks,
	}

	return getTokens(device)
}

// Create and register a new application associated with the provided client_id
// and client_secret.
// The channel that is returned will be written to once. If the read value is a nil,
// then the returned client pointer has been successfully inflated and is ready for use.
// Otherwise, the error will contain information about what went wrong.
func NewApplicationClient(clientId string, clientSecret string) (*Client, chan error) {
	application := &application{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	return getTokens(application)
}

// The method to use for making requests!
// `responseJSON` can be a struct modeling the expected JSON, or an arbitrary JSON map (map[string]interface{})
// that can be used with the helper method `GetValueFromJSONObject`.
// The channel that is returned will be written to once. If the read value is a nil,
// then the provided responseJSON has been successfully inflated and is ready for use.
// Otherwise, the error will contain information about what went wrong.
func (client *Client) Request(route string, params map[string]interface{}, responseJSON interface{}) chan error {
	errorChan := make(chan error)
	go client.session.geotriggerAPIRequest(route, params, responseJSON, errorChan)
	return errorChan
}

// Get info about the current session.
// If this is an application session, the following keys will be present:
// `access_token`
// `client_id`
// `client_secret`
// If this is a device session, the following keys will be present:
// `access_token`
// `refresh_token`
// `device_id`
// `client_id`
func (client *Client) GetSessionInfo() map[string]string {
	return client.session.getSessionInfo()
}

// Un-exported helper to just DRY up the client constructors above.
func getTokens(session session) (*Client, chan error) {
	errorChan := make(chan error)
	client := &Client{session: session}

	go session.requestAccess(errorChan)
	go session.tokenManager()
	return client, errorChan
}
