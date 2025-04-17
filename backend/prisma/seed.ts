import { Prisma, PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

async function createUserTriggers() {
    const queries = [
        Prisma.sql`
            CREATE OR REPLACE FUNCTION fnOnAuthNewUser()
            RETURNS trigger
            LANGUAGE plpgsql
            SECURITY DEFINER
            SET search_path = ''
            AS $$
            BEGIN
                INSERT INTO public.UserProfile (id, name, email)
                VALUES (new.id, '', '');

                -- update default data
                UPDATE public.UserProfile SET
                    name = COALESCE(new.raw_user_meta_data->>'full_name', new.raw_user_meta_data->>'name', ''),
                    email = COALESCE(new.email, '')
                WHERE id = new.id;

                -- update custom claims
                UPDATE auth.users SET
                    raw_app_meta_data = jsonb_set(
                        COALESCE(raw_app_meta_data, '{}'::jsonb),
                        '{custom_claims}',
                        '{"role": "PLAYER"}',
                        true
                    ) WHERE id = new.id;

                RETURN NEW;
            END;
            $$;
        `,
        Prisma.sql`
            CREATE OR REPLACE TRIGGER trOnAuthNewUser
                AFTER INSERT ON auth.users
                FOR EACH ROW EXECUTE PROCEDURE fnOnAuthNewUser();
        `,
        Prisma.sql`
            CREATE OR REPLACE FUNCTION fnOnUserUpdate()
            RETURNS TRIGGER
            LANGUAGE plpgsql
            SECURITY DEFINER
            SET search_path = ''
            AS $$
            BEGIN
                UPDATE auth.users
                SET raw_app_meta_data = jsonb_set(
                    COALESCE(raw_app_meta_data, '{}'::jsonb),
                    '{custom_claims}',
                    json_build_object(
                        'role', NEW.role
                    )::jsonb,
                    true
                ) WHERE id = NEW.id;

                RETURN NEW;
            END;
            $$;
        `,
        Prisma.sql`
            CREATE OR REPLACE TRIGGER trOnUserUpdate
                AFTER UPDATE OF role ON public.UserProfile
                FOR EACH ROW
                WHEN (
                    OLD.role IS DISTINCT FROM NEW.role
                ) EXECUTE PROCEDURE fnOnUserUpdate();
        `,
        Prisma.sql`
            CREATE OR REPLACE FUNCTION fnOnAuthDeleteUser()
            RETURNS trigger
            LANGUAGE plpgsql
            SECURITY DEFINER
            SET search_path = ''
            AS $$
            BEGIN
                DELETE FROM public.UserProfile WHERE id = old.id;
                RETURN old;
            END;
            $$;
        `,
        Prisma.sql`
            CREATE OR REPLACE TRIGGER trOnAuthDeleteUser
                AFTER DELETE ON auth.users
                FOR EACH ROW EXECUTE PROCEDURE fnOnAuthDeleteUser();
        `,
    ];

    for (const query of queries) {
        await prisma.$executeRaw(query);
    }
}

async function createUpdatedAtTriggers() {
    const queries = [
        Prisma.sql`
            CREATE OR REPLACE FUNCTION fnUpdateUpdatedAt()
            RETURNS trigger
            LANGUAGE plpgsql
            SECURITY DEFINER
            SET search_path = ''
            AS $$
            BEGIN
                NEW.updated_at = now();
                RETURN NEW;
            END;
            $$
        `,
        Prisma.sql`
            CREATE TRIGGER trUpdateUpdatedAtTeamLevelAssignment
                BEFORE UPDATE ON teamlevelassignment
                FOR EACH ROW EXECUTE PROCEDURE fnUpdateUpdatedAt();
        `,
        Prisma.sql`
            CREATE TRIGGER trUpdateUpdatedAtTeamSpellAttempt
                BEFORE UPDATE ON teamspellattempt
                FOR EACH ROW EXECUTE PROCEDURE fnUpdateUpdatedAt();
        `
    ];

    for (const query of queries) {
        await prisma.$executeRaw(query);
    }
}

async function seedLevelTable() {
    const levels = [
        { id: 1, name: 'Ex 1', targetScore: 10 },
        { id: 2, name: 'Ex 2', targetScore: 20 },
        { id: 3, name: 'Ex 3', targetScore: 30 },
        { id: 4, name: 'Ex 4', targetScore: 40 },
        { id: 5, name: 'Ex 5', targetScore: 50 },
        { id: 6, name: 'Ex 6', targetScore: 60 },
        { id: 7, name: 'Ex 7', targetScore: 120 },
    ];

    await prisma.level.createMany({
        data: levels
    });
}

async function seedSpellTable() {
    const spells = [
        { rewardScore: 10, numVertex: 2, cooldown: 1 },
        { rewardScore: 20, numVertex: 3, cooldown: 2 },
        { rewardScore: 30, numVertex: 3, cooldown: 2 },
        { rewardScore: 30, numVertex: 3, cooldown: 2 },
        { rewardScore: 40, numVertex: 4, cooldown: 3 },
        { rewardScore: 40, numVertex: 4, cooldown: 3 },
        { rewardScore: 50, numVertex: 4, cooldown: 3 },
        { rewardScore: 100, numVertex: 8, cooldown: 5 },
    ];

    await prisma.spell.createMany({
        data: spells
    });
}

async function seedLocationTable() {
    const locationsWithoutId = [
        { name: 'Tower 6 FruitShop', latitude: 28.52836141567019, longitude: 77.57777379378554 },
        { name: 'Research Block', latitude: 28.527471246410833, longitude: 77.57890336254033 },
        { name: 'A Block', latitude: 28.52692140098942, longitude: 77.57706407672785 },
        { name: 'CnD Atrium', latitude: 28.525519393291614, longitude: 77.57653461322104 },
        { name: 'G Block', latitude: 28.52799850829152, longitude: 77.5749038301257 },
        { name: 'Shopping Arcade', latitude: 28.52723498480195, longitude: 77.57292972432113 },
        { name: 'UAC', latitude: 28.523582249548888, longitude: 77.57437275274097 },
        { name: 'Cluster 1', latitude: 28.52440905859402, longitude: 77.57308311017349 },
        { name: 'Dibang Cycle Shop', latitude: 28.525235547619324, longitude: 77.57072373132858 },
        { name: 'Dining Hall 3', latitude: 28.52316003957018, longitude: 77.5696713403205 },
        { name: 'Indoor Sports Complex', latitude: 28.521496152454514, longitude: 77.5712575598366 },
    ];

    const locations = locationsWithoutId.map(l => {
        const id = Math.floor(100000 + Math.random() * 900000).toString();
        return { id, ...l }
    });

    await prisma.location.createMany({
        data: locations
    });
}

async function main() {
    await createUserTriggers()
        .then(() => console.log('✅ userTriggers created'))
        .catch((e) => console.error(`🚨 ${e}`));

    await createUpdatedAtTriggers()
        .then(() => console.log('✅ updatedAtTriggers created'))
        .catch((e) => console.error(`🚨 ${e}`));

    await seedLevelTable()
        .then(() => console.log('✅ seeded Level table'))
        .catch((e) => console.error(`🚨 ${e}`));

    await seedSpellTable()
        .then(() => console.log('✅ seeded Spell table'))
        .catch((e) => console.error(`🚨 ${e}`));

    await seedLocationTable()
        .then(() => console.log('✅ seeded Location table'))
        .catch((e) => console.error(`🚨 ${e}`));
}

main()
    .then(async () => {
        await prisma.$disconnect();
    })
    .catch(async (e) => {
        console.error(e);
        await prisma.$disconnect();
        process.exit(1);
    });