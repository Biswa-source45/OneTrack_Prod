-- Migration 010: Feedback
-- Any logged-in user (staff or customer) can report an issue with an optional
-- screenshot; only users with the 'superadmin' role triage it.
-- Safe to re-run.

-- NOT EXISTS instead of ON CONFLICT: doesn't depend on how the unique index
-- on role_name was originally created.
INSERT INTO master_roles (role_name, created_at)
SELECT 'superadmin', NOW()
WHERE NOT EXISTS (SELECT 1 FROM master_roles WHERE LOWER(TRIM(role_name)) = 'superadmin');

CREATE TABLE IF NOT EXISTS feedbacks (
    id             BIGSERIAL PRIMARY KEY,
    -- Reporter is either a users row or a contacts row, so no FK; name/email
    -- are snapshotted so the inbox still reads correctly if the reporter is deleted.
    reporter_id    BIGINT       NOT NULL,
    reporter_type  VARCHAR(20)  NOT NULL CHECK (reporter_type IN ('user', 'contact')),
    reporter_name  VARCHAR(255) NOT NULL,
    reporter_email VARCHAR(255),
    title          VARCHAR(200) NOT NULL,
    description    TEXT         NOT NULL,
    page_url       VARCHAR(500),
    -- Relative to the backend's working dir, e.g. uploads/feedback/<random>.png
    image_path     VARCHAR(500),
    status         VARCHAR(20)  NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'IN PROGRESS', 'COMPLETED')),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_status_created ON feedbacks (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedbacks_reporter ON feedbacks (reporter_type, reporter_id, created_at DESC);
