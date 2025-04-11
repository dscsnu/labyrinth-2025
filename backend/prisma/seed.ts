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

async function seedLevelTable() {
    const levels = [
        { name: 'Ex 1', targetScore: 10 },
        { name: 'Ex 2', targetScore: 20 },
        { name: 'Ex 3', targetScore: 30 },
        { name: 'Ex 4', targetScore: 40 },
        { name: 'Ex 5', targetScore: 50 },
        { name: 'Ex 6', targetScore: 60 },
        { name: 'Ex 7', targetScore: 120 },
    ];

    await prisma.level.createMany({
        data: levels
    });
}

async function seedSpellTable() {
    const spells = [
        { rewardScore: 10, numPoints: 2, cooldown: 1 },
        { rewardScore: 20, numPoints: 3, cooldown: 2 },
        { rewardScore: 30, numPoints: 3, cooldown: 2 },
        { rewardScore: 30, numPoints: 3, cooldown: 2 },
        { rewardScore: 40, numPoints: 4, cooldown: 3 },
        { rewardScore: 40, numPoints: 4, cooldown: 3 },
        { rewardScore: 50, numPoints: 4, cooldown: 3 },
        { rewardScore: 100, numPoints: 8, cooldown: 5 },
    ];

    await prisma.spell.createMany({
        data: spells
    });
}

async function main() {
    await createUserTriggers()
        .then(() => console.log('✅ userTriggers created'))
        .catch((e) => console.error(`🚨 ${e}`));

    await seedLevelTable()
        .then(() => console.log('✅ seeded Level table'))
        .catch((e) => console.error(`🚨 ${e}`));

    await seedSpellTable()
        .then(() => console.log('✅ seeded Spell table'))
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