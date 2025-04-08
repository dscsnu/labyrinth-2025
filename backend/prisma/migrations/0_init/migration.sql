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
    "team_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "is_ready" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "teammember_pkey" PRIMARY KEY ("team_id","user_id")
);

-- CreateTable
CREATE TABLE "level" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "description" TEXT,
    "target_score" INTEGER NOT NULL,

    CONSTRAINT "level_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "teamlevelassignment" (
    "team_id" UUID NOT NULL,
    "level_id" UUID NOT NULL,
    "sequence" INTEGER NOT NULL,
    "current_score" INTEGER NOT NULL DEFAULT 0,
    "is_finished" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "teamlevelassignment_pkey" PRIMARY KEY ("team_id","level_id")
);

-- CreateTable
CREATE TABLE "spell" (
    "spell_id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "reward_score" INTEGER NOT NULL,
    "num_points" INTEGER NOT NULL,
    "cooldown" INTEGER NOT NULL,

    CONSTRAINT "spell_pkey" PRIMARY KEY ("spell_id")
);

-- CreateTable
CREATE TABLE "pattern" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "num_points" INTEGER NOT NULL,
    "spellId" UUID NOT NULL,

    CONSTRAINT "pattern_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "location" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "latitude" DOUBLE PRECISION NOT NULL,
    "longitude" DOUBLE PRECISION NOT NULL,

    CONSTRAINT "location_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "patternlocationassignment" (
    "patternId" UUID NOT NULL,
    "locationId" UUID NOT NULL,

    CONSTRAINT "patternlocationassignment_pkey" PRIMARY KEY ("patternId","locationId")
);

-- CreateIndex
CREATE UNIQUE INDEX "userprofile_email_key" ON "userprofile"("email");

-- CreateIndex
CREATE UNIQUE INDEX "team_id_key" ON "team"("id");

-- CreateIndex
CREATE UNIQUE INDEX "teammember_user_id_key" ON "teammember"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "teamlevelassignment_team_id_sequence_key" ON "teamlevelassignment"("team_id", "sequence");

-- AddForeignKey
ALTER TABLE "teammember" ADD CONSTRAINT "teammember_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teammember" ADD CONSTRAINT "teammember_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "userprofile"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlevelassignment" ADD CONSTRAINT "teamlevelassignment_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "teamlevelassignment" ADD CONSTRAINT "teamlevelassignment_level_id_fkey" FOREIGN KEY ("level_id") REFERENCES "level"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "pattern" ADD CONSTRAINT "pattern_spellId_fkey" FOREIGN KEY ("spellId") REFERENCES "spell"("spell_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "patternlocationassignment" ADD CONSTRAINT "patternlocationassignment_patternId_fkey" FOREIGN KEY ("patternId") REFERENCES "pattern"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "patternlocationassignment" ADD CONSTRAINT "patternlocationassignment_locationId_fkey" FOREIGN KEY ("locationId") REFERENCES "location"("id") ON DELETE CASCADE ON UPDATE CASCADE;

