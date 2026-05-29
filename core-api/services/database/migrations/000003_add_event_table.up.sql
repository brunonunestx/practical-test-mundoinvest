CREATE TABLE "Events" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "event_id" VARCHAR(255) NOT NULL,
    "card_id" VARCHAR(255) NOT NULL,
    "client_email" VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_client_email ON "Events" (client_email);