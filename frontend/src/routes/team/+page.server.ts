import { TOKEN_NAME as TEAM_TOKEN_NAME, type ITeamData } from "$lib/stores/TeamStore";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const teamCookie = cookies.get(TEAM_TOKEN_NAME);
    if (teamCookie) {
        try {
            JSON.parse(teamCookie)
        } catch (e) {
            console.error(`Error: Malformed TeamCookie > ${e}`)
            cookies.delete(TEAM_TOKEN_NAME, { path: '/' });
            redirect(303, '/new');
        }
        // valid team cookie

        const teamData: ITeamData = JSON.parse(teamCookie);
        const allReady = teamData.members.every(m => m.isReady);

        // add check for redirecting to /game
        if (allReady) {
        }

        return {};
    } else {
        redirect(303, '/new');
    }

    return {};
}