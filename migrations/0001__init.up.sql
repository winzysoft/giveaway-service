-- ENUMS
CREATE TYPE "GiveawayStatus" AS ENUM ('CREATED', 'ACTIVE', 'DRAFT', 'ENDED', 'CANCELED', 'DELETED');
CREATE TYPE "GiveawayPlanningStatus" AS ENUM ('PENDING', 'IN_PROGRESS', 'SUCCESS', 'FAILED');
CREATE TYPE "ConditionTargetType" AS ENUM ('CHANNEL', 'CHAT');
CREATE TYPE "ConditionActionType" AS ENUM ('SUBSCRIBE');
CREATE TYPE "MediaType" AS ENUM ('IMAGE', 'VIDEO');
CREATE TYPE "GiveawayPlanningType" AS ENUM ('PUBLISH', 'RESULT');

-- TABLES
CREATE TABLE "users"
(
    "id"          UUID PRIMARY KEY,
    "telegram_id" TEXT         NOT NULL UNIQUE,
    "first_name"  TEXT         NOT NULL,
    "last_name"   TEXT,
    "avatar_url"  TEXT,
    "timezone"    TEXT         NOT NULL,
    "language"    TEXT         NOT NULL,
    "created_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "media"
(
    "id"         UUID PRIMARY KEY,
    "type"       "MediaType"  NOT NULL,
    "url"        TEXT         NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "chats"
(
    "id"                UUID PRIMARY KEY,
    "telegram_id"       TEXT         NOT NULL UNIQUE,
    "telegram_username" TEXT         NOT NULL,
    "created_at"        TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"        TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user_chats"
(
    "id"         UUID PRIMARY KEY,
    "user_id"    UUID         NOT NULL REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "chat_id"    UUID         NOT NULL REFERENCES "chats" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("user_id", "chat_id") -- Оставляем, так как пользователь может быть только один раз в чате
);

CREATE TABLE "giveaways"
(
    "id"                  UUID PRIMARY KEY,
    "status"              "GiveawayStatus" NOT NULL,
    "chat_id"             UUID             NOT NULL REFERENCES "chats" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "message_id"          TEXT,
    "should_edit_message" BOOLEAN          NOT NULL DEFAULT true,
    "should_send_new"     BOOLEAN          NOT NULL DEFAULT false,
    "win_count"           INTEGER          NOT NULL DEFAULT 1,
    "created_at"          TIMESTAMP(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"          TIMESTAMP(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "giveaway_managers"
(
    "id"          UUID PRIMARY KEY,
    "giveaway_id" UUID         NOT NULL REFERENCES "giveaways" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "user_id"     UUID         NOT NULL REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "created_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "giveaway_meta"
(
    "id"           UUID PRIMARY KEY,
    "giveaway_id"  UUID         NOT NULL UNIQUE REFERENCES "giveaways" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "title"        TEXT         NOT NULL,
    "description"  TEXT         NOT NULL,
    "button_title" TEXT         NOT NULL DEFAULT 'Participate',
    "media_id"     UUID REFERENCES "media" ("id") ON DELETE SET NULL ON UPDATE CASCADE,
    "created_at"   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "giveaway_planning"
(
    "id"          UUID PRIMARY KEY,
    "giveaway_id" UUID                     NOT NULL REFERENCES "giveaways" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "date"        TIMESTAMP(3)             NOT NULL,
    "type"        "GiveawayPlanningType"   NOT NULL,
    "status"      "GiveawayPlanningStatus" NOT NULL,
    "created_at"  TIMESTAMP(3)             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP(3)             NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "giveaway_conditions"
(
    "id"          UUID PRIMARY KEY,
    "giveaway_id" UUID                  NOT NULL REFERENCES "giveaways" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "target_type" "ConditionTargetType" NOT NULL,
    "action_type" "ConditionActionType" NOT NULL,
    "target_id"   TEXT                  NOT NULL,
    "created_at"  TIMESTAMP(3)          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP(3)          NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица участников
CREATE TABLE "giveaway_participants"
(
    "id"          UUID PRIMARY KEY,
    "giveaway_id" UUID         NOT NULL REFERENCES "giveaways" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "user_id"     UUID         NOT NULL REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "is_win"      BOOLEAN      NOT NULL DEFAULT false,
    "created_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX "idx_giveaway_managers_user_id" ON "giveaway_managers" ("user_id");
CREATE INDEX "idx_giveaway_participants_user_id" ON "giveaway_participants" ("user_id");
CREATE INDEX idx_participants_winners ON giveaway_participants (giveaway_id, user_id) WHERE is_win is true;
CREATE INDEX "idx_giveaway_planning_date" ON "giveaway_planning" ("date");
CREATE INDEX "idx_giveaways_chat_id" ON "giveaways" ("chat_id");