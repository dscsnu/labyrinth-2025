import { TOKEN_NAME } from "$lib/stores/TeamStore";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const teamCookie = cookies.get(TOKEN_NAME);

    // Check if user has a team
    if (!teamCookie) {
        throw redirect(303, '/newplayer');
    }

    // Try to parse team data
    let teamData
    try {
        teamData = JSON.parse(decodeURIComponent(teamCookie));
        // Check if team is ready and game start
        if (teamData.isReady) { //game start remaining
            throw redirect(303, '/game');
        }
    } catch (e) {
        console.error('Failed to parse team cookie:', e);
        cookies.delete(TOKEN_NAME, { path: '/' });
        throw redirect(303, '/newplayer');
    }

    return {};
};