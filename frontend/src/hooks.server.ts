import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from "$env/static/public";
import { TOKEN_NAME } from "$lib/stores/DeviceStore";
import { createServerClient } from "@supabase/ssr";

import { redirect, type Handle } from "@sveltejs/kit";
import { sequence } from "@sveltejs/kit/hooks";

const createSupabase: Handle = async ({ event, resolve }) => {
    event.locals.supabase = createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
        cookies: {
            getAll() {
                return event.cookies.getAll()
            },
            setAll(cookiesToSet) {
                cookiesToSet.forEach(({ name, value, options }) => {
                    event.cookies.set(name, value, { ...options, path: '/' })
                });
            },
        },
    })

    event.locals.safeGetSession = async () => {
        const {
            data: { user },
            error
        } = await event.locals.supabase.auth.getUser();
        if (error) {
            return { session: null, user: null }
        }

        const {
            data: { session }
        } = await event.locals.supabase.auth.getSession();
        if (!session) {
            return { session: null, user: null }
        }

        return { session, user }
    }

    event.locals.getSession = async () => {
        const { session } = await event.locals.safeGetSession();
        return session;
    }
    event.locals.session = await event.locals.getSession();

    event.locals.getUser = async () => {
        const { user } = await event.locals.safeGetSession();
        return user;
    }
    event.locals.user = await event.locals.getUser();

    return resolve(event, {
        filterSerializedResponseHeaders(name) {
            return name === 'content-range' || name === 'x-supabase-api-version'
        },
    });
}

const authGuard: Handle = async ({ event, resolve }) => {
    const currentPath = event.url.pathname;

    if (currentPath.startsWith('/api')) {
        return resolve(event);
    }

    if (currentPath !== '/' && !event.locals.user) {
        return redirect(303, '/');
    }

    return resolve(event);
}

const handleDevice: Handle = async({ event, resolve }) => {
    const response = await resolve(event, {
        transformPageChunk: ({ html }) => {
            const maxAge = 365 * 24 * 60 * 60;

            let currentDevice = event.cookies.get(TOKEN_NAME);
            if (!currentDevice) {
                const userOnMobile = event.request.headers.get('sec-ch-ua-mobile') === '?1';
                currentDevice = userOnMobile ? 'mobile' : 'desktop';

                event.cookies.set(TOKEN_NAME, currentDevice, {
                    path: '/',
                    expires: new Date(Date.now() + maxAge),
                    maxAge,
                    httpOnly: false,
                    sameSite: 'strict',
                })
            }

            return html
                .replace('data-device=""', `data-device="${currentDevice}"`)
        }
    })

    response.headers.set('Accept-CH', 'Sec-CH-UA-Mobile');
    response.headers.set('Vary', 'Sec-CH-UA-Mobile');
    response.headers.set('Critical-CH', 'Sec-CH-UA-Mobile');

    return response;
}

export const handle: Handle = sequence(createSupabase, authGuard, handleDevice)