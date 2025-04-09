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

                UPDATE public.UserProfile SET
                    name = COALESCE(new.raw_user_meta_data->>'full_name', new.raw_user_meta_data->>'name', ''),
                    email = COALESCE(new.email, '')
                WHERE id = new.id;

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

async function main() {
    await createUserTriggers()
        .then(() => console.log('✅ userTriggers created'))
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