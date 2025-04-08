import { redirect } from '@sveltejs/kit';

export const GET = async (event) => {
    const {
        url,
        locals: { supabase }
    } = event;

    const code = url.searchParams.get('code') as string;
    const next = url.searchParams.get('next') ?? '/';
    const errorCode = url.searchParams.get('error') as string;

    if (errorCode) {
        redirect(303, '/');
    }

    if (code) {
        const { error } = await supabase.auth.exchangeCodeForSession(code);
        if (!error) {
            redirect(303, `/${next.slice(1)}`)
        }
    }

    redirect(303, '/api/auth/error');
}