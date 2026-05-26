CREATE TYPE "request_status_enum" AS ENUM ('PENDING_ANALYSIS', 'PROCESSED');
CREATE TYPE "priority_enum" AS ENUM ('LOW', 'MEDIUM', 'HIGH');

CREATE TABLE "Clients" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "request_type" VARCHAR(50) NOT NULL,
    "status" "request_status_enum" NOT NULL,
    "priority" "priority_enum" NOT NULL,
    "amount" NUMERIC NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);