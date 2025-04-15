import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { TOKEN_NAME } from "$lib/stores/TeamStore";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const teamCookie = cookies.get(TOKEN_NAME);
    if (teamCookie) {
        // Check for malformed cookie
        try {
            JSON.parse(teamCookie);
        } catch (e) {
            console.error('ERROR: Malformed TeamCookie > ', e);
            cookies.delete(TOKEN_NAME, { path: '/' });
            return {};
        }
        
        redirect(303, '/team');
    }
};