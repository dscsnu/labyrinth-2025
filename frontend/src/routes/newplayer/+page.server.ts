import { TEAM_TOKEN_NAME } from "$lib/stores/TeamStore";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const teamCookie = cookies.get(TEAM_TOKEN_NAME);
    if (!teamCookie) {
        // get current team and set cookie and redirect or return
    } else {
        return redirect(303, '/team');
    }
}