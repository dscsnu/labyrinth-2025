-- CreateEnum
CREATE TYPE "UserRole" AS ENUM ('PLAYER', 'HELPER', 'ADMIN');

-- CreateTable
CREATE TABLE "gameconfig" (
    "property" TEXT NOT NULL,
    "value" TEXT NOT NULL,

    CONSTRAINT "gameconfig_pkey" PRIMARY KEY ("property")
);

-- CreateTable
CREATE TABLE "userprofile" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "role" "UserRole" NOT NULL DEFAULT 'PLAYER',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),

    CONSTRAINT "userprofile_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "team" (
    "id" VARCHAR(6) NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "team_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "teammember" (
    "team_id" VARCHAR(6) NOT NULL,
    "user_id" UUID NOT NULL,
    "is_ready" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "teammember_pkey" PRIMARY KEY ("team_id","user_id")
);

-- CreateTable
CREATE TABLE "level" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "image_url" TEXT,
    "target_score" INTEGER NOT NULL,

    CONSTRAINT "level_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "teamlevelassignment" (
    "team_id" VARCHAR(6) NOT NULL,
    "level_id" INTEGER NOT NULL,
    "sequence" INTEGER NOT NULL,
    "current_score" INTEGER NOT NULL DEFAULT 0,
    "is_finished" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "teamlevelassignment_pkey" PRIMARY KEY ("team_id","level_id")
);

-- CreateTable
CREATE TABLE "spell" (
    "spell_id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "reward_score" INTEGER NOT NULL,
    "num_vertex" INTEGER NOT NULL,
    "cooldown" INTEGER NOT NULL,
    "image_url" TEXT,

    CONSTRAINT "spell_pkey" PRIMARY KEY ("spell_id")
);

-- CreateTable
CREATE TABLE "pattern" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "num_vertex" INTEGER NOT NULL,
    "spell_id" UUID NOT NULL,

    CONSTRAINT "pattern_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "location" (
    "id" VARCHAR(6) NOT NULL,
    "name" TEXT NOT NULL,
    "latitude" DOUBLE PRECISION NOT NULL,
    "longitude" DOUBLE PRECISION NOT NULL,

    CONSTRAINT "location_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "patternlocationassignment" (
    "patternId" UUID NOT NULL,
    "locationId" VARCHAR(6) NOT NULL,

    CONSTRAINT "patternlocationassignment_pkey" PRIMARY KEY ("patternId","locationId")
);

-- CreateTable
CREATE TABLE "teamspellassignment" (
    "team_id" VARCHAR(6) NOT NULL,
    "spell_id" UUID NOT NULL,
    "remaining_cooldown" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "teamspellassignment_pkey" PRIMARY KEY ("team_id","spell_id")
);

-- CreateTable
CREATE TABLE "teamspellattempt" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "team_id" VARCHAR(6) NOT NULL,
    "spell_id" UUID NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),

    CONSTRAINT "teamspellattempt_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "teamlocationprogress" (
    "attempt_id" UUID NOT NULL,
    "location_id" VARCHAR(6) NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT now(),

    CONSTRAINT "teamlocationprogress_pkey" PRIMARY KEY ("attempt_id","location_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "userprofile_email_key" ON "userprofile"("email");

-- CreateIndex
CREATE UNIQUE INDEX "team_id_key" ON "team"("id");

-- CreateIndex
CREATE UNIQUE INDEX "teammember_user_id_key" ON "teammember"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "teamlevelassignment_team_id_sequence_key" ON "teamlevelassignment"("team_id", "sequence");

-- CreateIndex
CREATE UNIQUE INDEX "teamspellattempt_team_id_active_key" ON "teamspellattempt"("team_id", "active");

-- AddForeignKey
ALTER TABLE "teammember" ADD CONSTRAINT "teammember_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teammember" ADD CONSTRAINT "teammember_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "userprofile"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlevelassignment" ADD CONSTRAINT "teamlevelassignment_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlevelassignment" ADD CONSTRAINT "teamlevelassignment_level_id_fkey" FOREIGN KEY ("level_id") REFERENCES "level"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "pattern" ADD CONSTRAINT "pattern_spell_id_fkey" FOREIGN KEY ("spell_id") REFERENCES "spell"("spell_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "patternlocationassignment" ADD CONSTRAINT "patternlocationassignment_patternId_fkey" FOREIGN KEY ("patternId") REFERENCES "pattern"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "patternlocationassignment" ADD CONSTRAINT "patternlocationassignment_locationId_fkey" FOREIGN KEY ("locationId") REFERENCES "location"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamspellassignment" ADD CONSTRAINT "teamspellassignment_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamspellassignment" ADD CONSTRAINT "teamspellassignment_spell_id_fkey" FOREIGN KEY ("spell_id") REFERENCES "spell"("spell_id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamspellattempt" ADD CONSTRAINT "teamspellattempt_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamspellattempt" ADD CONSTRAINT "teamspellattempt_spell_id_fkey" FOREIGN KEY ("spell_id") REFERENCES "spell"("spell_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlocationprogress" ADD CONSTRAINT "teamlocationprogress_attempt_id_fkey" FOREIGN KEY ("attempt_id") REFERENCES "teamspellattempt"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlocationprogress" ADD CONSTRAINT "teamlocationprogress_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
