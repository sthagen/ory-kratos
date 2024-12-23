CREATE INDEX IF NOT EXISTS session_devices_id_nid_idx ON session_devices (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS session_devices_session_id_nid_idx ON session_devices (session_id ASC, nid ASC);
DROP INDEX IF EXISTS session_devices_nid_idx;
DROP INDEX IF EXISTS session_devices_session_id_idx;

CREATE INDEX IF NOT EXISTS session_token_exchanges_nid_code_idx ON session_token_exchanges (init_code ASC, nid ASC);
CREATE INDEX IF NOT EXISTS session_token_exchanges_nid_flow_id_idx ON session_token_exchanges (flow_id ASC, nid ASC);
DROP INDEX IF EXISTS session_token_exchanges_flow_id_nid_init_code_idx;
DROP INDEX IF EXISTS session_token_exchanges_nid_init_code_idx;

CREATE INDEX IF NOT EXISTS courier_messages_status_idx ON  courier_messages (status ASC);
CREATE INDEX IF NOT EXISTS courier_messages_nid_id_idx ON  courier_messages (nid ASC, id ASC);
CREATE INDEX IF NOT EXISTS courier_messages_id_nid_idx ON  courier_messages (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS courier_messages_nid_created_at_id_idx ON  courier_messages (nid ASC, created_at DESC);
DROP INDEX IF EXISTS courier_messages_status_id_idx;
DROP INDEX IF EXISTS courier_messages_nid_id_created_at_idx;

CREATE INDEX IF NOT EXISTS continuity_containers_nid_id_idx ON continuity_containers (nid ASC, id ASC);
CREATE INDEX IF NOT EXISTS continuity_containers_id_nid_idx ON continuity_containers (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS continuity_containers_identity_id_nid_idx ON continuity_containers (identity_id ASC, nid ASC);
DROP INDEX IF EXISTS continuity_containers_identity_id_idx;
DROP INDEX IF EXISTS continuity_containers_nid_idx;

CREATE INDEX IF NOT EXISTS identity_verification_codes_nid_flow_id_idx ON identity_verification_codes (nid ASC, selfservice_verification_flow_id ASC);
CREATE INDEX IF NOT EXISTS identity_verification_codes_id_nid_idx ON identity_verification_codes (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS identity_verification_codes_verifiable_address_nid_idx ON identity_verification_codes (identity_verifiable_address_id ASC, nid ASC);
DROP INDEX IF EXISTS identity_verification_codes_identity_verifiable_address_id_idx;
DROP INDEX IF EXISTS identity_verification_codes_nid_idx;

CREATE INDEX IF NOT EXISTS identity_verification_tokens_nid_id_idx ON identity_verification_tokens (nid ASC, id ASC);
CREATE INDEX IF NOT EXISTS identity_verification_tokens_id_nid_idx ON identity_verification_tokens (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS identity_verification_tokens_token_nid_used_flow_id_idx ON identity_verification_tokens (nid ASC, token ASC, used ASC, selfservice_verification_flow_id ASC);
DROP INDEX IF EXISTS identity_verification_tokens_nid_idx;

CREATE INDEX IF NOT EXISTS identity_registration_codes_nid_flow_id_idx ON identity_registration_codes (nid ASC, selfservice_registration_flow_id ASC);
CREATE INDEX IF NOT EXISTS identity_registration_codes_id_nid_idx ON identity_registration_codes (id ASC, nid ASC);
DROP INDEX IF EXISTS identity_registration_codes_nid_idx;

CREATE INDEX IF NOT EXISTS identity_recovery_tokens_nid_id_idx  ON identity_recovery_tokens (nid ASC, id ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_tokens_id_nid_idx  ON identity_recovery_tokens (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_tokens_token_nid_used_idx  ON identity_recovery_tokens (nid ASC, token ASC, used ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_tokens_identity_id_nid_idx  ON identity_recovery_tokens (identity_id ASC, nid ASC);
DROP INDEX IF EXISTS identity_recovery_tokens_identity_id_idx;
DROP INDEX IF EXISTS identity_recovery_tokens_nid_idx;

CREATE INDEX IF NOT EXISTS identity_recovery_codes_nid_flow_id_idx ON identity_recovery_codes (nid ASC, selfservice_recovery_flow_id ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_codes_id_nid_idx ON identity_recovery_codes (id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_codes_identity_id_nid_idx ON identity_recovery_codes (identity_id ASC, nid ASC);
CREATE INDEX IF NOT EXISTS identity_recovery_codes_identity_recovery_address_id_nid_idx ON identity_recovery_codes (identity_recovery_address_id ASC, nid ASC);
DROP INDEX IF EXISTS identity_recovery_codes_identity_recovery_address_id_idx;
DROP INDEX IF EXISTS identity_recovery_codes_identity_id_idx;
DROP INDEX IF EXISTS identity_recovery_codes_nid_idx;

CREATE INDEX IF NOT EXISTS identity_login_codes_nid_flow_id_idx ON identity_login_codes (nid ASC, selfservice_login_flow_id ASC);
CREATE INDEX IF NOT EXISTS identity_login_codes_id_nid_idx ON identity_login_codes (id ASC, nid ASC);
DROP INDEX IF EXISTS identity_login_codes_nid_idx;
