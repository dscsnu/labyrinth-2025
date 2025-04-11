import { TOKEN_NAME } from '$lib/stores/TokenStore';
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

// @ts-ignore
export const load: PageServerLoad = async ({ cookies }) => {
    const userCookie = cookies.get(TOKEN_NAME);
    if (!userCookie) {
        // get current team and set cookie and redirect or return
    } else {
        return redirect(303, '/newplayer');
    }
}