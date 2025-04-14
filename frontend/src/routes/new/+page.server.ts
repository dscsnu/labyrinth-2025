import { TeamStore, TOKEN_NAME, type ITeamData } from "$lib/stores/TeamStore";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { fetchWithAuth } from "$lib/utils/fetchWithAuth";

const resolveTeam = async (userId: string) => {
    const param = new URLSearchParams({ user_id: userId })
    const response = await fetchWithAuth(`/api/team?${param.toString()}`);

    if (response.ok) {
        const data = await response.json();
        console.log(data);

        const team: ITeamData = {
            id: data.id,
            name: data.name,
            allReady: data.members.every((m: any) => m.is_ready === true),
            members: data.members.map((m: any) => ({
                id: m.id,
                name: m.name,
                isReady: m.is_ready
            })),
        };

        TeamStore.set(team);
    } else {
        TeamStore.set(null);
    }
}

// @ts-ignore
export const load: PageServerLoad = async ({ cookies, locals }) => {
    const teamCookie = cookies.get(TOKEN_NAME);

    // If they already have a team, redirect to team page
    if (teamCookie) {
        // try catch is here to handle the case where the cookie is malformed
        try {
            JSON.parse(decodeURIComponent(teamCookie));
            throw redirect(303, '/team');
        } catch (e) {
            console.error('ERROR: Malformed TeamCookie > ', e);
            cookies.delete(TOKEN_NAME, { path: '/' });
            resolveTeam(locals.user?.id!);
        }
    } else {
        resolveTeam(locals.user?.id!);
    }
};