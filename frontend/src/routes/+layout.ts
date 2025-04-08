import { PUBLIC_SUPABASE_ANON_KEY, PUBLIC_SUPABASE_URL } from "$env/static/public";
import { createBrowserClient, createServerClient, isBrowser } from "@supabase/ssr";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ data, depends, fetch }) => {
    depends('supabase:auth');

    const supabase = isBrowser()
        ? createBrowserClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
            global: {
                fetch,
            },
        })
        : createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
            global: {
                fetch,
            },
            cookies: {
                getAll() {
                    return data.cookies
                },
            },
        });

    return {
        supabase,
        session: data.session,
        user: data.user,
    }
}