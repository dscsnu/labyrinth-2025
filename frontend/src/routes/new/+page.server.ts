import { redirect, type Cookies } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { TOKEN_NAME as TEAM_TOKEN_NAME, type ITeamData } from "$lib/stores/TeamStore";
import { TOKEN_NAME as JWT_TOKEN_NAME } from "$lib/stores/JwtTokenStore";
import { PUBLIC_BACKEND_URL } from "$env/static/public";
import type { User } from "@supabase/supabase-js";

const validateTeam = async (cookies: Cookies, user: User) => {
    const jwtCookie = cookies.get(JWT_TOKEN_NAME);
    console.log(jwtCookie)
    if (jwtCookie) {
        const cleanedUrl = PUBLIC_BACKEND_URL.replace(/\/+$/, "");
        const params = new URLSearchParams({ user_id: user.id });

        const response = await fetch(`${cleanedUrl}/api/team?${params.toString()}`, {
            headers: {
                'Authorization': `Bearer ${jwtCookie}`
            }
        });

        if (!response.ok) return;

        const data = await response.json();
        const teamData: ITeamData = {
            id: data.team_id,
            name: data.name,
            members: data.members.map((m: any) => ({
                id: m.id,
                name: m.name,
                email: m.email,
                isReady: m.isReady
            })),
        };

        cookies.set(TEAM_TOKEN_NAME, JSON.stringify(teamData), { path: '/', encode: (v) => v });
        redirect(303, '/team')
    }
}

// @ts-ignore
export const load: PageServerLoad = async ({ cookies, locals: { user } }) => {
    const teamCookie = cookies.get(TEAM_TOKEN_NAME);
    if (teamCookie) {
        // Check for malformed cookie
        try {
            JSON.parse(teamCookie);
        } catch (e) {
            console.error('ERROR: Malformed TeamCookie > ', e);
            cookies.delete(TEAM_TOKEN_NAME, { path: '/' });
            await validateTeam(cookies, user!);
            return {};
        }

        redirect(303, '/team');
    } else {
        await validateTeam(cookies, user!);
    }

    return {};
};