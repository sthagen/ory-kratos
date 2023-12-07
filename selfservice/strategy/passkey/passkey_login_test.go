// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package passkey_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/kratos/driver/config"
	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/internal/testhelpers"
	"github.com/ory/kratos/selfservice/flow"
	"github.com/ory/kratos/selfservice/flow/login"
	"github.com/ory/kratos/selfservice/strategy/passkey"
	"github.com/ory/kratos/text"
	"github.com/ory/kratos/ui/node"
)

var (
	//go:embed fixtures/login/success/identity.json
	loginSuccessIdentity []byte
	//go:embed fixtures/login/success/credentials.json
	loginPasswordlessCredentials []byte
	//go:embed fixtures/login/success/internal_context.json
	loginPasswordlessContext []byte
	//go:embed fixtures/login/success/response.json
	loginPasswordlessResponse []byte
)

func TestCompleteLogin(t *testing.T) {
	t.Parallel()
	fix := newLoginFixture(t)

	t.Run("flow=passwordless", func(t *testing.T) {
		t.Run("case=passkey button exists", func(t *testing.T) {
			client := testhelpers.NewClientWithCookies(t)
			f := testhelpers.InitializeLoginFlowViaBrowser(t, client, fix.publicTS, false, true, false, false)
			testhelpers.SnapshotTExcept(t, f.Ui.Nodes, []string{
				"0.attributes.value",
				"2.attributes.nonce",
				"2.attributes.src",
				"5.attributes.value",
			})
		})

		t.Run("case=passkey shows error if user tries to sign in but no such user exists", func(t *testing.T) {
			payload := func(v url.Values) {
				v.Set("method", "passkey")
				v.Set(node.PasskeyLogin, string(loginPasswordlessResponse))
			}

			check := func(t *testing.T, shouldRedirect bool, body string, res *http.Response) {
				fix.checkURL(t, shouldRedirect, res)
				assert.NotEmpty(t, gjson.Get(body, "id").String(), "%s", body)
				assert.Equal(t, text.NewErrorValidationSuchNoWebAuthnUser().Text, gjson.Get(body, "ui.messages.0.text").String(), "%s", body)
			}

			t.Run("type=browser", func(t *testing.T) {
				body, res := fix.loginViaBrowser(t, false, payload, testhelpers.NewClientWithCookies(t))
				check(t, true, body, res)
			})

			t.Run("type=spa", func(t *testing.T) {
				body, res := fix.loginViaBrowser(t, true, payload, testhelpers.NewClientWithCookies(t))
				check(t, false, body, res)
			})
		})

		t.Run("case=should fail if passkey login is invalid", func(t *testing.T) {
			payload := func(v url.Values) {
				v.Set("method", "passkey")
				v.Set("passkey_login", "invalid passkey data")
			}

			check := func(t *testing.T, shouldRedirect bool, body string, res *http.Response) {
				fix.checkURL(t, shouldRedirect, res)
				assert.NotEmpty(t, gjson.Get(body, "id").String(), "%s", body)
				assert.Equal(t, "Unable to parse WebAuthn response.", gjson.Get(body, "ui.messages.0.text").String(), "%s", body)
			}

			t.Run("type=browser", func(t *testing.T) {
				body, res := fix.loginViaBrowser(t, false, payload, testhelpers.NewClientWithCookies(t))
				check(t, true, body, res)
			})

			t.Run("type=spa", func(t *testing.T) {
				body, res := fix.loginViaBrowser(t, true, payload, testhelpers.NewClientWithCookies(t))
				check(t, false, body, res)
			})
		})

		t.Run("case=succeeds with passwordless login", func(t *testing.T) {
			run := func(t *testing.T, spa bool) {
				fix.conf.MustSet(fix.ctx, config.ViperKeySessionWhoAmIAAL, "aal1")
				// We load our identity which we will use to replay the webauth session
				id := fix.createIdentityWithPasskey(t, identity.Credentials{
					Config:  loginPasswordlessCredentials,
					Version: 1,
				})

				browserClient := testhelpers.NewClientWithCookies(t)
				body, res, f := fix.submitWebAuthnLoginWithClient(t, spa, loginPasswordlessContext, browserClient, func(values url.Values) {
					values.Set(node.PasskeyLogin, string(loginPasswordlessResponse))
				}, testhelpers.InitFlowWithAAL(identity.AuthenticatorAssuranceLevel1))

				prefix := ""
				if spa {
					assert.Contains(t, res.Request.URL.String(), fix.publicTS.URL+login.RouteSubmitFlow)
					prefix = "session."
				} else {
					assert.Contains(t, res.Request.URL.String(), fix.redirTS.URL)
				}

				assert.True(t, gjson.Get(body, prefix+"active").Bool(), "%s", body)
				assert.EqualValues(t, identity.AuthenticatorAssuranceLevel1, gjson.Get(body, prefix+"authenticator_assurance_level").String(), "%s", body)
				assert.EqualValues(t, identity.CredentialsTypePasskey, gjson.Get(body, prefix+"authentication_methods.#(method==passkey).method").String(), "%s", body)
				assert.EqualValues(t, id.ID.String(), gjson.Get(body, prefix+"identity.id").String(), "%s", body)

				actualFlow, err := fix.reg.LoginFlowPersister().GetLoginFlow(context.Background(), uuid.FromStringOrNil(f.Id))
				require.NoError(t, err)
				assert.Empty(t, gjson.GetBytes(actualFlow.InternalContext, flow.PrefixInternalContextKey(identity.CredentialsTypePasskey, passkey.InternalContextKeySessionData)))
			}

			t.Run("type=browser", func(t *testing.T) {
				run(t, false)
			})

			t.Run("type=spa", func(t *testing.T) {
				run(t, true)
			})
		})
	})
}