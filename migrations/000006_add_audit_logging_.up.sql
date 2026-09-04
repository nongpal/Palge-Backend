CREATE TABLE IF NOT EXISTS audit_logs (
    id bigserial PRIMARY KEY,
    user_id bigint REFERENCES users ON DELETE SET NULL,
    action text NOT NULL,
    resource text NOT NULL,
    resource_id bigint,
    ip_address inet,
    user_agent text,
    result text NOT NULL,
    details jsonb,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS audit_logs_user_id_idx
    ON audit_logs (user_id);

CREATE INDEX IF NOT EXISTS audit_logs_created_ad_idx
    ON audit_logs (created_at);

CREATE INDEX IF NOT EXISTS audit_logs_resource_idx
    ON audit_logs (resource, resource_id);