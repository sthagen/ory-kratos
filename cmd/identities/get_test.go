// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package identities_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ory/kratos/cmd/identities"
	"github.com/ory/x/assertx"

	"github.com/ory/kratos/x"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/driver/config"
	"github.com/ory/kratos/identity"
)

func TestGetCmd(t *testing.T) {
	reg, cmd := setup(t, identities.NewGetIdentityCmd)

	t.Run("case=gets a single identity", func(t *testing.T) {
		i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
		i.MetadataPublic = []byte(`"public"`)
		i.MetadataAdmin = []byte(`"admin"`)
		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))

		stdOut := cmd.ExecNoErr(t, i.ID.String())

		ij, err := json.Marshal(identity.WithCredentialsNoConfigAndAdminMetadataInJSON(*i))
		require.NoError(t, err)

		assertx.EqualAsJSONExcept(t, json.RawMessage(ij), json.RawMessage(stdOut), []string{"created_at", "updated_at", "AdditionalProperties"})
	})

	t.Run("case=gets three identities", func(t *testing.T) {
		is, ids := makeIdentities(t, reg, 3)

		stdOut := cmd.ExecNoErr(t, ids...)

		isj, err := json.Marshal(is)
		require.NoError(t, err)

		assertx.EqualAsJSONExcept(t, json.RawMessage(isj), json.RawMessage(stdOut), []string{"created_at", "updated_at", "AdditionalProperties"})
	})

	t.Run("case=fails with unknown ID", func(t *testing.T) {
		stdErr := cmd.ExecExpectedErr(t, x.NewUUID().String())

		assert.Contains(t, stdErr, "Unable to locate the resource", stdErr)
	})

	t.Run("case=gets a single identity with oidc credentials", func(t *testing.T) {
		applyCredentials := func(identifier, accessToken, refreshToken, idToken string, encrypt bool) identity.Credentials {
			toJson := func(c identity.CredentialsOIDC) []byte {
				out, err := json.Marshal(&c)
				require.NoError(t, err)
				return out
			}
			transform := func(token string) string {
				if encrypt {
					s, err := reg.Cipher(context.Background()).Encrypt(context.Background(), []byte(token))
					require.NoError(t, err)
					return s
				}
				return token
			}
			return identity.Credentials{
				Type:        identity.CredentialsTypeOIDC,
				Identifiers: []string{"bar:" + identifier},
				Config: toJson(identity.CredentialsOIDC{Providers: []identity.CredentialsOIDCProvider{
					{
						Subject:             "foo",
						Provider:            "bar",
						InitialAccessToken:  transform(accessToken + "0"),
						InitialRefreshToken: transform(refreshToken + "0"),
						InitialIDToken:      transform(idToken + "0"),
						Organization:        "foo-org-id",
					},
					{
						Subject:             "baz",
						Provider:            "zab",
						InitialAccessToken:  transform(accessToken + "1"),
						InitialRefreshToken: transform(refreshToken + "1"),
						InitialIDToken:      transform(idToken + "1"),
						Organization:        "bar-org-id",
					},
				}}),
			}
		}
		i := identity.NewIdentity(config.DefaultIdentityTraitsSchemaID)
		i.MetadataPublic = []byte(`"public"`)
		i.MetadataAdmin = []byte(`"admin"`)
		i.SetCredentials(identity.CredentialsTypeOIDC, applyCredentials("uniqueIdentifier", "accessBar", "refreshBar", "idBar", true))
		// duplicate identity with decrypted tokens
		di := i.CopyWithoutCredentials()
		di.SetCredentials(identity.CredentialsTypeOIDC, applyCredentials("uniqueIdentifier", "accessBar", "refreshBar", "idBar", false))

		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))

		stdOut := cmd.ExecNoErr(t, "--"+identities.FlagIncludeCreds, "oidc", i.ID.String())
		ij, err := json.Marshal(identity.WithCredentialsAndAdminMetadataInJSON(*di))
		require.NoError(t, err)

		ii := []string{"id", "schema_url", "state_changed_at", "created_at", "updated_at", "credentials.oidc.created_at", "credentials.oidc.updated_at", "credentials.oidc.version", "AdditionalProperties"}
		assertx.EqualAsJSONExcept(t, json.RawMessage(ij), json.RawMessage(stdOut), ii)
	})
}
