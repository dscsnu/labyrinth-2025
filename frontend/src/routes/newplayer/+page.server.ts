import { TEAM_TOKEN_NAME } from "$lib/stores/TeamStore";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const teamCookie = cookies.get(TEAM_TOKEN_NAME);

    // If they already have a team, redirect to team page
    if (teamCookie) {
        // try catch is here to handle the case where the cookie is malformed
        try {
            const teamData = JSON.parse(decodeURIComponent(teamCookie));
            throw redirect(303, '/team');
        } catch (e) {
            console.error('Failed to parse team cookie:', e);
            cookies.delete(TEAM_TOKEN_NAME, { path: '/newplayer' });
        }
    }

    return {};
};