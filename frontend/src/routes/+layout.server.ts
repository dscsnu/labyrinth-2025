import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals: { session, user }, cookies }) => {
    return {
        session,
        user,
        cookies: cookies.getAll()
    }
}